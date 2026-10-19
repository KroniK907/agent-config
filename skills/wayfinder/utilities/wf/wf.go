// wf - wayfinder map helper. Keeps large map bodies out of agent context:
// reads return only the sections asked for, and uploads go through gh with
// UTF-8 files and a verify fetch.
//
// Stdlib only, single file, so it runs from any project (even inside another
// Go module) without a go.mod:
//
//	go run <wayfinder-skill-dir>/utilities/wf/wf.go <command> [flags]
//
// Commands:
//
//	section <issue> <Section name>...
//	body get <issue> <file>
//	body put <issue> <file>
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
	case "section":
		err = cmdSection(args[1:])
	case "body":
		err = cmdBody(args[1:])
	default:
		err = fmt.Errorf("unknown command %q\n%s", args[0], usage)
	}
	if err != nil {
		fail(err.Error())
	}
}

const usage = `usage: wf [-R owner/repo] <command>
  section <issue> <Section name>...
  body get <issue> <file>
  body put <issue> <file>`

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "wf: "+msg)
	os.Exit(1)
}

// --- gh access ---

type issue struct {
	Body  string `json:"body"`
	Title string `json:"title"`
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

func fetchIssue(num string) (*issue, error) {
	out, err := gh("issue", "view", num, "--json", "body,title")
	if err != nil {
		return nil, err
	}
	var is issue
	if err := json.Unmarshal(out, &is); err != nil {
		return nil, err
	}
	is.Body = normalize(is.Body)
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
	back, err := fetchIssue(num)
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

// --- sections ---

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
	is, err := fetchIssue(args[0])
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

// --- body get / put ---

func cmdBody(args []string) error {
	if len(args) != 3 || (args[0] != "get" && args[0] != "put") {
		return errors.New("body get|put <issue> <file>")
	}
	num, file := args[1], args[2]
	if args[0] == "get" {
		is, err := fetchIssue(num)
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
	if err := validateBody(body); err != nil {
		return fmt.Errorf("not uploaded: %w", err)
	}
	if err := uploadBody(num, body); err != nil {
		return err
	}
	fmt.Printf("OK: #%s body uploaded and verified (%d lines)\n", num, strings.Count(body, "\n")+1)
	return nil
}

// validateBody catches encoding damage on any body, and collapsed or
// truncated map bodies.
func validateBody(body string) error {
	for _, r := range []string{"\u00c3", "\uFFFD"} {
		if strings.Contains(body, r) {
			return errors.New("mojibake sequences detected (encoding corruption)")
		}
	}
	if !strings.Contains(firstLine(body), ":Map") {
		return nil
	}
	lines := strings.Split(body, "\n")
	for _, sec := range []string{"Target outcome", "Decisions"} {
		if _, _, ok := sectionBounds(lines, sec); !ok {
			return fmt.Errorf("map body is missing its own-line heading: ## %s", sec)
		}
	}
	return nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
