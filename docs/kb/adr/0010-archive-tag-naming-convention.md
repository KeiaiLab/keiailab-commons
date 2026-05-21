# ADR-0010 — archive tag 명명 규약 표준화 (S2 implementation)

| 메타 | 값 |
|---|---|
| Status | Accepted |
| Date | 2026-05-21 |
| Supersedes | — |
| Extends | ADR-0009 (archive/* 브랜치 정리 정책) |
| Refs | docs/specs/2026-05-21-stale-branch-cleanup-design.md §4.2.1 |

## Context

S2 spec (PR #38) 의 §4.2.1 에서 archive 브랜치 보존본 tag 의 신규 명명 규약을 결정했다:
- 기존 tag (2026-05-21 오전 cycle): `archive/main-13-commits-merge-style-2026-05-21-final-tag` (slash prefix + `-final-tag` suffix)
- 신규 tag (본 ADR): `archive-merge-style-v0.7.0` (slash 제거 + 버전 anchor)

기존 tag (`archive/.../-final-tag`) 는 ADR-0009 시점 의 "최후 백업" 의도로 작성되었으나, 명명 규약이 *날짜 기반* 으로 시점-종속적이다. 향후 동일 패턴 archive 처리 시 *어떤 시점의 버전* 이었는지 파악하기 어렵다.

## Decision

1. archive 브랜치 보존본 tag 의 표준 명명: `archive-<context>-v<semver>` (slash 없음).
2. `archive-merge-style-v0.7.0` = main 의 v0.7.0 시점 머지-스타일 commit 백업 (commit `910042a`).
3. 기존 tag (`archive/main-13-commits-merge-style-2026-05-21-final-tag` 외 1건) 은 *immutable archive 의 archive* 로 보존 — 삭제 안 함.
4. 향후 신규 archive tag 는 본 규약 (slash 없음 + semver anchor) 강제.

## Consequences

- (+) tag dropdown UI 에서 의미론적 grouping (`archive-*` prefix) 가능.
- (+) 버전 anchor 로 *어느 release 시점인지* 즉각 파악.
- (+) ADR-0009 의 정책 본문은 유지 (cherry-equivalence 검증 + annotated tag + branch 삭제).
- (-) 동일 archive commit 을 가리키는 tag 2건 공존 (`archive/.../-final-tag` + `archive-merge-style-v0.7.0`) → 기존 tag 는 *historical-exempt* 로 보존.
- (-) 향후 동일 패턴 적용 시 명명 일관성 강제 (lefthook hook 추가 미정).

## Alternatives

- (A) 기존 tag 삭제 + 신규 tag 만 유지 — 기존 tag 작성 commit 의 history 손실 가능 (다른 ref 가 참조 시).
- (B) 신규 tag 생략 + 기존 tag 만 보존 — S2 spec §4.2.1 의 사용자 결정 사항 위배.
- (C) 두 명명 규약 모두 허용 — drift 부채.

## Cross-link

- ADR-0009 (archive/* 브랜치 cleanup 정책 본문)
- spec: docs/specs/2026-05-21-stale-branch-cleanup-design.md §4.2
- 본 S2 implementation cycle (chore/archive-tag-and-cleanup-2026-05-21)
