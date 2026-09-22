// wf - wayfinder tracker helper. Keeps large map and decision-log bodies out of
// agent context: reads return only what was asked for, edits happen in memory
// and upload through gh with UTF-8 files (no PowerShell string round-trips).
//
// Stdlib only, single file, so it runs from any project (even inside another
// Go module) without a go.mod:
//
//	go run <wayfinder-skill-dir>/utilities/wf/wf.go <command> [flags]
//
// Commands:
//
//	log <log-issue> [--ids A,B] [--global] [--list] [--next]
//	log-append <log-issue> <rows.md> [--amend]
//	section <issue> <Section name>...
//	map-edit <map-issue> [ops...] [--dry-run]
//	body get <issue> <file>
//	body put <issue> <file>
//	validate <file>
//
// Global flag: -R owner/repo (default: repo of the current directory).
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var repoFlag string

func main() {
	args := os.Args[1:]
	if len(args) >= 2 && (args[0] == "-R" || args[0] == "--repo") {
		repoFlag = args[1]
		args = args[2:]
	}
	if len(args) == 0 {
		fail(usage)
	}
	var err error
	switch args[0] {
	case "log":
		err = cmdLog(args[1:])
	case "log-append":
		err = cmdLogAppend(args[1:])
	case "section":
		err = cmdSection(args[1:])
	case "map-edit":
		err = cmdMapEdit(args[1:])
	case "body":
		err = cmdBody(args[1:])
	case "validate":
		err = cmdValidate(args[1:])
	default:
		err = fmt.Errorf("unknown command %q\n%s", args[0], usage)
	}
	if err != nil {
		fail(err.Error())
	}
}

const usage = `usage: wf [-R owner/repo] <command>
  log <log-issue> [--ids A,B] [--global] [--list] [--next]
  log-append <log-issue> <rows.md> [--amend]
  section <issue> <Section name>...
  map-edit <map-issue> [--complete N --gist TEXT] [--coverage ID=status[=link]]...
           [--append "Section=line"]... [--remove "Section=substring"]... [--dry-run]
  body get <issue> <file>
  body put <issue> <file>
  validate <file>`

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "wf: "+msg)
	os.Exit(1)
}

// --- gh access ---

