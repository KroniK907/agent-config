# Example: Southwood Manor Mobile Home Park (AK2211677, 2025)

Source PDF: `AK2211677 Southwood Manor Mobile Home Park - 2025 CCR.pdf` (5 pages).

This example shows how the skill maps a **compliant small system** report. Your output for other PDFs will differ when violations or conditionals appear.

---

## Extraction highlights

| Item | Finding |
|------|---------|
| System | Southwood Manor Mobile Home Park |
| PWSID | AK2211677 |
| Period | 2025 (Annual Drinking Water Quality Report) |
| Source | Two groundwater wells; treatment at well #1 (arsenic, iron, H₂S); chlorine residual |
| Electronic delivery | Cover promotes online report: `https://ccrwriter.com/report/AK2211677` and QR - **§ 141.156(c)(1) likely applies** if primary delivery is electronic |
| **(d)(6)** | Detected Impurities tables: all listed rows `Violation` = **N** |
| Testing scope | “Water Quality Testing”: routinely monitored for applicable impurities; **only detected** impurities listed in Detected Impurities table |
| **(d)(8)** | Lead 2023: 4.8 ppb (90th), 0 samples > AL - **no AL exceedance** |
| **(f)** | No monitoring/reporting violation narratives |
| **(h)(6)** | No Groundwater Rule deficiency section |
| **(h)(7)** | No RTCR assessment / E. coli section |
| **(h)(2)** contact | System: (907) 344-0111, 9499 Brayton #68, Anchorage; Operator: Northern Utility Services (907) 222-4084, info@nusalaska.com |
| Subpart Q | Not stated |
| Biannual / (e) | Small mobile home park - not ≥ 10,000 persons |

**Non-violation notes (do not summarize as violations):**

- Arsenic detected (5.4 ppb) but **below MCL 10**; educational paragraph on low-level arsenic
- SOC / asbestos **testing waivers**
- Lead/copper educational text; corrosion responsibility split

---

## Sample summary output

Layout assumptions: the report template supplies the section **title** and **“About this report”** heading; the summary body starts with plain paragraphs.

**Headings.** Every subsection below the opening paragraphs uses `#####` (markdown h5). See **Heading rule** in `SKILL.md`.

The **paper copy** subsection does not repeat online/URL content (covered elsewhere on the printed report). The **§ 141.156(f)** sharing paragraph has no heading when the layout provides one. **Deliver to the user in a copyable markdown code block** (see `SKILL.md` Deliverables).

```markdown
This report summarizes the quality of your drinking water for 2025. It describes where your water comes from and how our monitoring results compare to federal drinking water standards.

We test your water for many contaminants required by federal and state rules. The tables in this report list only contaminants that were **detected**. Contaminants we tested for but did not find are **not shown here**, even though those tests were completed.

##### Your water at a glance
Our 2025 monitoring results show that your water **met all drinking water standards**. We did not identify any violations during the period covered by this report. Several contaminants were detected, including arsenic and chlorine, but all were **below the allowable limits** listed in this report. We were not required to run some tests (such as asbestos and certain synthetic organic chemicals) because they do not apply to our system. Skipping those tests is **not** a violation.

##### How to get more information
**Southwood Manor Mobile Home Park** - Phone: (907) 344-0111 - Address: 9499 Brayton #68, Anchorage, AK 99507  

**Operator: Northern Utility Services** - Phone: (907) 222-4084 (office 8:00-5:00 Mon-Fri; 24-hour emergency answering service) - Email: info@nusalaska.com  

You may contact us to request lead tap sampling results or our service line inventory.

##### Paper copy of this report
To request a **paper copy** of this report, call Southwood Manor Mobile Home Park at (907) 344-0111 or Northern Utility Services at (907) 222-4084. You can also email info@nusalaska.com.

Please share this information with anyone who drinks this water (or their guardians), especially those who may not have received this report directly (for example, people in apartments, nursing homes, schools, and businesses). You can do this by posting this report in a public place or distributing copies by hand, mail, email, or another method.
```

---

## Items intentionally omitted (this system)

| Section | Reason |
|---------|--------|
| § 141.156(c)(2) | User scope: no LEP translation block |
| § 141.156(c)(3) | Report does not state subpart Q public notification use |
| § 141.156(e) | Not a ≥10,000-person biannual system |

---

## If this report had violations (plain-language contrast)

**Example (d)(6):** “**Chlorine levels exceeded the allowable limit** during our monitoring on [dates]. The annual average was [X] mg/L; the limit is 4 mg/L. Some customers may experience stomach discomfort. **We reduced chlorine feed rates**, and levels are back within the limit. See the Detected Impurities table in this report.”

**Example (d)(8):** “**Lead levels in our tap samples exceeded the action level** ([X] µg/L; action level is 15 µg/L). **Run your tap** before using water for drinking or cooking. Use **cold water** for drinking, cooking, and preparing infant formula. Contact us if you would like your home water tested. See the Lead section in this report.”

**Example (h)(7):** “We detected **coliform bacteria** in our monitoring, which required a **Level 2 assessment** of the system. We completed [N] assessments and [M] corrective actions. See the coliform section in this report.”
