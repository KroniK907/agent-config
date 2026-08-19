# CCR Summary - Regulatory Requirements Reference

Authority: **40 CFR § 141.156** (summary), with content drawn from **§ 141.153** and conditionals from **§ 141.155**.

---

## § 141.156 - Complete structure

### (a) - Placement and purpose

- Summary at the **beginning** of the report, **prominently** displayed.
- Include a **brief description of the nature of the report** (Consumer Confidence Report / drinking water quality report for the prior calendar year or applicable monitoring period).
- State that the system **tests for many contaminants** but the report **lists only detected results** - impurities tested for but not found are **omitted** from the tables (aligns with typical “Detected Impurities” / “Water Quality Testing” body text; see `SKILL.md` **Detected vs. tested impurities**).

### (b) - Always required in the summary

#### (b)(1) - Violations and compliance (from report body)

Summarize information **included in the report** per:

| Citation | Report content to summarize |
|----------|----------------------------|
| **§ 141.153(d)(6)** | MCL, MRDL, and treatment technique violations in detected contaminant data: identification, plain-language explanation (length, health effects per appendix A, actions taken). |
| **§ 141.153(d)(8)** | Lead **action level** exceedance at § 141.80(c): exceedance identification, explanation, steps consumers can take to reduce exposure, corrective actions taken/planned. |
| **§ 141.153(f)** | Other **NPDWR compliance** violations during the report period, with explanation, health effects, and correction steps. Categories in the rule: (1) monitoring/reporting; (2) filtration/disinfection subpart H; (3) lead/copper **control requirements** subpart I (failure to take required actions - appendix A language); (4) Acrylamide/Epichlorohydrin subpart K; (5) recordkeeping; (6) special monitoring §§ 141.40, 141.41; (7) variance, exemption, or order violations. |
| **§ 141.153(h)(6)** | Groundwater Rule (subpart S): uncorrected significant deficiency and/or fecal indicator-positive groundwater source (nature, dates, correction plan/progress, health effects language). |
| **§ 141.153(h)(7)** | Revised Total Coliform Rule (subpart Y): assessment/corrective-action triggers, Level 1/2 assessments, E. coli MCL-related statements, failure to complete assessments or fix sanitary defects. |

#### (b)(2) - Contact

- Telephone (and ideally name/role) of **owner, operator, or designee** per **§ 141.153(h)(2)**.

### (c) - Include only if applicable

| Citation | Trigger | Summary must include |
|----------|---------|----------------------|
| **(c)(1)** | Delivery under **§ 141.155(a)(1)(ii)** (mail notice + website link), **(iii)** (email link/report), or **(iv)** (primacy-approved method) | Directions to **request a paper copy** per **§ 141.155(a)(2)** (prominently displayed in notification; repeat in summary). |
| **(c)(2)** | **§ 141.153(h)(3)** - primacy determines large proportion of consumers with limited English proficiency | Where to get translated report or language assistance. **Omit** when user confirms systems are below LEP thresholds. |
| **(c)(3)** | Report also satisfies **subpart Q** public notification | State dual purpose; brief nature of notice(s); how to locate notices in the report. |

### (d) - Style

- **Plain language at roughly an 8th grade reading level** (clear sentences, plain words; standard terms like *standards*, *violations*, and *contaminants* are fine). See skill `SKILL.md` for vocabulary rules and “Your water at a glance” enforcement.
- Infographics encouraged.

### (e) - Biannual large systems

- Systems **≥ 10,000** persons with **§ 141.155(j)(2)** biannual delivery: describe report + **6-month update** and new information Jan-Jun.

### (f) - Mandatory verbatim text

```
Please share this information with anyone who drinks this water (or their guardians), especially those who may not have received this report directly (for example, people in apartments, nursing homes, schools, and businesses). You can do this by posting this report in a public place or distributing copies by hand, mail, email, or another method.
```

---

## PDF search patterns

Use case-insensitive search on extracted text:

**Violations / exceedances**

- `violation`, `exceeded`, `failure to`, `MCL violation`, `MRDL`, `treatment technique`
- Table columns: `Violation`, `Samples > AL`, `Yes`, `Y`
- Lead: `action level`, `90th percentile`, `exceedance`

**§ 141.153(f)-type narratives**

- `monitoring violation`, `reporting violation`, `failed to take`, `assessment`, `Level 1`, `Level 2`, `significant deficiency`, `fecal indicator`, `E. coli`, `variance`, `exemption`

**Contacts**

- `contact`, `phone`, `operator`, `owner`, `tel`, `@`, PWSID block at end

**Electronic delivery / paper copy**

- `online`, `website`, `QR`, `email`, `request a paper`, `paper copy`

**Public notification (subpart Q)**

- `public notification`, `notice of violation`, `Tier 1`, `Tier 2`, `subpart Q`

**Not violations (common false positives)**

- `waiver`, `below the MCL`, `no violation`, `Violation N`, educational `arsenic`/`lead` when tables show compliance

---

## § 141.153(d)(6) vs (d)(8) vs (f) - Quick distinction

| Issue | Section |
|-------|---------|
| Coliform/nitrate/chlorine/etc. over **MCL/MRDL/TT** in data tables | **(d)(6)** |
| Lead **90th percentile / AL exceedance** in lead/copper table | **(d)(8)** |
| Failed to monitor, report, assess, install treatment, keep records, etc. | **(f)** |
| Groundwater deficiency / fecal positive source | **(h)(6)** |
| RTCR assessments / E. coli MCL narrative | **(h)(7)** |

---

## § 141.155 - Delivery (for (c)(1))

Direct delivery options (§ 141.155(a)(1)):

1. Mail or hand-deliver **paper** report
2. Mail notification with **website direct link**
3. **Email** direct link or electronic report
4. Other primacy-approved direct method

Electronic methods (2-4) require **paper copy on request**; directions must be **prominent** in the notification and summarized per § 141.156(c)(1).

---

## Agent output rules

1. **Never fabricate** violations, dates, or contacts not supported by the PDF.
2. **Distinguish** “detected but compliant” from “violation.”
3. **Quote** § 141.156(f) sharing language exactly.
4. **Ask** the user if delivery method or subpart Q dual-use is unclear from the PDF alone.
5. When tables use `N` / `No` for violation columns, treat as **no violation** for that row unless narrative contradicts.
6. **Headings** - all summary subsections use `#####` (h5) only. Opening paragraphs have no heading. Full rule in `SKILL.md`.

**Done when:** summary passes the Step 7 compliance self-check in `SKILL.md`.