type issue struct {
	Body     string `json:"body"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Comments []struct {
		Body string `json:"body"`
	} `json:"comments"`
}

func gh(args ...string) ([]byte, error) {
	if repoFlag != "" {
		args = append(args, "--repo", repoFlag)
	}
	cmd := exec.Command("gh", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func fetchIssue(num string, withComments bool) (*issue, error) {
	fields := "body,title,url"
	if withComments {
		fields += ",comments"
	}
	out, err := gh("issue", "view", num, "--json", fields)
	if err != nil {
		return nil, err
	}
	var is issue
	if err := json.Unmarshal(out, &is); err != nil {
		return nil, err
	}
	is.Body = normalize(is.Body)
	for i := range is.Comments {
		is.Comments[i].Body = normalize(is.Comments[i].Body)
	}
	return &is, nil
}

// uploadBody writes body to a temp UTF-8 file, uploads it, and re-fetches to
// confirm GitHub stored exactly what was sent.
func uploadBody(num, body string) error {
	f, err := os.CreateTemp("", "wf-body-*.md")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(body); err != nil {
		f.Close()
		return err
	}
	f.Close()
	if _, err := gh("issue", "edit", num, "--body-file", f.Name()); err != nil {
		return err
	}
	back, err := fetchIssue(num, false)
	if err != nil {
		return fmt.Errorf("verify fetch: %w", err)
	}
	if strings.TrimSpace(back.Body) != strings.TrimSpace(body) {
		return fmt.Errorf("verify: body on #%s differs from upload (%d vs %d bytes)", num, len(back.Body), len(body))
	}
	return nil
}

func normalize(s string) string {
	s = strings.TrimPrefix(s, "\uFEFF")
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// --- decision log ---

var rowStart = regexp.MustCompile(`^\*\*([A-Z0-9]+(?:-[A-Z0-9]+)*-GM-(\d{3,}))\*\*`)

type gmRow struct {
	ID   string
	Num  int
	Text string
}

// parseRows returns GM rows in document order. A row runs from its **ID** line
// to the next row, heading, or horizontal rule.
func parseRows(text string) []gmRow {
	var rows []gmRow
	var cur *gmRow
	flush := func() {
		if cur != nil {
			cur.Text = strings.TrimSpace(cur.Text)
			rows = append(rows, *cur)
			cur = nil
		}
	}
	for _, line := range strings.Split(text, "\n") {
		if m := rowStart.FindStringSubmatch(line); m != nil {
			flush()
			n, _ := strconv.Atoi(m[2])
			cur = &gmRow{ID: m[1], Num: n, Text: line}
			continue
		}
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "---" {
			flush()
			continue
		}
		if cur != nil {
			cur.Text += "\n" + line
		}
	}
	flush()
	return rows
}

// logRows merges rows from the log body and its comments. A later entry with
// the same ID supersedes the earlier one (amendments).
func logRows(is *issue) (order []string, byID map[string]gmRow) {
	byID = map[string]gmRow{}
	add := func(text string) {
		for _, r := range parseRows(text) {
			if _, seen := byID[r.ID]; !seen {
				order = append(order, r.ID)
			}
			byID[r.ID] = r
		}
	}
	add(is.Body)
	for _, c := range is.Comments {
		add(c.Body)
	}
	return order, byID
}

func isGlobal(r gmRow) bool {
	return strings.Contains(r.Text, "[global]")
}

func nextID(order []string, byID map[string]gmRow) (string, error) {
	if len(order) == 0 {
		return "", errors.New("log has no GM rows yet; start at {MAP-SLUG}-GM-001")
	}
	max := byID[order[0]]
	for _, id := range order {
		if byID[id].Num > max.Num {
			max = byID[id]
		}
	}
	prefix := strings.TrimSuffix(max.ID, fmt.Sprintf("%03d", max.Num))
	return fmt.Sprintf("%s%03d", prefix, max.Num+1), nil
}

func cmdLog(args []string) error {
	if len(args) < 1 {
		return errors.New("log <log-issue> [--ids A,B] [--global] [--list] [--next]")
	}
	num := args[0]
	var ids []string
	global, list, next := false, false, false
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--ids":
			if i+1 >= len(args) {
				return errors.New("--ids needs a comma-separated list")
			}
			for _, id := range strings.Split(args[i+1], ",") {
				if id = strings.TrimSpace(id); id != "" {
					ids = append(ids, id)
				}
			}
			i++
		case "--global":
			global = true
		case "--list":
			list = true
		case "--next":
			next = true
		default:
			return fmt.Errorf("unknown flag %q", args[i])
		}
	}
	is, err := fetchIssue(num, true)
	if err != nil {
		return err
	}
	order, byID := logRows(is)
	if next {
		id, err := nextID(order, byID)
		if err != nil {
			return err
		}
		fmt.Println(id)
		return nil
	}
	if list {
		for _, id := range order {
			r := byID[id]
			gist := strings.TrimSpace(strings.TrimPrefix(firstLine(r.Text), "**"+id+"**"))
			gist = strings.TrimPrefix(gist, "- ")
			if len(gist) > 110 {
				gist = gist[:110] + "..."
			}
			tag := ""
			if isGlobal(r) {
				tag = " [global]"
			}
			fmt.Printf("%s%s - %s\n", id, tag, gist)
		}
		return nil
	}
	var missing []string
	printed := 0
	want := map[string]bool{}
	for _, id := range ids {
		want[id] = true
		if _, ok := byID[id]; !ok {
			missing = append(missing, id)
		}
	}
	for _, id := range order {
		r := byID[id]
		if want[id] || (global && isGlobal(r)) {
			fmt.Println(r.Text)
			fmt.Println()
			printed++
		}
	}
	if len(ids) == 0 && !global {
		for _, id := range order {
			fmt.Println(byID[id].Text)
			fmt.Println()
		}
		return nil
	}
	if len(missing) > 0 {
		return fmt.Errorf("not in log #%s: %s", num, strings.Join(missing, ", "))
	}
	if printed == 0 {
		fmt.Println("(no matching rows)")
	}
	return nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func cmdLogAppend(args []string) error {
	if len(args) < 2 {
		return errors.New("log-append <log-issue> <rows.md> [--amend]")
	}
	num, file := args[0], args[1]
	amend := len(args) > 2 && args[2] == "--amend"
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	text := normalize(string(data))
	if err := checkText(text); err != nil {
		return err
	}
	rows := parseRows(text)
	if len(rows) == 0 {
		return errors.New("no **{MAP-SLUG}-GM-NNN** rows found in file")
	}
	is, err := fetchIssue(num, true)
	if err != nil {
		return err
	}
	order, byID := logRows(is)
	if len(order) > 0 {
		prefix := strings.TrimSuffix(order[0], fmt.Sprintf("%03d", byID[order[0]].Num))
		for _, r := range rows {
			if !strings.HasPrefix(r.ID, prefix) {
				return fmt.Errorf("%s does not use this log's prefix %s", r.ID, prefix)
			}
		}
	}
	for _, r := range rows {
		_, exists := byID[r.ID]
		if exists && !amend {
			return fmt.Errorf("%s already exists in log #%s (use --amend to supersede it)", r.ID, num)
		}
		if !exists && amend {
			return fmt.Errorf("--amend given but %s is not in log #%s", r.ID, num)
		}
	}
	if !amend && len(order) > 0 {
		want, _ := nextID(order, byID)
		if rows[0].ID != want {
			return fmt.Errorf("first new row is %s; next free ID is %s", rows[0].ID, want)
		}
		for i := 1; i < len(rows); i++ {
			if rows[i].Num != rows[i-1].Num+1 {
				return fmt.Errorf("row IDs must be consecutive: %s follows %s", rows[i].ID, rows[i-1].ID)
			}
		}
	}
	if _, err := gh("issue", "comment", num, "--body-file", file); err != nil {
		return err
	}
	var posted []string
	for _, r := range rows {
		posted = append(posted, r.ID)
	}
	verb := "appended"
	if amend {
		verb = "amended"
	}
	fmt.Printf("OK: %s %s on log #%s\n", verb, strings.Join(posted, ", "), num)
	return nil
}

// --- map sections ---

// sectionBounds returns the line range [start, end) of a "## name" section,
// start being the heading line. ok is false when the heading is absent.
func sectionBounds(lines []string, name string) (start, end int, ok bool) {
	head := "## " + name
	start = -1
	for i, l := range lines {
		if strings.TrimSpace(l) == head {
			start = i
			break
		}
	}
	if start < 0 {
		return 0, 0, false
	}
	end = len(lines)
	for i := start + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") || strings.HasPrefix(lines[i], "# ") {
			end = i
			break
		}
	}
	return start, end, true
}

func cmdSection(args []string) error {
	if len(args) < 2 {
		return errors.New("section <issue> <Section name>...")
	}
	is, err := fetchIssue(args[0], false)
	if err != nil {
		return err
	}
	lines := strings.Split(is.Body, "\n")
	for _, name := range args[1:] {
		s, e, ok := sectionBounds(lines, name)
		if !ok {
			return fmt.Errorf("section %q not found on #%s", name, args[0])
		}
		fmt.Println(strings.TrimRight(strings.Join(lines[s:e], "\n"), "\n"))
		fmt.Println()
	}
	return nil
}

// --- map edits ---

type mapOps struct {
	complete string
	gist     string
	coverage []string
	appends  []string
	removes  []string
	dryRun   bool
}

func cmdMapEdit(args []string) error {
	if len(args) < 2 {
		return errors.New("map-edit <map-issue> [ops...]")
	}
	num := args[0]
	var ops mapOps
	for i := 1; i < len(args); i++ {
		val := func() (string, error) {
			if i+1 >= len(args) {
				return "", fmt.Errorf("%s needs a value", args[i])
			}
			i++
			return args[i], nil
		}
		var v string
		var err error
		switch args[i] {
		case "--dry-run":
			ops.dryRun = true
			continue
		case "--complete":
			v, err = val()
			ops.complete = strings.TrimPrefix(v, "#")
		case "--gist":
			v, err = val()
			ops.gist = v
		case "--coverage":
			v, err = val()
			ops.coverage = append(ops.coverage, v)
		case "--append":
			v, err = val()
			ops.appends = append(ops.appends, v)
		case "--remove":
			v, err = val()
			ops.removes = append(ops.removes, v)
		default:
			return fmt.Errorf("unknown flag %q", args[i])
		}
		if err != nil {
			return err
		}
	}
	if (ops.complete == "") != (ops.gist == "") {
		return errors.New("--complete and --gist go together")
	}
	for _, s := range append(append([]string{ops.gist}, ops.appends...), ops.coverage...) {
		if err := checkPunctuation(s); err != nil {
			return err
		}
	}
	is, err := fetchIssue(num, false)
	if err != nil {
		return err
	}
	body, notes, err := applyMapOps(is.Body, ops)
	if err != nil {
		return err
	}
	if err := validateMap(body); err != nil {
		return fmt.Errorf("edited map fails validation, not uploaded: %w", err)
	}
	if ops.dryRun {
		fmt.Println(strings.Join(notes, "\n"))
		fmt.Println("(dry run - nothing uploaded)")
		return nil
	}
	if err := uploadBody(num, body); err != nil {
		return err
	}
	notes = append(notes, fmt.Sprintf("OK: map #%s updated and verified", num))
	fmt.Println(strings.Join(notes, "\n"))
	return nil
}

var mdLink = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

// applyMapOps is the pure edit step, kept separate from gh for tests.
func applyMapOps(body string, ops mapOps) (string, []string, error) {
	lines := strings.Split(body, "\n")
	var notes []string

	if ops.complete != "" {
		ref := regexp.MustCompile(`(/issues/` + regexp.QuoteMeta(ops.complete) + `\)|\(#` + regexp.QuoteMeta(ops.complete) + `\))`)
		found := false
		for _, sec := range []string{"To Do", "Implementing"} {
			s, e, ok := sectionBounds(lines, sec)
			if !ok {
				continue
			}
			for i := s + 1; i < e; i++ {
				if !strings.HasPrefix(strings.TrimSpace(lines[i]), "|") || !ref.MatchString(lines[i]) {
					continue
				}
				link := mdLink.FindString(lines[i])
				if link == "" {
					return "", nil, fmt.Errorf("row for #%s in %s has no markdown link", ops.complete, sec)
				}
				lines = append(lines[:i], lines[i+1:]...)
				lines = appendToSection(lines, "Completed", fmt.Sprintf("- %s - %s", link, ops.gist))
				notes = append(notes, fmt.Sprintf("moved #%s: %s -> Completed", ops.complete, sec))
				found = true
				break
			}
			if found {
				break
			}
		}
		if !found {
			return "", nil, fmt.Errorf("no To Do or Implementing row links #%s", ops.complete)
		}
	}

	for _, c := range ops.coverage {
		parts := strings.SplitN(c, "=", 3)
		if len(parts) < 2 {
			return "", nil, fmt.Errorf("--coverage %q: want ID=status[=link]", c)
		}
		id, status := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		link := "-"
		if len(parts) == 3 && strings.TrimSpace(parts[2]) != "" {
			link = strings.TrimSpace(parts[2])
		}
		s, e, ok := sectionBounds(lines, "Decision coverage")
		if !ok {
			return "", nil, errors.New("map has no ## Decision coverage section")
		}
		row := fmt.Sprintf("| %s | %s | %s |", id, status, link)
		replaced := false
		for i := s + 1; i < e; i++ {
			cells := strings.Split(lines[i], "|")
			if len(cells) > 2 && strings.TrimSpace(cells[1]) == id {
				if len(parts) < 3 && len(cells) > 3 {
					if old := strings.TrimSpace(cells[3]); old != "" {
						row = fmt.Sprintf("| %s | %s | %s |", id, status, old)
					}
				}
				lines[i] = row
				replaced = true
				break
			}
		}
		if !replaced {
			lines = appendToSection(lines, "Decision coverage", row)
		}
		notes = append(notes, "coverage "+row)
	}

	for _, a := range ops.appends {
		sec, text, ok := strings.Cut(a, "=")
		if !ok {
			return "", nil, fmt.Errorf("--append %q: want Section=line", a)
		}
		text = strings.TrimSpace(text)
		if !strings.HasPrefix(text, "|") && !strings.HasPrefix(text, "- ") {
			text = "- " + text
		}
		if _, _, ok := sectionBounds(lines, strings.TrimSpace(sec)); !ok {
			return "", nil, fmt.Errorf("section %q not found", sec)
		}
		lines = appendToSection(lines, strings.TrimSpace(sec), text)
		notes = append(notes, fmt.Sprintf("appended to %s: %s", strings.TrimSpace(sec), text))
	}

	for _, r := range ops.removes {
		sec, sub, ok := strings.Cut(r, "=")
		if !ok || strings.TrimSpace(sub) == "" {
			return "", nil, fmt.Errorf("--remove %q: want Section=substring", r)
		}
		s, e, ok := sectionBounds(lines, strings.TrimSpace(sec))
		if !ok {
			return "", nil, fmt.Errorf("section %q not found", sec)
		}
		var hits []int
		for i := s + 1; i < e; i++ {
			if strings.Contains(lines[i], sub) {
				hits = append(hits, i)
			}
		}
		if len(hits) != 1 {
			return "", nil, fmt.Errorf("--remove %q matched %d lines in %s; make it unique", sub, len(hits), sec)
		}
		notes = append(notes, fmt.Sprintf("removed from %s: %s", strings.TrimSpace(sec), strings.TrimSpace(lines[hits[0]])))
		lines = append(lines[:hits[0]], lines[hits[0]+1:]...)
	}

	return strings.Join(lines, "\n"), notes, nil
}

