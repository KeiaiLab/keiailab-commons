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
