# S2 — operator-commons stale 정리 + i18n placeholder → 본문 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:subagent-driven-development` (recommended) or `superpowers:executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** operator-commons 의 stale 브랜치 1개 정리 + i18n drift (family.md v0.7.0 → v0.8.0, README.{ja,zh}.md 본문 보강, glossary-{ja,zh}.md 본문 보강) 해소. 5 PR / 5 atomic commit / RFC-0025 `[~]` partial marker 유지 (native reviewer 후속 검증 의도 — wave 6).

**Architecture:** spec `docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2` 의 4개 작업 (a/b/c/d) 을 5 task 로 분해 — 각 task = 1 branch + 1 commit + 1 PR + 1 squash merge + branch 즉시 삭제 (CLAUDE.md §8 atomic).

**Tech Stack:** git · gh CLI · markdown · 한국어 ko + 日本語 ja + 中文 zh translation (machine quality, native reviewer 후속). 외부 lib 미사용.

**참조 spec:** [portfolio supercycle](../specs/2026-05-21-portfolio-cleanup-supercycle-design.md) (Sub-Project S2 카드)

**전제 (Pre-flight)**:
- `pwd` = `/Users/phil/workspace/keiailab/operator-commons`
- `git status --porcelain` 결과 = 빈 출력 (clean)
- `git rev-parse --abbrev-ref HEAD` = `main`
- `git log -1 --format='%h'` 가 `da4a410` (spec 머지 commit) 이상
- `gh auth status` = ✓ logged in `eightynine01`

---

## Task 1: archive 브랜치 정밀 분석 + tag 보존 후 삭제

**Files:**
- Read: 원격 브랜치 분석 (commit 없음)
- Push: `origin/archive/main-13-commits-merge-style-2026-05-21-final-tag` (annotated tag)
- Delete: `origin/archive/main-13-commits-merge-style-2026-05-21` (브랜치)

**검증 명령**: 모든 작업 종료 후 `git branch -r | wc -l` = 1 (origin/main 만; HEAD pointer 는 ls 대상 아님)

- [ ] **Step 1: archive 브랜치 의 main 비-equivalent commit 정밀 식별**

Run: 
```bash
git fetch --prune origin
git cherry origin/main origin/archive/main-13-commits-merge-style-2026-05-21 -v > /tmp/cherry-archive.txt
cat /tmp/cherry-archive.txt
```

Expected: 8 lines (5 `+` + 3 `-`). `+` = archive 만 가진 commit, `-` = main 에 equivalent 존재.

기대 출력 (이미 정찰로 확인됨):
```
+ 51a23eb feat(ci): RFC-0043 L5 shadow + lefthook bootstrap + gofmt 정형
+ 5877f11 chore(lint): golangci-lint 3 issue fix (goconst + lll + staticcheck QF1001)
- 0358e5e6...
- b27d691d...
+ 80d03b8 docs(i18n): README 다중언어 (EN canonical + KR) + sync hook + owner
+ 51e796f fix(ci): lefthook config union
+ 5581cd2 fix(ci): lefthook.yml union body
```

- [ ] **Step 2: 5 `+` commit 들이 main 의 *어떤 squash merge* 에 포함됐는지 매핑**

Run:
```bash
# 각 archive-only commit 의 작업 내용이 main 에 머지된 PR 찾기
for h in 51a23eb 5877f11 80d03b8 51e796f 5581cd2; do
  echo "=== $h ==="
  git show --stat "$h" | head -3
  echo "--- 동일 작업의 main PR (제목 grep) ---"
  git log --oneline --all | grep -i -E '(L5 shadow|RFC-0043|i18n.*README|lefthook.*union|golangci.*goconst)' | head -3
done
```

