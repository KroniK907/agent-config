---
name: ccr-summary
description: CCR summary, summarize a CCR, Consumer Confidence Report PDF, water quality report PDF, 40 CFR 141.156, § 141.156 summary
disable-model-invocation: true
---

# CCR Summary (40 CFR § 141.156)

Generate the **prominent opening summary** for a CCR by reading the full report PDF and mapping findings to [40 CFR § 141.156](https://www.law.cornell.edu/cfr/text/40/141.156). The summary is **not** a rewrite of the whole report - it distills what consumers need first.

**Scope note:** This workflow **excludes § 141.156(c)(2)** (limited English proficiency / translated reports). Only include that block if the user explicitly asks and the system is subject to § 141.153(h)(3).

**Local references** (in this repo): `40 CFR § 141.153…html`, `40 CFR § 141.155…html`, `40 CFR § 141.156…html`.

For the full regulatory checklist and PDF search patterns, read [requirements.md](requirements.md). For a worked small-system example, read [example.md](example.md).

---

## Heading rule (do not skip)

**All subsection headings in summary output MUST be h5 - five hash marks (`#####`) in markdown.**

| ✓ Correct | ✗ Wrong |
|-----------|---------|
| `##### Your water at a glance` | `## Your water at a glance` |
| `##### How to get more information` | `### How to get more information` |
| `##### Paper copy of this report` | `#### Paper copy of this report` |

- **Opening paragraphs:** no heading (the printed report supplies the title).
- **Subsections below the opening:** `#####` only - never `#`, `##`, `###`, or `####`.
- **§ 141.156(f) sharing block:** omit the heading if the layout supplies it; if you include one, it must still be `##### Please share this report`.

This matches the CCR layout template. Using `##` or other levels will break formatting.

---

## Workflow

Copy this checklist and mark items as you go:

```
CCR summary progress:
- [ ] Step 1: Extract full PDF text
- [ ] Step 2: Identify system metadata and delivery context
- [ ] Step 3: Scan for violations/compliance (required sources)
- [ ] Step 4: Extract contact information
- [ ] Step 5: Evaluate conditional § 141.156(c) and (e) items
- [ ] Step 6: Draft summary using template (**subsection headings = `#####` / h5 only**)
- [ ] Step 7: Compliance self-check (**verify h5 headings before sending**)
```

### Step 1: Extract full PDF text

Do not rely on partial reads. Extract **every page** to a UTF-8 text file (avoids Windows console encoding errors):

```python
from pypdf import PdfReader

path = r"<absolute-path-to-ccr.pdf>"
out = r"<absolute-path-to-extract.txt>"
reader = PdfReader(path)
with open(out, "w", encoding="utf-8") as f:
    for i, page in enumerate(reader.pages):
        f.write(f"--- PAGE {i+1} ---\n")
        f.write(page.extract_text() or "")
        f.write("\n\n")
```

If `pypdf` is missing: `pip install pypdf`. If extraction is empty or garbled, tell the user the PDF may be scanned - OCR will be needed before summarizing.

### Step 2: System metadata and delivery context

Record from the report (usually cover, footer, or “System Contact”):

| Field | Where to look |
|-------|----------------|
| System name | Title / header |
| PWSID | Header or contact block |
| Report year / monitoring period | Title (“For The Year Of …”) |
| Population served | Sometimes in header or certification (not always in PDF) |
| Delivery method | Cover “online version,” QR code, website URL, mail/email notices |

**Delivery drives conditionals:**

- **Paper mail/hand delivery only** → § 141.156(c)(1) usually **not** needed.
- **Website link, email link, or other electronic delivery** (§ 141.155(a)(1)(ii)-(iv)) → summary **must** include directions to request a **paper copy** per § 141.155(a)(2). Use the system’s actual request process if stated; otherwise use the regulatory minimum (contact phone/email/mail).

**Biannual / 6-month update (§ 141.156(e)):** Required only if the system serves **≥ 10,000 persons** and distributes twice per year (§ 141.155(j)(2)). Small systems like mobile home parks typically **skip (e)** unless the report itself states biannual delivery.

### Step 3: Violations and compliance summary (§ 141.156(b)(1))

Summarize **only what appears in the report** for the monitoring period covered. Source sections in the **full report** (not only the summary you are writing):

| CFR source | What to find in the PDF |
|------------|-------------------------|
| **§ 141.153(d)(6)** | Contaminant tables: `Violation` = Y/Yes; values over MCL/MRDL/TT; narrative “violation,” “exceeded,” “failed to meet.” Include length, health effects language, and corrective actions if provided. |
| **§ 141.153(d)(8)** | Lead **action level** exceedance: 90th percentile or `Samples > AL` > 0; lead exceedance explanation; consumer steps to reduce exposure; corrective actions. |
| **§ 141.153(f)** | NPDWR compliance violations **other than** contaminant-level exceedances already in (d)(6): monitoring/reporting failures, filtration/disinfection (subpart H), lead/copper **rule** violations (failure to complete required actions - not the same as AL exceedance in (d)(8)), subpart K, recordkeeping, §§ 141.40/141.41 monitoring, variance/exemption/order violations. |
| **§ 141.153(h)(6)** | Groundwater Rule (subpart S): uncorrected **significant deficiency**, fecal indicator-positive source, state-approved correction plan - often a dedicated section. |
| **§ 141.153(h)(7)** | Revised Total Coliform Rule (subpart Y): Level 1/2 assessments, E. coli MCL, missed assessments/corrective actions - often a dedicated section. |

**If no violations in any category:** State clearly that the report shows **no violations** of drinking water standards for the period covered (still required content for (b)(1)).

**Do not treat as violations:**

- Educational arsenic/lead text when MCL/AL are met
- Monitoring **waivers** (e.g., SOC, asbestos) unless the report says a violation occurred
- Detected contaminants **below** MCL/AL with `Violation` = N

**Testing waivers (optional one-sentence summary):** If the report has a “Testing Waivers” (or similar) section and the system had **no violations**, you may add **one short plain-language sentence** in **Your water at a glance** so readers know skipped tests are not problems. **Do not** name EPA, the state, primacy agency, or any regulator as who granted the waiver - the report may say “EPA,” “State,” or “ADEC,” but the summary must stay **agency-neutral**.

| Avoid | Use instead |
|-------|-------------|
| “when EPA allowed us to skip tests” | “we were not required to run some tests” |
| “EPA granted a waiver” | “some tests do not apply to our system” |
| “ADEC/the State approved…” | (omit agency name entirely) |

**Example (plain language):** “We also did not run some tests (for example, asbestos and certain chemicals) because they do not apply to our system. Skipping those tests is **not** a violation.”

See [requirements.md](requirements.md) for subsection detail and search keywords.

### Step 4: Contact information (§ 141.156(b)(2) / § 141.153(h)(2))

Include **telephone** (required) for owner, operator, or designee. Also include name, email, and address when present. Prefer the block labeled “System Contact,” “Operator Contact,” or end-of-report contact.

### Step 5: Conditional blocks

| Section | Include when |
|---------|----------------|
| **§ 141.156(c)(1)** | Electronic delivery per § 141.155(a)(1)(ii)-(iv) |
| **§ 141.156(c)(2)** | **Skip** unless user requests (large LEP population per § 141.153(h)(3)) |
| **§ 141.156(c)(3)** | Report explicitly states it **also serves as public notification** under 40 CFR part 141 **subpart Q** |
| **§ 141.156(e)** | System ≥ 10,000 persons with biannual report + 6-month update |

### Step 6: Draft the summary

Keep it **short and prominent** - consumers should understand the report’s purpose, compliance status, and who to call without reading all tables.

**Detected vs. tested impurities (§ 141.156(a)) - REQUIRED:** Many CCR templates include a “Water Quality Testing” or “Detected Impurities” section explaining that only impurities with a measurable result appear in the tables below. The opening summary **must** state this plainly so readers do not think an unlisted impurity was never tested.

Include **all three ideas** (1-2 short sentences in the opening paragraphs - usually the second or third paragraph):

1. We **test for many** contaminants/impurities required by law.
2. This report **lists only** impurities where we got a result.
3. Impurities we tested for but **did not detect** are **not included** in the report - even though we ran those tests.

**Scan the PDF for:** `Water Quality Testing`, `Detected Impurities`, `only those detected`, `routinely monitored`. Use the report’s wording when helpful, but keep the summary at ~8th grade (see Plain language rules).

**Example (opening paragraph):** “We test your water for many contaminants required by federal and state rules. The tables in this report list only contaminants that were **detected**. Contaminants we tested for but did not find are **not shown here**, even though those tests were completed.”

**Do not** bury this only in **Your water at a glance** - it belongs in the § 141.156(a) opening block about the nature of the report.

## Plain language (§ 141.156(d)) - REQUIRED

**Target reading level: 8th grade.** This is not optional. § 141.156(d) requires plain language - clear and direct, but not oversimplified. Write so a typical middle-school reader can follow it without water-industry training. **Heavily enforce** readability in **Your water at a glance** - that section must be the clearest part of the summary.

### Rules (apply to all summary text)

| Do | Don't |
|----|-------|
| Clear sentences (aim for **≤ 20 words**; rarely exceed 25) | Long, nested sentences or dense legal phrasing |
| Plain, everyday language; standard terms readers expect in a water report (**standards**, **violations**, **contaminants**, **monitoring**) | Regulatory jargon, CFR citations, or acronym stacks (MCL, MRDL, NPDWR, RTCR, AL, TT) |
| Active voice (“We found…”, “Your water met…”) | Passive voice (“Monitoring was conducted…”) |
| “We / your water / this report” | “The system,” “consumers,” “persons served” (unless a proper name) |
| One main idea per sentence | Packing unrelated facts into one sentence |
| Brief explanation when a technical term is needed | Copying table headers or appendix language verbatim |

**Acronyms:** Avoid in the summary when possible. Prefer plain words or short phrases consumers will understand:

| Instead of | Write |
|------------|-------|
| MCL / maximum contaminant level | **allowable limit** / **drinking water standard** |
| MRDL | **limit on disinfectant** (e.g., chlorine) |
| violation | **violation** / **did not meet standards** |
| compliance | **met standards** / **in compliance** |
| action level (lead/copper) | **action level** (briefly: level that triggers follow-up steps) |
| treatment technique | **required treatment** |
| disinfection byproducts (TTHMs, HAA5) | **disinfection byproducts** (may add: formed when water is treated with chlorine) |
| contaminants / impurities | **contaminants** |
| monitoring / sampling | **monitoring** / **water tests** |
| corrective actions | **corrective actions** / **steps we took to fix the problem** |
| detected | **detected** / **found** |
| exceed / exceedance | **exceeded the limit** |
| 90th percentile | state the result plainly (e.g., lead level in tap samples); skip the statistic name unless the report emphasizes it |
| fecal indicator-positive | **signs of contamination in the source** (only if the report says this) |
| significant deficiency | **significant deficiency in the system** (plain follow-up if the report explains it) |

### Your water at a glance - extra strict

This subsection summarizes § 141.156(b)(1). It must be **easier to read than the rest of the report**, not harder.

1. **Lead with the bottom line** in the first sentence: Did water meet standards or were there violations?
2. **No violations:** State clearly that monitoring shows water **met drinking water standards** for the period covered. Do not list every contaminant unless it helps the reader. One short paragraph is often enough. If the report describes **testing waivers**, you may add one agency-neutral sentence (see Step 3) - never attribute waivers to EPA or any named agency.
3. **Violations:** Use **short bullets**. Each bullet: what happened → duration (if known) → possible health effects (plain words from the report) → corrective actions. No appendix or CFR references.
4. **Never** copy table headers or regulatory phrases into this section - **translate** them into plain language at an 8th-grade level.
5. Read the draft aloud. If it sounds like a legal brief or engineering memo, simplify - but do not talk down to the reader.

**Bad (too complex):**
> Based on monitoring results and compliance information, no violations of maximum contaminant levels, disinfectant residual levels, treatment techniques, or lead action levels were identified.

**Bad (too simplified):**
> Our water tests show your water met all safety rules. We did not find any problems that broke EPA limits.

**Good (8th grade, plain language):**
> Our 2025 monitoring results show that your water **met all drinking water standards**. We did not identify any violations during the period covered by this report.

**Formatting (user layout) - h5 headings are mandatory:**

- **Do not** include a document title or an “About this report” heading - the printed report already provides those.
- **Begin** with 1-3 plain paragraphs (no heading) describing the nature of the report (§ 141.156(a)). **Always** explain that we test for many contaminants but only list those detected (see **Detected vs. tested impurities** above).
- **Subsection headings = h5 only.** Write `##### Your water at a glance`, not `## Your water at a glance`. See **Heading rule** at the top of this skill. **Never** use `#`, `##`, `###`, or `####` in summary output - agents often default to `##`; override that habit every time.
- **Paper copy (§ 141.156(c)(1)):** State only how to **request a paper copy** (phone, mail, email). **Do not** mention the online report, website URL, QR code, or “online version” - that information is already elsewhere on printed reports.

**Output template:**

```markdown
[Opening paragraphs only - no heading. Plain language at ~8th grade: annual water quality report; source; how results compare to drinking water standards; monitoring period. **Required:** we test for many contaminants; tables list only detected results; undetected impurities are omitted even though we tested for them.]

##### Your water at a glance
[§ 141.156(b)(1). **Clearest section in the whole summary.** Bottom line first. Plain language at ~8th grade - see rules above. Short bullets if there are violations; one short paragraph if none.]

##### How to get more information
[§ 141.156(b)(2): who to call. Short lines: name, phone, address, email. Plain words for anything extra (e.g., “Ask us for lead test results.”).]

##### [If applicable] Paper copy of this report
[§ 141.156(c)(1): how to request a paper copy - contact method only; no online/website/QR references.]

##### [If applicable] Public notification
[§ 141.156(c)(3): state report doubles as public notice; brief nature of notice(s); where to find detail in report.]

##### [If applicable] Mid-year update
[§ 141.156(e): only for biannual large systems.]

##### Please share this report
Please share this information with anyone who drinks this water (or their guardians), especially those who may not have received this report directly (for example, people in apartments, nursing homes, schools, and businesses). You can do this by posting this report in a public place or distributing copies by hand, mail, email, or another method.
```

**§ 141.156(f):** The “Please share this report” heading may be omitted if the layout supplies it; use `##### Please share this report` when a heading is included. The paragraph text must still appear **verbatim** as the final content.

### Step 7: Compliance self-check

Before delivering:

- [ ] Brief description of report nature (§ 141.156(a)), including **tested vs. listed** impurities (many tested; only detected shown; undetected omitted)
- [ ] (b)(1) All five source areas addressed or explicitly “none”
- [ ] (b)(2) Phone contact included
- [ ] (c) Only applicable subsections; (c)(2) omitted unless requested
- [ ] (f) Standard sharing language **exact**
- [ ] No invented violations or contacts
- [ ] **Plain language:** entire summary at ~8th grade level - clear, not oversimplified
- [ ] **Your water at a glance:** clearest section; no unexplained acronyms; readable sentences
- [ ] **Headings:** all subsection headings use `#####` (h5) only
- [ ] Read **Your water at a glance** aloud - it should be clear to someone without a water background
- [ ] **Output format:** finished summary placed in a **copyable markdown fenced code block** in chat (see Deliverables)
- [ ] Flag ambiguities (e.g., electronic delivery unclear, partial PDF extraction)

---

## Deliverables

1. **Finished summary** - deliver in chat inside a **single markdown fenced code block** so the user can copy it in one click:

   ````markdown
   [Opening paragraphs - no heading]

   ##### Your water at a glance
   ...

   ##### How to get more information
   ...
   ````

 - **Every subsection heading inside the block must be `#####` (h5).** Scan the draft for `##`, `###`, or `####` before sending - replace any you find.
 - Use a ` ```markdown ` fence (or ` ``` `) containing **only** the summary text - no commentary inside the block.
 - Put extraction notes, assumptions, and compliance flags **outside** the block (plain text after the fence).
 - Do **not** save to a file unless the user asks.
2. **Extraction notes** (optional, brief): page/section references for violations and contacts - **outside** the code block.
3. **Gaps/assumptions**: delivery method, population, missing sections - **outside** the code block.

---

## Additional resources

- [requirements.md](requirements.md) - full § 141.156/153/155 mapping
- [example.md](example.md) - Southwood Manor AK2211677 (2025) sample