// appendToSection inserts line after the last non-blank line of the section,
// so tables and lists stay contiguous and the blank line before the next
// heading is preserved.
func appendToSection(lines []string, name, line string) []string {
	s, e, _ := sectionBounds(lines, name)
	at := s + 1
	for i := e - 1; i > s; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			at = i + 1
			break
		}
	}
	if at == s+1 {
		// Empty section: keep one blank line under the heading.
		out := append([]string{}, lines[:s+1]...)
		out = append(out, "", line)
		return append(out, lines[s+1:]...)
	}
	out := append([]string{}, lines[:at]...)
	out = append(out, line)
	return append(out, lines[at:]...)
}

// --- body get / put ---

func cmdBody(args []string) error {
	if len(args) != 3 || (args[0] != "get" && args[0] != "put") {
		return errors.New("body get|put <issue> <file>")
	}
	num, file := args[1], args[2]
	if args[0] == "get" {
		is, err := fetchIssue(num, false)
		if err != nil {
			return err
		}
		if dir := filepath.Dir(file); dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
		}
		if err := os.WriteFile(file, []byte(is.Body), 0o644); err != nil {
			return err
		}
		fmt.Printf("OK: #%s body -> %s (%d lines, %d bytes)\n", num, file, strings.Count(is.Body, "\n")+1, len(is.Body))
		return nil
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	body := normalize(string(data))
	if err := checkText(body); err != nil {
		return err
	}
	if looksLikeMap(body) {
		if err := validateMap(body); err != nil {
			return fmt.Errorf("map body fails validation, not uploaded: %w", err)
		}
	}
	if err := uploadBody(num, body); err != nil {
		return err
	}
	fmt.Printf("OK: #%s body uploaded and verified (%d lines)\n", num, strings.Count(body, "\n")+1)
	return nil
}

