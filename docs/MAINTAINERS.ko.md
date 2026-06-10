# Maintainers

> [English](MAINTAINERS.md) | **한국어** | [日本語](MAINTAINERS.ja.md) | [中文](MAINTAINERS.zh.md)

본 문서는 `keiailab/keiailab-commons` 의 의사결정 권한을 가진 메인테이너
명단을 관리합니다.

## 현재 메인테이너

| 이름 / 팀 | GitHub | 역할 | 담당 영역 |
|---|---|---|---|
| keiailab maintainers | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | Lead | 전체 |

GitHub team `@keiailab/maintainers` 가 본 라이브러리의 모든 영역에 대한 머지 ·
릴리스 태그 권한을 보유합니다.

## 메인테이너 자격

downstream consumer operator 의 메인테이너이거나, 다음 조건을 6 개월 이상
만족한 contributor:

- 머지된 PR ≥ 10 건 (라이브러리 특성상 PR 빈도가 낮으므로 일반 operator 의
  절반 기준).
- 리뷰한 PR ≥ 20 건 (downstream consumer 측 PR 포함 가능).
- `pkg/` 하위 패키지 중 하나 이상의 영역 (security / labels / webhook /
  monitoring / networkpolicy / version / status / finalizer / storageclass /
  events / probes / pvc / topology) 에 대한 깊은 이해.

## 추가 절차

1. 기존 메인테이너 또는 candidate 본인이 Issue 또는 ADR 로 제안.
2. `@keiailab/maintainers` 팀의 lazy consensus (7일 코멘트 윈도우).
3. 반대 없으면 GitHub team 에 추가 + 본 파일 PR 갱신.

## 비활성 메인테이너

연속 6 개월간 활동이 없는 메인테이너는 emeritus 로 이동합니다 (권한 회수,
명예 명단 유지).

## Cross-repo 합의

본 라이브러리의 *공개 API breaking change* 는 downstream consumer
메인테이너의 LGTM 이 ADR 단계에서 동반되어야 합니다 — [GOVERNANCE.md](GOVERNANCE.md)
참조.

## 다국어 문서 책임자 (i18n owner)

| 언어 | 책임자 | 담당 파일 | 책임 |
|---|---|---|---|
| English (canonical) | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | `README.md` 외 canonical 문서 | 진본 (single source of truth) |
| 한국어 (Korean) | TaeHwan Park ([@eightynine01](https://github.com/eightynine01)) | `README.ko.md` 외 ko 번역 | EN canonical sync + 번역 검수 |
| 日本語 (Japanese) | (모집 중 — Issue 로 자원) | `*.ja.md` | AI 번역 native review |
| 中文 (Chinese) | (모집 중 — Issue 로 자원) | `*.zh.md` | AI 번역 native review |

**drift 검증**: `bash scripts/check-readme-sync.sh` — 양 file 존재 + section
header 수 일치 + line count ±15 % + 양방향 cross-link. lefthook pre-push
`readme-i18n-sync` hook 가 자동 강제합니다.

## Emeritus

(아직 없음)

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
