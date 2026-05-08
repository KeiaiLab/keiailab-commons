# Governance

본 문서는 `keiailab/operator-commons` 라이브러리의 의사결정 절차를 정의합니다. 본 라이브러리는 3 consumer operator (`mongodb-operator`, `postgres-operator`, `valkey-operator`) 가 공통으로 import 하므로, *공개 API* 변경은 3 repo 전반의 호환성 영향을 고려해야 합니다.

## 원칙

1. **개방성**: 모든 의사결정은 공개 채널(GitHub issue/PR/RFC)에서 이뤄집니다.
2. **최소 합의(Lazy Consensus)**: 일상적 변경은 반대 없으면 진행됩니다.
3. **명시적 합의(Explicit Consensus)**: 공개 API breaking, 새 패키지 도입, 라이선스 변경은 RFC 후 메인테이너 **2/3 supermajority** 승인 + 3 consumer operator 메인테이너 LGTM 동시 필요.
4. **공동 책임**: 메인테이너는 라이브러리 안정성, downstream 운영 영향, 보안에 대해 공동 책임을 집니다.

## 의사결정 분류

### 일상 변경 (Lazy Consensus)
- 버그 픽스, 문서 개선, 테스트 추가, 의존성 minor/patch 업그레이드, 내부 리팩터링(공개 API 무변경)
- 절차: PR → 1명 이상 메인테이너 LGTM → 머지
- 시한: 별도 윈도우 없음 (CI 그린이면 즉시 머지 가능)

### 중간 변경 (Explicit Consensus)
- 새 공개 API 함수/타입 추가, 의존성 major 업그레이드, 새 `pkg/<sub>` 패키지 도입
- 절차: 이슈로 제안 → 7일 코멘트 윈도우 → 메인테이너 다수 LGTM → 머지
- 거부 1건이 있을 시 메인테이너 회의에서 토론

### 공개 API breaking (RFC 필수)
- 함수 시그니처 변경, 타입 제거, 모듈 경로 변경, 라이선스 변경
- 절차:
  1. `docs/kb/adr/NNNN-title.md` 또는 ai-dev `rfcs/NNNN-title.md` 제출
  2. 14일 코멘트 윈도우
  3. 메인테이너 2/3 + 3 consumer operator 측 1명 이상씩 LGTM
  4. ADR/RFC `Status: Draft → Accepted` 후 구현 PR

## 보안 결정

CVE 보고는 [SECURITY.md](SECURITY.md) 절차에 따라 비공개로 우선 처리. consumer operator 측 fix release 가능 시점까지 embargo 유지.

## 릴리스 결정

- v0.x: 메인테이너 1 인 lazy consensus 로 minor / patch tag 가능
- v1.0+ (stable): SemVer 엄격 — major bump 는 RFC + 2/3 supermajority

## 변경 이력

| Date | Change | Refs |
|---|---|---|
| 2026-05-09 | 본 문서 신설 — 4-repo 거버넌스 자산 정합 | plan: mongodb-operator-operator-commons-postgr-tranquil-horizon |
