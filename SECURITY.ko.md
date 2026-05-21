# Security Policy

> [English](SECURITY.md) | **한국어** | [日本語](SECURITY.ja.md) | [中文](SECURITY.zh.md)

`keiailab/operator-commons` 는 downstream Kubernetes operator 가 import 하는
라이브러리이므로, 본 라이브러리의 취약점은 *downstream operator 의 운영
보안* 에 직접 영향을 줄 수 있습니다.

## 취약점 보고

**공개 issue 로 보고하지 마세요.**

### 보고 경로

다음 중 하나로 비공개 보고:

1. **GitHub Security Advisory** (권장):
   `https://github.com/keiailab/operator-commons/security/advisories/new`
2. **이메일**: `security@keiailab.com` (PGP 옵션):
   - PGP fingerprint: `89A4 0947 6828 CB99 2338  C378 651E 51AF 520B CB78`.

### 포함 정보

- 영향받는 버전 (release tag 또는 commit SHA).
- 영향받는 패키지 (`pkg/security`, `pkg/webhook` 등).
- 재현 단계 (가능한 한 minimal repro; downstream 환경 의존 시 명시).
- 영향 평가 — downstream consumer 에 미치는 영향 범위.
- CVSS 자체 평가 시 포함.

## 응답 SLA

| 단계 | 시간 |
|---|---|
| 초기 응답 (수신 확인) | 72시간 이내 |
| 심각도 평가 | 7일 이내 |
| 패치 release | severity 따라 (Critical: 14일, High: 30일, Medium: 60일) |
| 공개 disclosure | 패치 release 후 14일 (downstream consumer 동시 fix release 가능 시점) |

## Embargo 처리

공개 API 가 영향받는 취약점은 downstream consumer 측에 동시 fix release
가능 시점까지 embargo. consumer 측 메인테이너에게 비공개 advisory 사전 공유.

## 지원 버전

| Version | Supported |
|---------|-----------|
| 0.x (alpha) | ✅ 최신 minor 만 |
| 1.0+ (stable) | TBD — 첫 stable release 후 갱신 |

현재 v0.x 단계. *공개 API breaking 가능* — 보안 패치는 *최신* minor 에만
적용됩니다.

## 의존성 보안

본 라이브러리 의존성 추가 / 업그레이드 시 라이선스 검증 + CVE 검토 결과를
PR 본문에 인용합니다. Dependabot / Renovate 자동 업데이트 PR 은 우선 review.

## 라이선스 / 공급망

본 라이브러리는 **Apache-2.0 only** 정책 — *AGPL / BUSL transitive 의존성 0건*
목표 (charter ADR-0001). 매 minor release 시 license 감사.

## 보안 모범 사례 (downstream consumer 측)

본 라이브러리를 import 하는 operator 는 다음을 보장해야 합니다:

1. **`pkg/security` 사용** — PodSecurity restricted SecurityContext 빌더
   직접 호출 (자체 구현 금지).
2. **`pkg/webhook` 사용** — 버전 validation 직접 구현 금지.
3. **`pkg/networkpolicy` 사용** — deny-by-default NetworkPolicy 빌더 활용.
4. 의존성 버전 동기화 — `go.mod` 의 `github.com/keiailab/operator-commons`
   는 *최신 patch* 항상 추적 (Renovate 자동 PR).

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
