# i18n — operator-commons multilingual policy

> **English** | [한국어](README.ko.md) | [日本語](README.ja.md)

This document defines how `operator-commons` documentation is
maintained in multiple languages. The project keeps an English canonical
plus Korean / Japanese / Chinese translations.

## §1 Policy

### 1.1 Core principles

- **Canonical language**: English (ground truth for every source).
- **Translated languages (three additional)**: Korean (`ko`), Japanese
  (`ja`), Chinese (`zh`, Simplified).
- **Four-language skeleton**: `README.md` (English) + `README.ko.md` +
  `README.ja.md` + `README.zh.md`.
- **No language beyond these four** — adding `es` / `fr` / `de` / …
  requires a separate ADR.

### 1.2 In scope vs. out of scope

**In scope** (user-facing or external documents):

- `README.{md, ko, ja, zh}.md`
- Top-level governance / branding: `BRANDING.md`, `STABILITY.md`,
  `coverage-report.md` and similar four-language matrices.
- `docs/i18n/glossary-{ko, ja, zh}.md` (this SSOT).

**Out of scope** (internal, governance, or legal):

- `docs/kb/adr/*.md` (decision records — English canonical only).
- `docs/kb/deps/*.md` (dependency audit — generated).
- `LICENSE`, `NOTICE`, `CITATION.cff`, `.gitignore` (legal or config).

## §2 Glossary

### 2.1 Location

SSOT: [`glossary-ko.md`](glossary-ko.md) / [`glossary-ja.md`](glossary-ja.md)
/ [`glossary-zh.md`](glossary-zh.md).

### 2.2 Priority

1. Glossary definitions take priority.
2. The official Kubernetes Korean / Japanese / Chinese guides.
3. General dictionaries.
4. **Code identifiers stay in English** (no translation) — `pkg/probes`,
   `Reconciler`, `kubectl`, etc.

### 2.3 Consistency rules (all four languages)

- Do not mix formal and informal register inside a single document.
- User-facing documents use formal register.
- Internal documents (AGENTS) may use informal register.
- On first occurrence, include the English original plus the
  parenthesised translation; subsequent occurrences may use the
  translation alone.

## §3 Automated translation SOP

### 3.1 Engine

**Default**: Claude direct (an AI subagent reads the source and
translates).

Rationale: Korean quality and prompt caching cost. Japanese and
Chinese translations are followed by native review.

### 3.2 Commands

```bash
# Manual single-file translation
# The subagent reads the source, translates, and writes the output
# with the AI-translation warning banner.

# Automation (scripted)
./scripts/i18n-translate.sh <source.md> --lang all --engine claude

# Dry run
./scripts/i18n-translate.sh README.md --dry-run
```

### 3.3 Result marker

Every AI-translated file must include the warning banner:

```markdown
> ⚠️ This translation is AI-generated and pending native review.
```

Additional markers: `[needs review]` (review pending), `[reviewed]`
(promoted by a native reviewer through a PR).

## §4 Native review SOP

### 4.1 Reviewer recruitment

Default policy: AI translation + `[needs review]` marker + warning
banner. The placeholder remains until a native reviewer is recruited.

### 4.2 Promotion criteria

When promoting `[needs review]` → `[reviewed]`, the native reviewer
confirms:

1. **Idiomatic vs. literal balance** — natural reading first, but
   technical terminology stays precise.
2. **Glossary consistency** — every glossary term is applied
   identically.
3. **Register consistency** — formal / informal register is unified
   within a document.
4. **Link integrity** — cross-links resolve per language.
5. **Code identifiers untranslated**.

### 4.3 Promotion PR format

```
title: docs(i18n): <lang> <doc> native reviewer promotion
body:
  ## Reviewer
  - <name / handle>

  ## Scope reviewed
  - <file path>

  ## Changes
  - [needs review] → [reviewed] marker change
  - Warning banner removed (or kept as partial)
  - <other idiomatic edits>

  ## Verification
  - [x] Glossary consistency
  - [x] Register consistency
  - [x] Cross-link integrity
```

## §5 Drift control

### 5.1 lefthook hook

`lefthook.yml` enforces `readme-i18n-sync` at pre-push:

```yaml
readme-i18n-sync:
  glob: "README*.md"
  run: bash scripts/check-readme-sync.sh
```

The pre-push position keeps work-flow friction low.

### 5.2 Thresholds

| Lang | Line-diff threshold | Reason |
|---|---|---|
| ko | ≤ 15 % | Korean is roughly the same LOC as English. |
| ja | ≤ 25 % | Japanese mixes kana and kanji — ~15–20 % shorter than English. |
| zh | ≤ 30 % | Chinese is fully Hanzi — ~25 % shorter than English. |

These are recommendations; adjust after empirical measurement.

### 5.3 4-lang matrix

`scripts/check-readme-sync.sh` automatically checks the four-language
matrix (EN ↔ {ko, ja, zh}):

- Target-language file absent → skip (handles partial skeletons).
- Per-language bypass available
  (`SKIP_CHECK_README_SYNC_JA=1` and similar).

## §6 References

- Glossaries: [`glossary-ko.md`](glossary-ko.md) /
  [`glossary-ja.md`](glossary-ja.md) / [`glossary-zh.md`](glossary-zh.md).
- Check script: [`../../scripts/check-readme-sync.sh`](../../scripts/check-readme-sync.sh).
- Translation script: [`../../scripts/i18n-translate.sh`](../../scripts/i18n-translate.sh).
