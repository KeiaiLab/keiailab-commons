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

## 2. 로고 및 시각 자산

| 자산 | URL | 사용처 |
|---|---|---|
| Current primary logo | `docs/branding/symbol.png` | README header, slide decks |
| Current favicon | `https://keiailab.com/favicon.ico` | Favicon, social cards |
| Planned SVG kit | `https://keiailab.com/assets/{logo,mark,wordmark}.svg` | Future replacement after URLs return 200 |

**로고 배치**: README 상단 중앙, width 96 px. 항상 `https://keiailab.com` 으로 링크.

**Clear space**: 로고 주위 최소 padding 은 로고 width 의 25 % 입니다.

**금지**:

- 로고 색상 변경
- drop shadow / filter 추가
- 대비 부족한 배경에 배치
- keiailab 브랜드 승인 없이 다른 로고와 결합

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
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-0EA5A8.svg" alt="License"/></a>
  <!-- 추가 shield.io 배지 -->
</p>

<p align="center">
  <a href="README.md">English</a> |
  <b>한국어</b> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>
```

## 7. README Footer 표준

모든 README 와 root-level `.md` 파일은 다음 단일 attribution 라인으로
끝납니다:

```markdown
---

<p align="center">© 2026 keiailab · MIT · <a href="https://keiailab.com">keiailab.com</a></p>
```

추가 cross-link 블록은 두지 않습니다. footer 를 최소화하여 문서를
self-contained 로 유지합니다.

## 8. 배지 순서

README 의 shield.io 배지는 다음 순서 (좌→우):

1. License (MIT)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `keiailab-commons` 는 *라이브러리* 이므로 container image,
> Helm chart, Kubernetes deployment 배지는 본 라이브러리에 부착하지 않습니다 —
> 이미지나 chart 를 배포하는 downstream operator README 에 둡니다.

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
