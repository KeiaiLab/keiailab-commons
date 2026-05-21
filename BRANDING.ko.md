# 브랜드 가이드 — `operator-commons`

> ⚠️ This translation is AI-generated and pending native review. — 본 번역은 Claude 기계 번역 결과입니다.

> keiailab 오퍼레이터 패밀리의 visual identity, voice, tone.

본 문서는 `operator-commons` 브랜딩 결정의 canonical reference 입니다. README, release note, 마케팅 자료, 프로젝트를 대표하는 third-party 커뮤니케이션에 적용됩니다.

## 1. Identity

**Organization**: [keiailab](https://keiailab.com) — Kubernetes-native 데이터 플랫폼 오퍼레이터 (Apache-2.0, license-clean, vanilla-upstream 호환).

**Project**: `operator-commons` — keiailab 오퍼레이터의 공유 Go 라이브러리 — finalizer / labels / status / version / security / monitoring partial.

**Family**: [`operator-commons`](https://github.com/keiailab/operator-commons) 공유 라이브러리를 공유하는 4 자매 오퍼레이터 중 하나:

| 프로젝트 | 데이터베이스 | 저장소 |
|---|---|---|
| `postgres-operator` | PostgreSQL 18+ | https://github.com/keiailab/postgres-operator |
| `mongodb-operator` | MongoDB 7.0+ | https://github.com/keiailab/mongodb-operator |
| `valkey-operator` | Valkey 8.0+ (Redis fork, BSD-3) | https://github.com/keiailab/valkey-operator |
| `operator-commons` | 공유 Go 라이브러리 | https://github.com/keiailab/operator-commons |

## 2. 로고 & 시각 자산

| 자산 | URL | 사용처 |
|---|---|---|
| Primary 로고 (SVG) | `https://keiailab.com/assets/logo.svg` | README header, slide |
| Mono mark | `https://keiailab.com/assets/mark.svg` | Favicon, social card |
| Wordmark | `https://keiailab.com/assets/wordmark.svg` | Footer, dark background |

**로고 배치**: README 상단 중앙, width 120px. 항상 https://keiailab.com 으로 링크.

**Clear space**: 로고 주변 최소 padding = 로고 width 의 25%.

**금지**:
- 로고 색상 변경
- drop shadow / filter 추가
- 대비 부족 배경에 배치
- keiailab 브랜드 승인 없이 타 로고와 결합

## 3. 색상 팔레트

| 역할 | Hex | 사용처 |
|---|---|---|
| Primary (keiailab teal) | `#0EA5A8` | 헤더, primary action, 링크 |
| Secondary (deep navy) | `#0F172A` | dark 배경, 코드 블록 |
| Accent (warm amber) | `#F59E0B` | 강조, 배지 accent |
| Neutral grey | `#64748B` | light 배경의 body text |
| Background light | `#F8FAFC` | 문서 페이지 배경 |
| Background dark | `#020617` | 코드 에디터 테마, dark mode |

GitHub README 의 shield.io badge 는 위 hex 사용 권장.

## 4. 타이포그래피

- **Heading**: 시스템 기본 (GitHub 의 default `-apple-system, BlinkMacSystemFont, Segoe UI, ...`)
- **Body**: 동일 (GitHub-native 정합)
- **Code**: `ui-monospace, SFMono-Regular, Consolas, ...`

별도 webfont 미사용 (GitHub README rendering 정합).

## 5. Voice & Tone

**Audience**: Kubernetes 플랫폼 엔지니어 / DBA / SRE / Go 라이브러리 consumer.

**Voice 원칙**:
- **Direct** — 가능하면 문단보다 bullet point
- **Evidence-based** — claim 은 benchmark / SLA / link 포함
- **Library-focused** — `operator-commons` 는 *라이브러리* 다 — controller-runtime / CRD / reconciler 는 *consumer operator 의 책임*
- **License-aware** — Apache-2.0 only, AGPL/BUSL transitive 의존성 0건 목표 (ADR-0001 charter)

**Avoid**:
- 마케팅 최상급 ("blazing fast", "revolutionary", "best-in-class")
- 모호한 비교 ("X-class quality") — *구체적 metric 또는 benchmark 로 qualify*
- 로드맵의 시간 기반 마감일 (`standards/roadmap.md §1.1` — 기능 체크리스트 대신)

## 6. README Header 표준

모든 README 의 첫 문단은 다음 형식 (Wave 3 표준):

```markdown
<p align="center">
  <img src="https://keiailab.com/assets/logo.svg" alt="keiailab" width="120"/>
</p>

# operator-commons

> **keiailab 오퍼레이터의 공유 Go 라이브러리 — finalizer / labels / status / version / security / monitoring partial**

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"/></a>
</p>

<p align="center">
  <a href="README.md">English</a> |
  <b>한국어</b> |
  <a href="README.ja.md">日本語</a> |
  <a href="README.zh.md">中文</a>
</p>
```

## 7. README Footer 표준

모든 README + root-level .md 파일의 마지막에 다음 footer 부착 (Wave 3 표준):

```markdown
---

<p align="center">
  <b>keiailab 오퍼레이터 패밀리</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
```

## 8. Badge 표준 순서

README 의 shield.io badge 순서 (좌→우):

1. License (Apache-2.0)
2. Go Version (1.25+)
3. Go Reference (pkg.go.dev)
4. OpenSSF Scorecard
5. GitHub Discussions

> **Note**: `operator-commons` 는 *라이브러리* 라 Kubernetes / Container Image / Helm Chart badge 는 *consumer operator 측 README 에 위치*. 본 라이브러리는 `pkg.go.dev` + `OpenSSF Scorecard` 중심.

## 9. Discussions / Issues / PR Template

- **Discussions**: `https://github.com/keiailab/operator-commons/discussions` — pkg API 질문, integration 사례, 새 helper 제안
- **Issues**: bug report + use case 가 있는 구체적 feature request (consumer operator 측 영향 명시 권장)
- **PR template**: `.github/PULL_REQUEST_TEMPLATE.md` 표준 (사용자 시나리오 + 검증 명령 인용 의무, `standards/checklist.md §3`)

## 10. Social & External

- **Website**: https://keiailab.com
- **GitHub Org**: https://github.com/keiailab
- **pkg.go.dev**: https://pkg.go.dev/github.com/keiailab/operator-commons

## 11. License & Attribution

- License: [Apache-2.0](LICENSE)
- Copyright: © 2026 keiailab contributors
- Third-party attribution: [NOTICE](NOTICE) 참조 (해당 시)

---

<p align="center">
  <b>keiailab 오퍼레이터 패밀리</b><br/>
  <a href="https://github.com/keiailab/postgres-operator">postgres-operator</a> ·
  <a href="https://github.com/keiailab/mongodb-operator">mongodb-operator</a> ·
  <a href="https://github.com/keiailab/valkey-operator">valkey-operator</a> ·
  <a href="https://github.com/keiailab/operator-commons">operator-commons</a>
</p>

<p align="center">
  © 2026 keiailab · <a href="LICENSE">Apache-2.0</a> · <a href="https://keiailab.com">keiailab.com</a>
</p>