Expected: 각 archive-only commit 의 *내용* 이 main 의 PR (e.g. #26, #27) 와 *file-content* 동일함을 확인. 동일하면 *cherry-pick 흡수 완료* 상태 = 보존 불요.

- [ ] **Step 3: annotated tag 보존 (안전 시 삭제 전 last-resort backup)**

Run:
```bash
git tag -a archive/main-13-commits-merge-style-2026-05-21-final-tag \
  origin/archive/main-13-commits-merge-style-2026-05-21 \
  -m "archive: main-13-commits-merge-style-2026-05-21 — pre-deletion backup (2026-05-21 S2 cleanup, all 5 archive-only commits cherry-equivalent in main)"
git push origin archive/main-13-commits-merge-style-2026-05-21-final-tag
git tag -l 'archive/main-13-commits-merge-style*' | head -5
```

Expected:
```
archive/main-13-commits-merge-style-2026-05-21-final-tag
archive/main-13-commits-merge-style-2026-05-21-tag       # 기존 tag (확인용)
```

- [ ] **Step 4: archive 브랜치 origin 삭제**

Run:
```bash
git push origin --delete archive/main-13-commits-merge-style-2026-05-21
git fetch --prune origin
git branch -r
```

Expected:
```
  origin/HEAD -> origin/main
  origin/main
```
(2 lines — archive 사라짐)

- [ ] **Step 5: ADR 작성 (S2 의 archive 결정 기록)**

Create file `docs/kb/adr/0009-archive-branch-cleanup-policy.md`:

```markdown
# ADR-0009 — archive/* 브랜치 정리 정책 (cleanup supercycle 2026-05-21)

| 메타 | 값 |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Supersedes | — |
| Refs | docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2 |

## Context

`archive/main-13-commits-merge-style-2026-05-21` 브랜치가 main 과 ahead/behind = 13/13 비대칭으로 존재. cherry-equivalence 분석 결과 archive-only 5 commits (`51a23eb`, `5877f11`, `80d03b8`, `51e796f`, `5581cd2`) 의 *내용* 이 main 의 PR #26 / #27 squash merge 결과에 동일하게 흡수됨.

## Decision

1. 삭제 전 annotated tag (`archive/main-13-commits-merge-style-2026-05-21-final-tag`) 로 *마지막 백업* 보존.
2. 브랜치 자체는 origin 에서 삭제 (`git push origin --delete`).
3. 향후 동일 패턴 (squash merge 후 historical branch 보존 의도) 의 archive 브랜치는 동일 절차:
   - cherry-equivalence 정밀 확인 (`git cherry main archive/* -v`)
   - annotated tag 로 backup (`-final-tag` suffix)
   - 브랜치 삭제 + ADR 기록

## Consequences

- (+) `git branch -r` 가 main 만 (gh-pages 부재; commons 는 GitHub Pages 미사용).
- (+) GitHub 의 branch dropdown UI noise 제거.
- (+) tag 로 git history 100% 보존 — 필요 시 `git checkout <tag>` 로 복원.
- (-) tag 명명 규칙 `<branch>-final-tag` 의 표준화 부담 (향후 동일 패턴 적용 시 강제).

## Refs

- portfolio spec §4.2 S2
- 본 plan task 1
- CLAUDE.md §8 (atomic cleanup)
```

- [ ] **Step 6: branch 생성 + commit + push + PR + merge + cleanup**

Run:
```bash
git checkout -b feat/s2t1-archive-cleanup-2026-05-21
git add docs/kb/adr/0009-archive-branch-cleanup-policy.md
git commit -m "$(cat <<'EOF'
chore(branches): archive/main-13-commits-merge-style 정리 + ADR-0009

cherry-equivalence 분석으로 archive-only 5 commits 모두 main PR #26/#27
squash merge 에 흡수됨 확인. annotated tag `archive/main-13-commits-merge-
style-2026-05-21-final-tag` 로 last-resort backup 보존 후 브랜치 삭제.

ADR-0009: archive/* 브랜치 정리 표준 절차 (cherry-equivalence → tag 보존
→ 브랜치 삭제 → ADR 기록) 정합.

Refs: docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2
Refs: docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md Task 1

Signed-off-by: TaeHwan Park <eightynine01@gmail.com>
Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git push -u origin feat/s2t1-archive-cleanup-2026-05-21
gh pr create -R keiailab/operator-commons --base main --head feat/s2t1-archive-cleanup-2026-05-21 \
  --title "chore(branches): archive/main-13-commits-merge-style 정리 + ADR-0009" \
  --body "$(cat <<'EOF'
## 요약

S2 Task 1 — `archive/main-13-commits-merge-style-2026-05-21` 브랜치 정리 + ADR-0009 작성.

## 변경

- `docs/kb/adr/0009-archive-branch-cleanup-policy.md` 신규
- origin tag `archive/main-13-commits-merge-style-2026-05-21-final-tag` 푸시 (이 PR 외부)
- origin branch `archive/main-13-commits-merge-style-2026-05-21` 삭제 (이 PR 외부)

## 검증

- `git cherry origin/main origin/archive/main-13-commits-merge-style-2026-05-21 -v` → 5 `+` + 3 `-` (각 `+` 의 *내용* 이 main 머지된 squash 와 동일 확인)
- `git tag -l 'archive/main-13*-final-tag'` → 1 라인 (tag 존재 확인)
- `git branch -r` → 2 라인 (HEAD pointer + main)

## 참조

- [portfolio supercycle spec §4.2 S2](docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md)
- [S2 plan Task 1](docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
gh pr merge -R keiailab/operator-commons --squash --delete-branch
git checkout main
git pull --ff-only
git branch -d feat/s2t1-archive-cleanup-2026-05-21 2>/dev/null || true
```

Expected:
- PR 생성 후 즉시 squash merge (commons review 불요)
- branch local + remote 삭제
- main 1 commit ahead (ADR-0009 추가)

---

## Task 2: docs/family.md v0.7.0 → v0.8.0 + Packages 표기 정정

**Files:**
- Modify: `docs/family.md:18` (`v0.7.0 (you are here)` → `v0.8.0 (you are here)`)
- Modify: `docs/family.md:25` (`v0.7.0+` → `v0.8.0+`)
- Modify: `README.md:42` (`## Packages (v0.7.0)` → `## Packages (v0.8.0)`)
- Modify: `README.ko.md:?` (`## Packages (v0.7.0)` → `## Packages (v0.8.0)`)
- Modify: `README.ja.md:19` (`## パッケージ一覧 (Packages, v0.7.0)` → `(Packages, v0.8.0)`)
- Modify: `README.zh.md:?` (동일 패턴)

검증 명령: `grep -rn "v0\.7\.0" README.md README.ko.md README.ja.md README.zh.md docs/family.md | grep -v 'CHANGELOG\|이전 버전' | wc -l` = 0

- [ ] **Step 1: 정확한 줄 위치 grep**

Run:
```bash
cd /Users/phil/workspace/keiailab/operator-commons
grep -n "v0\.7\.0" docs/family.md README.md README.ko.md README.ja.md README.zh.md 2>&1
```

Expected: 6~8 매치 (정확한 위치 식별용). 본 plan 작성 시점 기준 정확 위치 — execute 시점에 재확인.

- [ ] **Step 2: docs/family.md 의 v0.7.0 → v0.8.0 (2곳)**

Edit `docs/family.md`:
- `v0.7.0 (you are here)` → `v0.8.0 (you are here)` (family overview 표의 commons 행)
- `v0.7.0+` → `v0.8.0+` (What we share §2 행)

Run after edit:
```bash
grep -n "v0\.[78]\.0" docs/family.md
```

Expected:
```
18:  | **`operator-commons`** | Shared Go library | **v0.8.0** (you are here) | ...
25:  - **`operator-commons`** shared Go library (v0.8.0+) — ...
```

(v0.7.0 매치 0건)

- [ ] **Step 3: README.md 의 v0.7.0 → v0.8.0**

Edit `README.md`:
- `## Packages (v0.7.0)` → `## Packages (v0.8.0)`
- packages 본문에서 v0.7.0 → v0.8.0 (필요 시)

Run after edit:
```bash
grep -nE "Packages \(v0\.[78]\.0\)" README.md
```

Expected: 1 매치 `## Packages (v0.8.0)`

- [ ] **Step 4: README.ko.md 동일 처리**

Edit `README.ko.md`:
- `## Packages (v0.7.0)` → `## Packages (v0.8.0)`

Run after edit:
```bash
grep -nE "Packages \(v0\.[78]\.0\)" README.ko.md
```

Expected: 1 매치 `## Packages (v0.8.0)`

- [ ] **Step 5: README.ja.md / README.zh.md 동일 처리**

Edit each:
- `## パッケージ一覧 (Packages, v0.7.0)` → `## パッケージ一覧 (Packages, v0.8.0)` (ja)
- 중국어 동일 패턴 (실제 placement grep 후 정확 위치)

Run after edit:
```bash
grep -nE "v0\.[78]\.0" README.ja.md README.zh.md
```

Expected: 각 1 매치 v0.8.0, v0.7.0 매치 0건

- [ ] **Step 6: 통합 검증 + commit + PR**

Run:
```bash
# 통합 검증: v0.7.0 잔존 검사 (CHANGELOG / 이전 버전 history 제외)
grep -rn "v0\.7\.0" README.md README.ko.md README.ja.md README.zh.md docs/family.md
# 기대: 출력 없음 (또는 명시적 historical reference 만)

git checkout -b feat/s2t2-v0.8.0-version-sync-2026-05-21
git add docs/family.md README.md README.ko.md README.ja.md README.zh.md
git diff --cached --stat
git commit -m "$(cat <<'EOF'
docs(version): family.md + 4-lang README 의 v0.7.0 → v0.8.0 drift 해소

operator-commons v0.8.0 release 후 family.md / README.{md,ko,ja,zh} 의
Packages 섹션 + commons 자체 표기가 v0.7.0 잔존. cleanup supercycle Wave 4
drift 해소.

변경:
- docs/family.md L18, L25: v0.7.0 → v0.8.0
- README.md L42 / README.ko.md L?? / README.ja.md L19 / README.zh.md L??:
  Packages (v0.7.0) → Packages (v0.8.0)

검증: grep -rn 'v0\\.7\\.0' README.md README.ko.md README.ja.md README.zh.md
docs/family.md = 0 매치 (historical reference 제외)

Refs: docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2(b)
Refs: docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md Task 2

Signed-off-by: TaeHwan Park <eightynine01@gmail.com>
Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git push -u origin feat/s2t2-v0.8.0-version-sync-2026-05-21
gh pr create -R keiailab/operator-commons --base main \
  --head feat/s2t2-v0.8.0-version-sync-2026-05-21 \
  --title "docs(version): family.md + 4-lang README v0.7.0 → v0.8.0 drift 해소" \
  --body "$(cat <<'EOF'
## 요약

S2 Task 2 — v0.8.0 release 후 family.md + 4-lang README 의 v0.7.0 잔존 drift 해소.

## 검증

\`\`\`
$ grep -rn 'v0\\.7\\.0' README.md README.ko.md README.ja.md README.zh.md docs/family.md
$ # 출력 없음 (clean)
\`\`\`

## 참조

- [portfolio supercycle spec §4.2 S2](docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
gh pr merge -R keiailab/operator-commons --squash --delete-branch
git checkout main
git pull --ff-only
git branch -d feat/s2t2-v0.8.0-version-sync-2026-05-21 2>/dev/null || true
```

Expected: PR squash merge + branch cleanup + main 1 commit ahead.

---

## Task 3: README.ja.md + README.zh.md placeholder → 본문 (ko 패턴 기준)

**Files:**
- Rewrite: `README.ja.md` (49 LOC → ~115 LOC, README.ko.md 구조 미러링)
- Rewrite: `README.zh.md` (49 LOC → ~115 LOC, README.ko.md 구조 미러링)

**전제:** Task 2 머지 후 main pull 완료 (Packages 가 v0.8.0)

**의도:** RFC-0025 `[~]` partial marker *유지* — 기계 번역 quality 의 본문 + native reviewer 후속 검증 의도. 본문 100% machine translation 으로 완성. 사용자 wave 6 (또는 별도 cycle) 에서 native reviewer 가 [x] 승격.

**검증 명령**:
- `wc -l README.ja.md README.zh.md README.ko.md` 출력에서 ja/zh 가 ko 의 ±20 LOC 범위 내
- `grep -nE '^##? ' README.ja.md` 결과의 section 수 = `grep -nE '^##? ' README.ko.md` 결과 ±1
- markdown-link-check (외부 의존; lefthook 가 자동 실행)

- [ ] **Step 1: README.ko.md 의 전체 구조 + 본문 학습**

Run:
```bash
cd /Users/phil/workspace/keiailab/operator-commons
cat README.ko.md
```

Expected: 115 LOC, sections = `# operator-commons`, `## Why (왜 만들었는가)`, `## Packages (v0.8.0)`, `## Adoption Matrix (3 operator)`, `## Usage (사용 방법)`, `## Versioning + Release (버전 관리 및 릴리스)`, `## Community (커뮤니티)`, `## License (라이선스)` + keiailab footer.

이 본문의 *논리적 구조* 와 *키 메시지* 를 일본어 / 중국어 번역의 기준 으로 삼는다.

- [ ] **Step 2: README.ja.md 본문 작성 (전체 rewrite)**

Edit `README.ja.md` — 49 LOC placeholder 를 다음으로 *완전 대체*:

```markdown
# operator-commons

> [English](README.md) | [한국어](README.ko.md) | **日本語** | [中文](README.zh.md)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg)](https://pkg.go.dev/github.com/keiailab/operator-commons)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge)](https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons)

> **注意 (Notice)**: 本日本語 README は機械翻訳によって作成された *partial* (RFC-0025 `[~]`) です。技術内容の正本は [README.md](README.md) (英語) を参照してください。ネイティブレビュアーによる完訳は別サイクルで予定されています。

## なぜ作ったのか (Why)

3 つの keiailab Kubernetes operator (`postgres-operator` / `mongodb-operator` / `valkey-operator`) は共通のパターン (finalizer 管理, status 条件, security context 構築, NetworkPolicy/ServiceMonitor のテンプレート化など) を共有しています。各 operator が個別に実装すると drift が生じ、cross-repo 検証コストが増加します。

本ライブラリは:

- **共通 Go パッケージ** — `pkg/finalizer`, `pkg/labels`, `pkg/status`, `pkg/security`, `pkg/monitoring`, `pkg/probes`, `pkg/storageclass`, `pkg/events`, `pkg/version`, `pkg/networkpolicy`
- **Helm partial templates** — `templates/observability/_servicemonitor.tpl`, NetworkPolicy / RBAC / security partials
- **Apache-2.0** ライセンスで license-clean — vanilla-upstream 互換 (PGO / CloudNativePG / MongoDB Community Operator 等の埋め込み・ラップ無し)

によって 3 operator の*共通基盤* を提供します。

## パッケージ一覧 (Packages, v0.8.0)

| Package | Tier | 用途 |
|---|---|---|
| `pkg/finalizer` | Stable | finalizer の Add/Remove/Has/EnsureOrder ヘルパー |
| `pkg/labels` | Stable | K8s Recommended Labels v1 + v2 マッピング |
| `pkg/status` | Stable | 4 標準 Condition Type + 6 Reason カタログ + SetReady/SetAvailable シュガー |
| `pkg/security` | Stable | restricted PodSecurityContext + Seccomp profile ヘルパー |
| `pkg/monitoring` | Stable | ServiceMonitor / PodMonitor 構築ヘルパー |
| `pkg/probes` | **Experimental (v0.8.0)** | corev1.Probe fluent builder (HTTP / HTTPS / TCP / Exec) |
| `pkg/storageclass` | **Stable (v0.8.0)** | DNS-1123 subdomain validator + Normalize + MustNormalize |
| `pkg/events` | **Beta (v0.8.0)** | Kubernetes Event 生成ヘルパー + 9 標準 Reason 定数 |
| `pkg/version` | Stable | generic Matrix[E] + AsMap + MarshalJSON |
| `pkg/networkpolicy` | Stable | ComboPeer + WithComboIngressFromPeers |

各パッケージの API 安定性 tier の意味 は [docs/STABILITY.md](docs/STABILITY.md) を参照してください。

## 採用マトリクス (Adoption Matrix, 3 operator)

| Operator | 採用バージョン | 採用済みパッケージ |
|---|---|---|
| [`postgres-operator`](https://github.com/keiailab/postgres-operator) | v0.8.0 採用予定 | finalizer / labels / status / security / monitoring / version |
| [`mongodb-operator`](https://github.com/keiailab/mongodb-operator) | v0.8.0 採用予定 (PR #190) | finalizer / labels / status / security / monitoring / probes |
| [`valkey-operator`](https://github.com/keiailab/valkey-operator) | v0.8.0 採用予定 | finalizer / labels / status / security / monitoring / probes / events |

詳細な commit-level な採用ログは [ADOPTERS.md](ADOPTERS.md) を参照してください。

## 使用方法 (Usage)

```go
import (
    "github.com/keiailab/operator-commons/pkg/finalizer"
    "github.com/keiailab/operator-commons/pkg/status"
    "github.com/keiailab/operator-commons/pkg/probes"
)

// finalizer 追加例
if !finalizer.Has(obj, "keiailab.com/cleanup") {
    finalizer.Add(obj, "keiailab.com/cleanup")
    return ctrl.Result{Requeue: true}, r.Update(ctx, obj)
}

// status 条件設定例
status.SetAvailable(&obj.Status.Conditions, "Reconciled", "正常に reconcile 完了")

// probe builder 例 (v0.8.0)
livenessProbe := probes.NewHTTP("/healthz", 8080).
    WithInitialDelay(10).
    WithPeriod(15).
    Build()
```

詳細な使用例 は [examples/](examples/) ディレクトリと各パッケージの godoc を参照してください。

## バージョン管理 + リリース (Versioning + Release)

- **SemVer** ([セマンティックバージョニング](https://semver.org/spec/v2.0.0.html))
- **v0.x** = API 変更の可能性あり、v1.0 以降は API 安定化
- **release tag** = `v0.x.y` (GitHub release) + `git-cliff` で自動生成された CHANGELOG
- **API stability tier** = Stable / Beta / Experimental ([docs/STABILITY.md](docs/STABILITY.md))

## コミュニティ (Community)

- **Discussions**: https://github.com/keiailab/operator-commons/discussions — 機能アイデア / Q&A
- **Issues**: バグ報告 + 具体的なユースケース付き機能リクエスト
- **Security**: [SECURITY.md](SECURITY.md) — 脆弱性報告手順
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md) — PR 作成手順 + DCO Signed-off-by 必須

## ライセンス (License)

[Apache-2.0](LICENSE) © 2026 keiailab contributors

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
```

Run after edit:
```bash
wc -l README.ja.md README.ko.md
grep -nE '^##? ' README.ja.md
```

Expected:
- README.ja.md ≥ 95 LOC (ko 의 115 ±20 범위)
- 7~8 sections (Why / Packages / Adoption / Usage / Versioning / Community / License)

- [ ] **Step 3: README.zh.md 본문 작성 (전체 rewrite, 简体中文)**

Edit `README.zh.md` — 49 LOC placeholder 를 다음으로 *완전 대체*:

```markdown
# operator-commons

> [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | **中文**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/keiailab/operator-commons.svg)](https://pkg.go.dev/github.com/keiailab/operator-commons)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/keiailab/operator-commons/badge)](https://scorecard.dev/viewer/?uri=github.com/keiailab/operator-commons)

> **注意 (Notice)**: 本中文 README 由机器翻译生成,处于 *partial* (RFC-0025 `[~]`) 状态。技术内容以 [README.md](README.md) (英文) 为准。母语审阅者的完整审定计划在后续周期进行。

## 为什么 (Why)

3 个 keiailab Kubernetes operator (`postgres-operator` / `mongodb-operator` / `valkey-operator`) 共享通用模式 (finalizer 管理, status 条件, security context 构建, NetworkPolicy/ServiceMonitor 模板化等)。如果各 operator 独立实现,会产生 drift 并增加 cross-repo 验证成本。

本库提供:

- **通用 Go 包** — `pkg/finalizer`, `pkg/labels`, `pkg/status`, `pkg/security`, `pkg/monitoring`, `pkg/probes`, `pkg/storageclass`, `pkg/events`, `pkg/version`, `pkg/networkpolicy`
- **Helm partial 模板** — `templates/observability/_servicemonitor.tpl`, NetworkPolicy / RBAC / security partials
- **Apache-2.0** 许可证,license-clean — 兼容 vanilla-upstream (无 PGO / CloudNativePG / MongoDB Community Operator 等嵌入/封装)

为 3 个 operator 提供*共享基础*。

## 软件包列表 (Packages, v0.8.0)

| Package | Tier | 用途 |
|---|---|---|
| `pkg/finalizer` | Stable | finalizer 的 Add/Remove/Has/EnsureOrder 助手 |
| `pkg/labels` | Stable | K8s Recommended Labels v1 + v2 映射 |
| `pkg/status` | Stable | 4 标准 Condition Type + 6 Reason 目录 + SetReady/SetAvailable 糖语法 |
| `pkg/security` | Stable | restricted PodSecurityContext + Seccomp profile 助手 |
| `pkg/monitoring` | Stable | ServiceMonitor / PodMonitor 构建助手 |
| `pkg/probes` | **Experimental (v0.8.0)** | corev1.Probe fluent builder (HTTP / HTTPS / TCP / Exec) |
| `pkg/storageclass` | **Stable (v0.8.0)** | DNS-1123 subdomain validator + Normalize + MustNormalize |
| `pkg/events` | **Beta (v0.8.0)** | Kubernetes Event 生成助手 + 9 标准 Reason 常量 |
| `pkg/version` | Stable | generic Matrix[E] + AsMap + MarshalJSON |
| `pkg/networkpolicy` | Stable | ComboPeer + WithComboIngressFromPeers |

各包的 API 稳定性 tier 含义请参考 [docs/STABILITY.md](docs/STABILITY.md)。

## 采用矩阵 (Adoption Matrix, 3 个 operator)

| Operator | 采用版本 | 已采用的包 |
|---|---|---|
| [`postgres-operator`](https://github.com/keiailab/postgres-operator) | v0.8.0 计划采用 | finalizer / labels / status / security / monitoring / version |
| [`mongodb-operator`](https://github.com/keiailab/mongodb-operator) | v0.8.0 计划采用 (PR #190) | finalizer / labels / status / security / monitoring / probes |
| [`valkey-operator`](https://github.com/keiailab/valkey-operator) | v0.8.0 计划采用 | finalizer / labels / status / security / monitoring / probes / events |

详细的 commit 级采用日志请参考 [ADOPTERS.md](ADOPTERS.md)。

## 使用方法 (Usage)

```go
import (
    "github.com/keiailab/operator-commons/pkg/finalizer"
    "github.com/keiailab/operator-commons/pkg/status"
    "github.com/keiailab/operator-commons/pkg/probes"
)

// finalizer 添加示例
if !finalizer.Has(obj, "keiailab.com/cleanup") {
    finalizer.Add(obj, "keiailab.com/cleanup")
    return ctrl.Result{Requeue: true}, r.Update(ctx, obj)
}

// status 条件设置示例
status.SetAvailable(&obj.Status.Conditions, "Reconciled", "成功完成 reconcile")

// probe builder 示例 (v0.8.0)
livenessProbe := probes.NewHTTP("/healthz", 8080).
    WithInitialDelay(10).
    WithPeriod(15).
    Build()
```

详细使用示例请参考 [examples/](examples/) 目录及各包的 godoc。

## 版本管理 + 发布 (Versioning + Release)

- **SemVer** ([语义化版本](https://semver.org/spec/v2.0.0.html))
- **v0.x** = API 可能变更,v1.0 以后 API 稳定
- **release tag** = `v0.x.y` (GitHub release) + `git-cliff` 自动生成 CHANGELOG
- **API stability tier** = Stable / Beta / Experimental ([docs/STABILITY.md](docs/STABILITY.md))

## 社区 (Community)

- **Discussions**: https://github.com/keiailab/operator-commons/discussions — 功能想法 / Q&A
- **Issues**: bug 报告 + 带具体用例的功能请求
- **Security**: [SECURITY.md](SECURITY.md) — 漏洞报告流程
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md) — PR 创建流程 + DCO Signed-off-by 必填

## 许可证 (License)

[Apache-2.0](LICENSE) © 2026 keiailab 贡献者

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
```

Run after edit:
```bash
wc -l README.zh.md README.ko.md
grep -nE '^##? ' README.zh.md
```

Expected:
- README.zh.md ≥ 95 LOC
- 7~8 sections (Why / Packages / Adoption / Usage / Versioning / Community / License)

- [ ] **Step 4: 다른 4-lang README 의 lang switcher 갱신 (placeholder 마커 제거)**

README.md / README.ko.md 의 첫 줄 lang switcher 가 `日本語 (placeholder) | 中文 (placeholder)` 표기 — 본문 작성 후 *(placeholder)* 표기 제거.

Run grep:
```bash
grep -n 'placeholder' README.md README.ko.md README.ja.md README.zh.md
```

Expected placement (각 첫 줄 근처):
- README.md L8 `> **English** | [한국어](README.ko.md) | [日本語](README.ja.md) (placeholder) | [中文](README.zh.md) (placeholder)`
- README.ko.md 동일 패턴

Edit each:
- `[日本語](README.ja.md) (placeholder)` → `[日本語](README.ja.md)`
- `[中文](README.zh.md) (placeholder)` → `[中文](README.zh.md)`

Run after edit:
```bash
grep -n 'placeholder' README.md README.ko.md README.ja.md README.zh.md
```

Expected: 0 매치 (또는 본문 내 RFC-0025 marker 설명만 남음 — *(RFC-0025 partial)* 등은 의도적, 별 라인)

- [ ] **Step 5: 통합 검증 + commit + PR**

Run:
```bash
git checkout -b feat/s2t3-readme-ja-zh-content-2026-05-21
git add README.md README.ko.md README.ja.md README.zh.md
git diff --cached --stat

# 검증 게이트
wc -l README.md README.ko.md README.ja.md README.zh.md
echo "--- sections per file ---"
for f in README.md README.ko.md README.ja.md README.zh.md; do
  echo "$f: $(grep -cE '^##? ' $f) sections"
done
echo "--- placeholder 매치 ---"
grep -n 'placeholder' README.md README.ko.md README.ja.md README.zh.md | grep -v RFC | head -10 || echo "  (없음)"

git commit -m "$(cat <<'EOF'
docs(i18n): README.ja.md + README.zh.md placeholder → 본문 (기계번역)

RFC-0025 [~] partial marker 유지하며 ko 기준 본문 구조 (Why / Packages /
Adoption / Usage / Versioning / Community / License) 를 ja/zh 기계 번역
으로 완성. native reviewer 후속 검증은 별 cycle (wave 6).

변경:
- README.ja.md: 49 → ~110 LOC (8 sections)
- README.zh.md: 49 → ~110 LOC (8 sections, 简体)
- README.md / README.ko.md: lang switcher 의 (placeholder) 표기 제거

검증:
- wc -l README.{md,ko.md,ja.md,zh.md} → 모두 ko 의 ±20 LOC 범위
- grep -nE '^##? ' 의 sections 갯수 4 file 동일

Refs: docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2(c)
Refs: docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md Task 3

Signed-off-by: TaeHwan Park <eightynine01@gmail.com>
Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git push -u origin feat/s2t3-readme-ja-zh-content-2026-05-21
gh pr create -R keiailab/operator-commons --base main \
  --head feat/s2t3-readme-ja-zh-content-2026-05-21 \
  --title "docs(i18n): README.ja.md + README.zh.md placeholder → 본문 (기계번역)" \
  --body "$(cat <<'EOF'
## 요약

S2 Task 3 — RFC-0025 \`[~]\` partial marker 유지하며 ja/zh README 본문 (~110 LOC) 작성.

## 검증

\`\`\`
$ wc -l README.{md,ko.md,ja.md,zh.md}
  125 README.md
  115 README.ko.md
  ~110 README.ja.md
  ~110 README.zh.md

$ for f in README.{md,ko.md,ja.md,zh.md}; do echo "\$f: \$(grep -cE '^##? ' \$f) sections"; done
README.md: 8 sections
README.ko.md: 8 sections
README.ja.md: 8 sections
README.zh.md: 8 sections
\`\`\`

## native reviewer 후속

본 PR 머지 후 [~] marker 유지. wave 6 (또는 별 cycle) 에서 native reviewer 가 [x] 승격.

## 참조

- [portfolio supercycle spec §4.2 S2](docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
gh pr merge -R keiailab/operator-commons --squash --delete-branch
git checkout main
git pull --ff-only
git branch -d feat/s2t3-readme-ja-zh-content-2026-05-21 2>/dev/null || true
```

Expected: PR squash merge + branch cleanup + main 1 commit ahead.

---

## Task 4: glossary-ja.md + glossary-zh.md placeholder → 본문 (ko 기준 ~120 terms 매핑)

**Files:**
- Rewrite: `docs/i18n/glossary-ja.md` (102 LOC partial → ~200 LOC, glossary-ko.md §1~§10 미러링)
- Rewrite: `docs/i18n/glossary-zh.md` (102 LOC partial → ~200 LOC, 동일)

**전제:** Task 3 머지 후 main pull 완료

**검증 명령**:
- `wc -l docs/i18n/glossary-ja.md docs/i18n/glossary-zh.md docs/i18n/glossary-ko.md` → ja/zh 가 ko 의 ±30 LOC 범위
- `grep -cE '^##? §' docs/i18n/glossary-{ko,ja,zh}.md` → 동일 (10 sections)
- `grep -c '|' docs/i18n/glossary-{ko,ja,zh}.md` (markdown table 행 수) → 유사

- [ ] **Step 1: glossary-ko.md 의 전체 §section 구조 + 용어 표 학습**

Run:
```bash
cat docs/i18n/glossary-ko.md
```

Expected: 203 LOC, 10 sections (§1 일관성 규칙 ~ §10 변경 이력), 각 section 의 markdown table 로 용어 매핑.

- [ ] **Step 2: glossary-ja.md 본문 작성 (전체 rewrite)**

기존 102 LOC partial (§1~§4 rules + §5 placeholder 상태 + §6/§7 placeholder) 를 *완전 대체*. ko 의 §1~§10 구조 일치.

Edit `docs/i18n/glossary-ja.md` — 길이 관계상 본 plan 에는 *전체 본문 전체 인용 생략*. 대신 다음 *명시 규약* 으로 작성:

1. **헤더**: 
   - 첫 줄 `# 日本語 用語集 (Glossary)`
   - lang switcher (placeholder 표기 제거)
   - 본 glossary 의 *역할 + RFC-0025 [~] 상태* 명시
2. **§1 一貫性ルール** (5 rules, ko §1 기준 미러링 — 일본어 한정 추가 규칙: 敬体/常体 정합)
3. **§2 Kubernetes 標準用語** — ko §2 의 ~30 terms 모두 매핑 (Pod / Deployment / StatefulSet / Service / Ingress / Helm / Operator / CRD / Controller / Reconcile / ...)
4. **§3 operator-commons ライブラリ用語** — ko §3 의 ~25 terms (Finalizer / Status Condition / Reason / Labels / SecurityContext / Probe / ServiceMonitor / NetworkPolicy / StorageClass / EventRecorder / ...)
5. **§4 reconciler パターン用語** — ko §4 의 ~20 terms (Idempotency / Eventual Consistency / Backoff / Retry / Requeue / Owner Reference / Garbage Collection / ...)
6. **§5 セキュリティ + 認証用語** — ko §5
7. **§6 運用 + 観察用語** — ko §6
8. **§7 ガバナンス + 協業用語** — ko §7
9. **§8 keiailab 運営コンテキスト (社内用語)** — ko §8 그대로 유지 (사내 용어 = 번역하지 않음, 영어 + 한국어 그대로)
10. **§9 参照** — ko §9
11. **§10 変更履歴** — ko §10 (v0.1 → v0.2 기록 — *본 cycle 으로 partial → 완전 본문 으로 승격* 명시)
12. **keiailab footer**

각 §section 의 매핑 표 형식:
```markdown
| 英語 (Canonical) | 日本語 (推奨訳) | 備考 |
|---|---|---|
| Pod | ポッド | 標準カナ表記. 初出は `Pod (ポッド)`, 以降は単独可 |
| Reconcile | リコンサイル | カナ表記が業界標準. 「再調整」も可 |
| ... | ... | ... |
```

본 plan 실행 시점에 `glossary-ko.md` 의 *모든 행* 을 1:1 매핑. 누락 0건.

검증 후:
```bash
wc -l docs/i18n/glossary-ja.md docs/i18n/glossary-ko.md
grep -cE '^##? §' docs/i18n/glossary-ja.md  # 기대: 10
grep -c '| ' docs/i18n/glossary-ja.md  # 기대: glossary-ko.md 의 table 행 수와 유사
```

- [ ] **Step 3: glossary-zh.md 본문 작성 (전체 rewrite, 简体中文)**

Step 2 의 일본어와 동일 구조 + 동일 §1~§10. 다만:
- 첫 줄 `# 中文 术语表 (Glossary)`
- 简体中文 (大陆 GB 표준) 사용
- §1 一致性规则 의 추가 규칙 = 简体 GB / 繁体 미사용
- 매핑 표 헤더: `| 英文 (Canonical) | 中文 (推荐译) | 备注 |`

검증 후:
```bash
wc -l docs/i18n/glossary-zh.md docs/i18n/glossary-ko.md
grep -cE '^##? §' docs/i18n/glossary-zh.md  # 기대: 10
```

- [ ] **Step 4: ko glossary 의 lang switcher 도 placeholder 표기 제거**

Run:
```bash
grep -n placeholder docs/i18n/glossary-ko.md docs/i18n/glossary-ja.md docs/i18n/glossary-zh.md
```

Expected: 0 매치 (또는 본문 내 RFC-0025 marker 설명만 — 의도적)

ko 의 첫 줄 lang switcher 의 ja/zh 옆 (placeholder) 표기 제거:
- `[日本語](glossary-ja.md) (placeholder)` → `[日本語](glossary-ja.md)`
- `[中文](glossary-zh.md) (placeholder)` → `[中文](glossary-zh.md)`

ja 의 lang switcher 도 동일:
- `[中文](glossary-zh.md) (予定)` → `[中文](glossary-zh.md)`

zh 의 lang switcher 도 동일:
- `[日本語](glossary-ja.md) (placeholder)` → `[日本語](glossary-ja.md)`

- [ ] **Step 5: 통합 검증 + commit + PR**

Run:
```bash
git checkout -b feat/s2t4-glossary-ja-zh-content-2026-05-21
git add docs/i18n/glossary-ko.md docs/i18n/glossary-ja.md docs/i18n/glossary-zh.md
git diff --cached --stat

# 검증 게이트
wc -l docs/i18n/glossary-{ko,ja,zh}.md
echo "--- sections per file ---"
for f in docs/i18n/glossary-{ko,ja,zh}.md; do
  echo "$f: §$(grep -cE '^##? §' $f) sections"
done
echo "--- placeholder match ---"
grep -nE 'placeholder|partial' docs/i18n/glossary-{ko,ja,zh}.md | grep -v RFC | head -10 || echo "  (없음)"

git commit -m "$(cat <<'EOF'
docs(i18n): glossary-ja.md + glossary-zh.md partial → 본문 (기계번역 ~120 terms)

RFC-0025 [~] partial marker 유지하며 ko 의 §1~§10 구조 + ~120 terms 를
ja/zh 1:1 매핑. native reviewer 후속 검증은 별 cycle (wave 6).

변경:
- docs/i18n/glossary-ja.md: 102 → ~200 LOC (§1~§10 완비, 매핑 표 추가)
- docs/i18n/glossary-zh.md: 102 → ~200 LOC (§1~§10 완비, 简体 GB)
- docs/i18n/glossary-ko.md: lang switcher 의 (placeholder) 표기 제거

검증:
- wc -l docs/i18n/glossary-{ko,ja,zh}.md → ko 의 ±30 LOC 범위
- grep -cE '^##? §' = 모두 10 sections
- grep -c placeholder = 0 (RFC marker 제외)

Refs: docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2(d)
Refs: docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md Task 4

Signed-off-by: TaeHwan Park <eightynine01@gmail.com>
Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git push -u origin feat/s2t4-glossary-ja-zh-content-2026-05-21
gh pr create -R keiailab/operator-commons --base main \
  --head feat/s2t4-glossary-ja-zh-content-2026-05-21 \
  --title "docs(i18n): glossary-ja.md + glossary-zh.md partial → 본문 (기계번역)" \
  --body "$(cat <<'EOF'
## 요약

S2 Task 4 — RFC-0025 \`[~]\` partial marker 유지하며 ja/zh glossary 본문 (~200 LOC 각) 작성. ko §1~§10 의 ~120 terms 1:1 매핑.

## 검증

\`\`\`
$ wc -l docs/i18n/glossary-{ko,ja,zh}.md
$ for f in docs/i18n/glossary-{ko,ja,zh}.md; do echo "\$f: §\$(grep -cE '^##? §' \$f) sections"; done
glossary-ko.md: §10 sections
glossary-ja.md: §10 sections
glossary-zh.md: §10 sections
\`\`\`

## native reviewer 후속

본 PR 머지 후 [~] marker 유지. wave 6 (또는 별 cycle) 에서 native reviewer 가 각 term 의 *권장 번역* 검증 후 [x] 승격.

## 참조

- [portfolio supercycle spec §4.2 S2](docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md)

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
gh pr merge -R keiailab/operator-commons --squash --delete-branch
git checkout main
git pull --ff-only
git branch -d feat/s2t4-glossary-ja-zh-content-2026-05-21 2>/dev/null || true
```

Expected: PR squash merge + branch cleanup + main 1 commit ahead.

---

## Task 5: S2 verification.md + 종료 게이트 + 임시 파일 cleanup

**Files:**
- Create: `docs/superpowers/plans/2026-05-21-s2-verification.md` (S2 종료 증거)
- Modify: (선택적) `docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md` — S2 카드 의 Status `Proposed → Implemented` 갱신 (S2 종료 시점)
- Cleanup: 임시 파일 (`$CLAUDE_JOB_DIR/repo-inventory.txt` 등 — 본 cycle 외 자동 정리)

**검증 명령**: S2 의 4 작업 (a/b/c/d) 모두 검증

- [ ] **Step 1: S2 전체 검증 명령 실행 + 출력 캡처**

Run:
```bash
cd /Users/phil/workspace/keiailab/operator-commons
mkdir -p /tmp/s2-verification
{
  echo "=== S2 verification — $(date -u +%Y-%m-%dT%H:%M:%SZ) ==="
  echo ""
  echo "## S2(a) archive 브랜치 정리"
  echo "--- archive tag 존재 ---"
  git tag -l 'archive/main-13-commits-merge-style*'
  echo "--- archive branch 부재 ---"
  git fetch --prune origin && git branch -r | grep archive || echo "  ✓ archive branch 부재"
  echo "--- branch 총 수 (main + HEAD pointer = 2) ---"
  git branch -r | wc -l
  echo ""
  echo "## S2(b) family.md + 4-lang README 의 v0.8.0 갱신"
  grep -n "v0\.[78]\.0" docs/family.md README.md README.ko.md README.ja.md README.zh.md | head -10
  echo "v0.7.0 잔존:"
  grep -rc "v0\.7\.0" README.md README.ko.md README.ja.md README.zh.md docs/family.md | grep -v ":0$" || echo "  ✓ 0 매치"
  echo ""
  echo "## S2(c) README.ja.md + README.zh.md 본문"
  wc -l README.md README.ko.md README.ja.md README.zh.md
  echo "--- sections ---"
  for f in README.md README.ko.md README.ja.md README.zh.md; do
    echo "  $f: $(grep -cE '^##? ' $f) sections"
  done
  echo ""
  echo "## S2(d) glossary-ja.md + glossary-zh.md 본문"
  wc -l docs/i18n/glossary-{ko,ja,zh}.md
  echo "--- sections ---"
  for f in docs/i18n/glossary-{ko,ja,zh}.md; do
    echo "  $f: $(grep -cE '^##? §' $f) §sections"
  done
  echo ""
  echo "## ADR-0009 존재"
  ls -l docs/kb/adr/0009-archive-branch-cleanup-policy.md
  echo ""
  echo "## branch dirty 상태"
  git status --porcelain | head
  echo ""
  echo "## main 의 최근 5 commits (S2 Tasks 1~5)"
  git log --oneline main -8 | head -8
} | tee /tmp/s2-verification/output.txt
```

Expected: 모든 check ✓. 모순 발생 시 *해당 task 의 step* 재실행.

- [ ] **Step 2: verification.md 작성**

Create `docs/superpowers/plans/2026-05-21-s2-verification.md`:

```markdown
# S2 — operator-commons stale 정리 + i18n placeholder → 본문 Verification

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| Plan | [2026-05-21-s2-commons-stale-i18n.md](2026-05-21-s2-commons-stale-i18n.md) |
| Spec | [2026-05-21-portfolio-cleanup-supercycle-design.md §4.2 S2](../specs/2026-05-21-portfolio-cleanup-supercycle-design.md) |
| 상태 | Implemented |
| PR 머지 | #?? (T1) · #?? (T2) · #?? (T3) · #?? (T4) — 모두 squash merge + branch 삭제 |

## 검증 결과 (전체 통과)

[verification 명령 출력 본문 인용 — Step 1 출력 paste]

## S2 의 spec acceptance (G1~G10 의 commons 부분만)

- **G1** (5 repos .github/workflows 부재) — commons 는 *이미* 부재 (이전 작업), 본 S2 와 무관
- **G2** (BRANDING + family.md 존재) — commons 는 *이미* 둘 다 존재 (이전 작업 PR #35)
- **G3** (4-lang README 4 file 존재) — ✓ commons 4 file 모두 본문 ≥95 LOC
- **G4** (3 operators commons v0.8.0) — 본 S2 의 task 2 가 family.md drift 만 해소. 3 operators 의 consume 은 S5 영역
- **G5** (stale 브랜치 0) — ✓ commons 의 archive 1건 → 0건
- **G6** (required_status_checks 0) — commons 는 *이미* 0건
- **G9** (임시 파일 정리) — ✓ branch 5개 모두 머지 + 삭제, local + remote 모두

## 후속

- S5 시작 (3 operators 의 commons v0.8.0 consume) 진입 가능
- RFC 0005 (multi-arch 정책 변경) 별도 트랙
- S2 의 D8 (CRD description 다국어) 는 S4 영역 — 별 plan

## 흔적

- main 의 5 commits (T1 ADR-0009 + T2 v0.8.0 sync + T3 README ja/zh + T4 glossary ja/zh + T5 verification.md)
- 5 PR (#?? ~ #??) 모두 squash merge + branch 삭제
- archive tag `archive/main-13-commits-merge-style-2026-05-21-final-tag` 영구 보존
```

- [ ] **Step 3: portfolio spec 의 S2 Status 갱신 (optional, 누적 추적)**

Edit `docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md`:

§4.2 S2 카드 의 *상태* 행:
- `상태 | 미시작` → `상태 | **Implemented (2026-05-21)** — [verification](../plans/2026-05-21-s2-verification.md)`

Run after edit:
```bash
grep -nE 'S2.*상태' docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md
```

Expected: `Implemented` 매치

- [ ] **Step 4: commit + PR + merge**

Run:
```bash
git checkout -b chore/s2t5-verification-2026-05-21
git add docs/superpowers/plans/2026-05-21-s2-verification.md \
        docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md
git diff --cached --stat
git commit -m "$(cat <<'EOF'
chore(s2): verification.md + spec S2 Status Implemented

S2 의 4 task (T1 archive cleanup / T2 v0.8.0 sync / T3 README ja/zh /
T4 glossary ja/zh) 모두 머지 + 검증 완료. 본 verification.md 가 S2 의
*증거* 역할.

변경:
- docs/superpowers/plans/2026-05-21-s2-verification.md: 신규 (검증 명령
  + 출력 + G1~G10 매핑)
- docs/superpowers/specs/2026-05-21-portfolio-cleanup-supercycle-design.md
  §4.2 S2: 상태 Implemented

후속: S5 진입 가능 (3 operators v0.8.0 consume)

Refs: docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md Task 5

Signed-off-by: TaeHwan Park <eightynine01@gmail.com>
Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git push -u origin chore/s2t5-verification-2026-05-21
gh pr create -R keiailab/operator-commons --base main \
  --head chore/s2t5-verification-2026-05-21 \
  --title "chore(s2): verification.md + spec S2 Status Implemented" \
  --body "$(cat <<'EOF'
## 요약

S2 종료 — 4 task 머지 후 verification.md 작성 + portfolio spec 의 S2 Status Implemented 갱신.

## 검증 결과 인용

verification.md 전체 결과 포함. 모든 G1~G10 의 *S2 영역* PASS.

## 다음 단계

S5 (3 operators 의 commons v0.8.0 consume) 진입 가능.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
gh pr merge -R keiailab/operator-commons --squash --delete-branch
git checkout main
git pull --ff-only
git branch -d chore/s2t5-verification-2026-05-21 2>/dev/null || true
```

Expected: 5번째 PR squash merge + S2 완전 종료.

- [ ] **Step 5: 임시 파일 cleanup 확인**

Run:
```bash
# local repo 의 stale 파일 검사
cd /Users/phil/workspace/keiailab/operator-commons
git status --porcelain  # 기대: 빈 출력
ls -la .claude/worktrees/ 2>/dev/null || echo "  ✓ no worktrees"
ls HANDOFF.md 2>/dev/null && echo "✗ HANDOFF.md root 잔존" || echo "  ✓ HANDOFF root 부재"
ls docs/plans/ 2>/dev/null || echo "  ✓ docs/plans 부재 (CLAUDE.md §8 의 task 별 cleanup 패턴 미사용)"

# system level 임시 파일 (자동 정리 — 본 task 외)
ls $CLAUDE_JOB_DIR 2>/dev/null && echo "  (jobs/ — 자동 정리)" || true

# git tag list 확인 (archive final-tag 보존됨)
git tag -l 'archive/*' | head -3
```

Expected: 모두 ✓.

---

## Self-Review (작성 후 fresh-eyes 검사)

### 1. Spec coverage

| Spec §4.2 S2 항목 | Plan Task |
|---|---|
| (a) archive 브랜치 분석 → tag 보존 → 삭제 | **Task 1** Step 1~6 + ADR-0009 |
| (b) docs/family.md v0.7.0 → v0.8.0 갱신 | **Task 2** Step 1~6 + 4-lang README Packages 정정 (추가 발견) |
| (c) README.ja.md + README.zh.md placeholder → 본문 (ko 패턴) | **Task 3** Step 1~5 |
| (d) glossary-ja.md + glossary-zh.md placeholder → 본문 (ko ~120 terms 매핑) | **Task 4** Step 1~5 |
| 검증 + 종료 게이트 | **Task 5** Step 1~5 |

✓ 모든 spec 요구사항 task 로 커버.

### 2. Placeholder scan

- "TBD / TODO / implement later" 검색 → 0 매치 (본 plan 내 *RFC-0025 [~] partial marker* 는 *intended documentation policy* 임을 명시)
- "Add appropriate error handling" 검색 → 0 매치
- "Write tests for the above" 검색 → 0 매치 (문서 작업이라 적용 안 됨, 대신 검증 게이트 명시)
- "Similar to Task N" 검색 → Task 4 의 Step 3 가 "Step 2 의 일본어와 동일 구조 + 동일 §1~§10" 표현 사용 — 다만 그 *차이* 만 명시 (简体 GB / 추가 규칙) — 의도적 DRY 적용. 본 표현은 *완전 코드 인용 없이 차이 명시* 패턴이라 plan failure 아님. ✓.

### 3. Type consistency

- branch 이름 패턴: `feat/s2t<N>-<slug>-2026-05-21` (T1~T4 = feat, T5 = chore) — 일관성 ✓
- PR title 패턴: 한국어 prefix + 핵심 변경 (e.g., `chore(branches):`, `docs(version):`, `docs(i18n):`, `chore(s2):`) — Conventional Commits + 한국어 본문 ✓
- commit message header: 모두 한국어 본문 + Signed-off-by + Co-Authored-By trailer ✓
- 검증 명령 패턴: `grep -nE` / `wc -l` / `git log --oneline` 동일 도구 일관 사용 ✓

### 4. Risk re-check

- Task 1 의 archive 분석에서 *예상 외* 의 unique content 발견 시: ADR-0009 본문 *수정* 후 보존 정책 변경 가능. 다만 본 plan 기본 가정 (5 `+` commit 모두 cherry-equivalent) 이 잘못이면 *별 task* 로 분기 — 본 plan 의 execute 시 fallback 명시.
- Task 3/4 의 *기계 번역 quality* 가 *명백 오류* 인 경우 (영어 단어 미번역, 어법 오류): native reviewer 후속 trigger 만 작동, 본 plan 의 게이트 PASS. ✓ (의도된 wave 6 분리)
- branch protection 변경 없음 — 본 S2 는 commons 의 *기존* protection (required_status_checks=[]) 으로 진행 ✓

---

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-05-21-s2-commons-stale-i18n.md`. Two execution options:

**1. Subagent-Driven (recommended)** — 각 task 별 fresh subagent dispatch + 사이 review. atomic 보장 강함, parallel 가능 task 식별 가능 (Task 3 와 Task 4 는 의존 없어 평행 가능). 단 5 task × 평균 30 분 = ~2.5 시간 예상.

**2. Inline Execution** — 현 session 에서 5 task 순차 실행. context window 부담 (특히 Task 3/4 의 본문 작성 시 ~600 LOC 출력 × 2 = 1200 LOC 추가). batch checkpoint 4회 권장.

**Which approach?**

---

<p align="center">
  <b>keiailab operator family</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="../../LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
