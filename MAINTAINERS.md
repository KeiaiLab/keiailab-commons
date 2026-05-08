# Maintainers

본 문서는 `keiailab/operator-commons` 의 의사결정 권한을 가진 메인테이너 명단을 관리합니다.

## 현재 메인테이너

| 이름/팀 | GitHub | 역할 | 담당 영역 |
|---|---|---|---|
| keiailab maintainers | [@keiailab/maintainers](https://github.com/orgs/keiailab/teams/maintainers) | Lead | 전체 |

GitHub team `@keiailab/maintainers` 가 본 라이브러리의 모든 영역에 대한 머지·release tag 권한을 보유합니다.

## 메인테이너 자격

3 consumer operator (`mongodb-operator`, `postgres-operator`, `valkey-operator`) 중 어느 한 쪽 메인테이너이거나, 다음 조건을 6개월 이상 만족한 contributor:

- 머지된 PR ≥ 10건 (라이브러리 특성상 PR 빈도가 낮으므로 operator repo 절반 기준)
- 리뷰한 PR ≥ 20건 (consumer operator 측 PR 포함 가능)
- `pkg/` 하위 패키지 중 한 영역(security / labels / webhook / monitoring / networkpolicy / version) 이상 깊은 이해

## 추가 절차

1. 기존 메인테이너 또는 candidate 본인이 issue 또는 ADR 로 제안
2. `@keiailab/maintainers` 팀의 lazy consensus (7일 코멘트 윈도우)
3. 반대 없으면 GitHub team 추가 + 본 파일 PR 갱신

## 비활성 메인테이너

연속 6개월간 활동이 없는 메인테이너는 emeritus 로 이동합니다 (권한 회수, 명예 명단 유지).

## Cross-repo 합의

본 라이브러리의 *공개 API breaking change* 는 3 consumer operator 메인테이너 (각 repo에서 1명 이상) 의 LGTM 이 동시에 필요합니다 — RFC 로 명시.

## Emeritus

(아직 없음)