// --- validation ---

func cmdValidate(args []string) error {
	if len(args) != 1 {
		return errors.New("validate <file>")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	body := normalize(string(data))
	if err := validateMap(body); err != nil {
		return err
	}
	fmt.Printf("OK: map body valid (%d lines, required sections present, ASCII punctuation)\n", strings.Count(body, "\n")+1)
	return nil
}

func looksLikeMap(body string) bool {
	return strings.Contains(body, "\n## To Do") || strings.Contains(firstLine(body), ":Map")
}

// validateMap catches collapsed bodies, missing sections, and encoding damage.
func validateMap(body string) error {
	if n := strings.Count(body, "\n") + 1; n < 40 {
		return fmt.Errorf("body has %d lines; collapsed maps are often ~10 lines", n)
	}
	lines := strings.Split(body, "\n")
	for _, sec := range []string{"To Do", "Completed", "Decision coverage"} {
		if _, _, ok := sectionBounds(lines, sec); !ok {
			return fmt.Errorf("missing section header on its own line: ## %s", sec)
		}
	}
	var heads []string
	for _, l := range lines {
		if strings.HasPrefix(l, "## ") {
			heads = append(heads, strings.TrimSpace(l))
		}
	}
	if len(heads) > 0 && heads[len(heads)-1] != "## Decision coverage" {
		return errors.New("## Decision coverage must be the last section")
	}
	return checkText(body)
}

func checkText(s string) error {
	if err := checkPunctuation(s); err != nil {
		return err
	}
	for _, r := range []string{"\u00c3", "\u0393", "\u256c", "\u251c", "\u2524", "\u252c", "\u2510", "\uFFFD"} {
		if strings.Contains(s, r) {
			return errors.New("mojibake sequences detected (encoding corruption)")
		}
	}
	return nil
}

func checkPunctuation(s string) error {
	for _, r := range []string{"\u2013", "\u2014", "\u00b7", "\u2022"} {
		if strings.Contains(s, r) {
			return errors.New("unicode dash or middle dot found - use ASCII hyphen only")
		}
	}
	return nil
}