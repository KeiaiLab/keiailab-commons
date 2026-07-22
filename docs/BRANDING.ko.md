# 브랜드 가이드 — `keiailab-commons`

> [English](BRANDING.md) | **한국어** | [日本語](BRANDING.ja.md) | [中文](BRANDING.zh.md)

> `keiailab-commons` 라이브러리의 visual identity, voice, tone.

본 문서는 `keiailab-commons` 브랜딩 결정의 canonical reference 입니다.
README, release note, 프로젝트에 관한 외부 커뮤니케이션에 적용됩니다.

## 1. Identity

**Organization**: [keiailab](https://keiailab.com).

**Project**: `keiailab-commons` — Kubernetes operator 공통 scaffolding
(finalizer / labels / status / version / security / monitoring partial) 을
위한 Go 라이브러리입니다.

본 라이브러리는 Go 모듈 `github.com/keiailab/keiailab-commons` 와
Helm library chart (`charts/keiailab-commons`) 로 배포됩니다. downstream
operator 가 표준 Go module import 로 사용합니다 — 특정 consumer 를
지명하거나 endorsement 하지 않습니다.

### 1.1 패밀리 charter (본 문서가 SSOT)

본 Branding Guide 는 keiailab operator 패밀리가 공유하는 시각 문법의
**단일 진본(SSOT)** 입니다. 색상 팔레트 (§3), 링 그래머 및 글리프 스왑 절차
(§2b), repo 별 글리프 레지스트리 (§2b), README 헤더 / footer 템플릿 (§6 / §7),
배지 순서 (§8) 는 **여기서** 정의하고 각 패밀리 저장소가 재유도하지 않고
복사합니다. `keiailab-commons` 가 charter 를 소유하는 이유는 모든 operator 가
이미 import 하는 공유 의존성이기 때문입니다.

패밀리는 4 개의 sister operator 와 본 공유 라이브러리로 구성됩니다:

| Project | Focus | Repository |
|---|---|---|
| `mongodb-operator` | MongoDB 8.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ | https://github.com/keiailab/valkey-operator |
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `qdrant-operator` | Qdrant vector database | https://github.com/keiailab/qdrant-operator |
| `keiailab-commons` | Shared Go library | https://github.com/keiailab/keiailab-commons |

본 가이드의 번역이 뒤처질 경우, 여기의 영어 텍스트가 canonical 로 유지됩니다.

## 2. 로고 및 시각 자산

| 자산 | URL | 사용처 |
|---|---|---|
| Primary 로고 | `docs/branding/symbol.png` | README 헤더, 슬라이드 |
| Keiailab base symbol | `docs/branding/base-symbol.png` | 외곽 회전 화살표 마크의 소스 레퍼런스 |
| Favicon | `https://keiailab.com/favicon.ico` | Favicon, social card |
| 예정 SVG kit | `https://keiailab.com/assets/{logo,mark,wordmark}.svg` | URL 이 200 반환 후 교체 예정 |

**로고 배치**: README 상단 중앙, width 96 px. 항상 `https://keiailab.com` 으로 링크.

**Clear space**: 로고 주위 최소 padding 은 로고 width 의 25 % 입니다.

**금지**:

- 로고 색상 변경
- drop shadow / filter 추가
- 대비 부족한 배경에 배치
- keiailab 브랜드 승인 없이 다른 로고와 결합

## 2b. 링 그래머 & 글리프 스왑 절차

패밀리는 하나의 외곽 심볼 — `docs/branding/base-symbol.png` 의 3색 회전 화살표
링 (460×460; 다크네이비 → 블루 → 그린 화살표, 중앙 구체) — 을 공유합니다.
**모든** 패밀리 저장소는 이 파일을 byte 단위로 동일하게 vendor 하며, 링은 절대
재색·재묘사·재생성하지 않습니다.

각 프로젝트의 `docs/branding/symbol.png` 는 링을 유지하고 **중앙 글리프만**
교체하여 만듭니다:

1. 공유 `base-symbol.png` 에서 시작 — 링은 건드리지 않습니다.
2. 제품 글리프를 중앙 clear zone (~185 px 지름) 에 alpha-composite — 링과
   겹치지 않습니다.
3. 소스 해상도로 `docs/branding/symbol.png` 로 export.
4. README 헤더에서 width 96 px 로 참조 (§6), operator 는 Helm chart `icon` 으로도
   참조.

공유 링 + 글리프 교체가 마크를 하나의 패밀리로 읽히게 하면서도 제품별 식별성을
유지하는 핵심입니다. 링 재색 또는 새 외곽 마크 작도는 금지입니다.

### 글리프 레지스트리

| Repository | Center glyph |
|---|---|
| `mongodb-operator` | leaf |
| `valkey-operator` | hexagon + keyhole |
| `postgres-operator` | elephant |
| `qdrant-operator` | kNN constellation |
| `keiailab-commons` | node grid |

## 3. 색상 팔레트

| 역할 | Hex | 사용처 |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | 헤더, primary action, 링크 |
| Secondary (deep navy) | `#0F172A` | dark 배경, 코드 블록 |
| Accent (warm amber) | `#F59E0B` | 강조, 배지 accent |
| Neutral grey | `#64748B` | light 배경의 body text |
| Background light | `#F8FAFC` | 문서 페이지 배경 |
| Background dark | `#020617` | dark-mode 코드 에디터 테마 |

GitHub README shield.io 배지는 같은 hex 값을 사용합니다.

## 4. 타이포그래피

- **Heading**: 시스템 기본 (GitHub 기본 `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: 시스템 기본 (GitHub-native 정합)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...` (GitHub 기본 monospace)

별도 web font 미사용 — GitHub native rendering 그대로 유지합니다.

## 5. Voice & Tone

**Audience**: Kubernetes 플랫폼 엔지니어, DBA, SRE, Go 라이브러리 consumer.

**Voice 원칙**:

- **Direct** — 가능하면 문단보다 bullet point 를 사용합니다.
- **Evidence-based** — 모든 claim 은 benchmark, SLA, 또는 link 를 동반합니다.
- **Library-focused** — `keiailab-commons` 는 *라이브러리* 입니다.
  controller-runtime, CRD, reconciler 는 downstream consumer 의 책임이며
  본 라이브러리의 책임이 아닙니다.
- **License-aware** — MIT only. charter 목표는 AGPL / BUSL transitive
  의존성 0 건입니다 (`docs/kb/adr/0001-charter.md` 참조).

**Avoid**:

- 마케팅 최상급 ("blazing fast", "revolutionary", "best-in-class").
- 모호한 비교 ("enterprise-grade quality") — 각 claim 을 구체적 metric 또는
  benchmark 로 qualify 합니다.
- 로드맵의 시간 기반 마감 — [ROADMAP.md](ROADMAP.md) 의 기능 체크리스트
  사용.

## 6. README Header 표준

모든 README 의 첫 블록은 다음 형식을 따릅니다:

```markdown
<p align="center">
  <img src="docs/branding/symbol.png" alt="keiailab" width="96"/>
</p>

# keiailab-commons

> **Kubernetes operator 공통 scaffolding 을 위한 Go 라이브러리 — finalizer / labels / status / version / security / monitoring partials.**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/<license placeholder>-blue.svg" alt="License"/></a>
  <!-- §8 순서에 따른 추가 shield.io 배지 -->
</p>

<p align="center">
  <a href="README.md">English</a> |
  <b>한국어</b> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>
```

> **템플릿 placeholder**: `<license placeholder>` 를 복사하는 저장소의 실제 SPDX
> 배지 세그먼트로 교체합니다 — 예: MIT 저장소(패밀리 대부분)는 `License-MIT`,
> 라이선스가 다른 멤버는 해당 `License-<SPDX-id>`. 공유 템플릿에 라이선스를
> 하드코딩하면 라이선스가 다른 멤버에 오배지가 재생산되므로, 토큰은 여기서
> generic 으로 두고 저장소별로 채웁니다.

## 7. README Footer 표준

모든 README 와 root-level `.md` 파일은 다음 단일 attribution 라인으로
끝납니다:

```markdown
---

<p align="center">© 2026 keiailab · <license placeholder> · <a href="https://keiailab.com">keiailab.com</a></p>
```

> **템플릿 placeholder**: `<license placeholder>` 를 복사하는 저장소의 실제
> 라이선스 이름으로 교체합니다 — §6 배지 토큰과 동일한 저장소별 치환
> 규칙입니다. `keiailab-commons` 자신의 실값: MIT.

추가 cross-link 블록은 두지 않습니다. footer 를 최소화하여 문서를
self-contained 로 유지합니다.

## 8. 배지 순서

배지는 **MUST** (항상 표시) 와 **SHOULD** (백엔드 인프라가 라이브일 때만 표시)
로 나뉩니다. 워크플로나 레지스트리가 아직 없는데 8 개를 강제하면 배지가 상시
red 로 남으므로, 패밀리 표준은 정직한 4 + 4 입니다.

**MUST (4)** — 모든 패밀리 저장소, 좌→우:

1. License
2. Go Version
3. Product (MongoDB / Valkey / PostgreSQL / Qdrant)
4. Kubernetes Version

**SHOULD (4)** — 백엔드 표면이 활성일 때만 각각 추가:

5. Container image (`ghcr.io/keiailab/<repo>`)
6. Helm chart (Artifact Hub)
7. OpenSSF Scorecard
8. GitHub Discussions

> **라이브러리 예외**: `keiailab-commons` 는 container image, Helm application
> chart, Kubernetes 워크로드 어느 것도 배포하지 않으므로, MUST 세트는 product 및
> Kubernetes 배지 대신 License / Go Version / **Go Reference** (pkg.go.dev) 이고,
> SHOULD 세트는 OpenSSF Scorecard + GitHub Discussions 입니다. container / Helm /
> Kubernetes 배지는 이미지나 chart 를 배포하는 downstream operator 에 둡니다.

## 9. Discussions / Issues / PR Template

- **Discussions**: `https://github.com/keiailab/keiailab-commons/discussions` — 패키지 API 질문, integration 사례, 새 helper 제안.
- **Issues**: 버그 보고 및 use case 가 있는 구체적 feature request. 관련 시
  downstream consumer 영향을 명시합니다.
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` — Conventional Commits +
  사용자 시나리오 + 검증 명령 출력 인용.

## 10. Social & External

- **Website**: <https://keiailab.com>
- **GitHub Org**: <https://github.com/keiailab>
- **pkg.go.dev**: <https://pkg.go.dev/github.com/keiailab/keiailab-commons>

## 11. License & Attribution

- License: [MIT](../LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: [NOTICE](../NOTICE) 참조

---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
