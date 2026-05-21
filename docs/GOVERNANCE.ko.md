# Governance

> [English](GOVERNANCE.md) | **한국어** | [日本語](GOVERNANCE.ja.md)

본 문서는 `keiailab/operator-commons` 라이브러리의 의사결정 절차를 정의합니다.
본 라이브러리는 downstream consumer operator 들이 공통으로 import 하므로,
*공개 API* 변경은 downstream 호환성에 영향을 줍니다.

## 원칙

1. **개방성** — 모든 의사결정은 공개 채널 (GitHub Issue / PR / ADR) 에서
   진행됩니다.
2. **최소 합의 (Lazy Consensus)** — 일상적 변경은 반대가 없으면 진행됩니다.
3. **명시적 합의 (Explicit Consensus)** — 공개 API breaking change, 새 패키지
   도입, 라이선스 변경은 ADR 후 메인테이너 **2/3 supermajority** 승인 +
   downstream consumer 측 1명 이상의 LGTM 이 동시에 필요합니다.
4. **공동 책임** — 메인테이너는 라이브러리 안정성, downstream 운영 영향,
   보안에 대해 공동 책임을 집니다.

## 의사결정 분류

### 일상 변경 (Lazy Consensus)

- 버그 픽스, 문서 개선, 테스트 추가, minor / patch 의존성 업그레이드, 내부
  리팩터링 (공개 API 무변경).
- 절차: PR → 메인테이너 1명 이상 LGTM → 머지.
- 시한: 별도 윈도우 없음. 로컬 게이트 (pre-commit / pre-push hook + Makefile)
  통과 시점에 머지 가능합니다. GitHub Actions 는 사용하지 않으며 — 모든
  품질 게이트는 로컬 4 계층 (`lefthook.yml`, `Makefile`, 리뷰어 증거,
  ADR coverage) 으로 강제합니다.

### 중간 변경 (Explicit Consensus)

- 새 공개 API 함수 / 타입 추가, 의존성 major 업그레이드, 새 `pkg/<sub>` 패키지
  도입.
- 절차: Issue 또는 ADR 로 제안 → 7일 코멘트 윈도우 → 메인테이너 다수 LGTM →
  머지.
- 거부 의견 1건 이상 시 메인테이너 회의로 토의 진행.

### 공개 API breaking change (ADR 필수)

- 함수 시그니처 변경, 타입 제거, 모듈 경로 변경, 라이선스 변경.
- 절차:
  1. `docs/kb/adr/NNNN-<slug>.md` 양식으로 ADR 제출.
  2. 14일 코멘트 윈도우.
  3. 메인테이너 2/3 + downstream consumer 측 1명 이상 LGTM.
  4. ADR `Status: Draft → Accepted` 후 구현 PR 머지.

## 보안 결정

CVE 보고는 [SECURITY.ko.md](../SECURITY.ko.md) 절차에 따라 비공개로 우선 처리합니다.
downstream consumer 측 fix release 가능 시점까지 embargo 를 유지합니다.

## 릴리스 결정

- **v0.x**: 메인테이너 1 인 lazy consensus 로 minor / patch tag 가능.
- **v1.0+ (stable)**: SemVer 엄격 — major bump 는 ADR + 2/3 supermajority.

## 변경 이력

| Date | Change |
|---|---|
| 2026-05-09 | 본 문서 신설 — 거버넌스 자산 정합. |

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
