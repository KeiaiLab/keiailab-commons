# Contributing to operator-commons

> [English](CONTRIBUTING.md) | **한국어** | [日本語](CONTRIBUTING.ja.md) | [中文](CONTRIBUTING.zh.md)

`keiailab/operator-commons` 는 downstream Kubernetes operator 가 공통으로
import 하는 Go 라이브러리입니다. 모든 기여는 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
와 [docs/GOVERNANCE.ko.md](docs/GOVERNANCE.ko.md) 를 따릅니다.

## 기여 흐름

1. **Issue 또는 ADR 로 의도 공유** (큰 변경의 경우).
2. **Fork & feature branch** — `feat/<slug>`, `fix/<slug>`, `docs/<slug>`,
   `refactor/<slug>` 형식.
3. **로컬 게이트 통과 확인**:
   - `lefthook install --force`
   - `lefthook run pre-commit --all-files`
   - `make lint test`
4. **PR 작성** — Conventional Commits + 한국어 본문 허용.
5. **리뷰 SLA**: 메인테이너 24시간 이내 1차 응답.

## PR 체크리스트 (작성자)

- [ ] PR 제목: Conventional Commits 형식 (`feat`, `fix`, `docs`, `refactor`,
  `test`, `chore`).
- [ ] PR 본문: 변경 요약 + 검증 명령 + 출력 인용.
- [ ] 단위 테스트 추가 또는 갱신 (`pkg/<sub>` 변경 시 의무).
- [ ] 공개 API 변경 시 GoDoc 갱신.
- [ ] **공개 API breaking** 시 ADR 링크 + downstream consumer 측 영향 분석.
- [ ] `go.mod` / `go.sum` drift 없음 (`go mod tidy` 후 변경 없음).
- [ ] 의존성 추가 시 라이선스 검증 인용 + CVE 검토.

## 로컬 개발 (downstream consumer 와의 동시 변경)

downstream operator 와 commons 를 동시에 수정해야 하는 *cross-cut* 변경:

```fish
# 1. consumer operator 의 go.mod 에 replace directive 추가 (커밋 금지, 로컬 한정)
# go.mod 끝줄에:
#   replace github.com/keiailab/operator-commons => ../operator-commons

# 2. 양쪽 동시 수정 + go test ./... 양쪽 모두 PASS

# 3. PR 분리:
#    - operator-commons 측 PR 머지 + tag (예: v0.9.0)
#    - consumer operator 측 PR 에서 require 버전 bump (replace 제거)
```

## 새 `pkg/<sub>` 패키지 추가

[docs/GOVERNANCE.ko.md](docs/GOVERNANCE.ko.md) "중간 변경" 절차 적용:

1. Issue 또는 ADR 로 제안 — *왜* commons 에 들어가야 하는가, *어떤* downstream
   consumer 가 사용할 것인가.
2. 7일 코멘트 윈도우.
3. 다수 LGTM 후 PR 머지.

## 릴리스

- v0.x SemVer: 매 minor 가 *공개 API* 또는 *의미 있는 동작* 변경.
- 릴리스 절차:
  1. `git tag v0.X.Y` (annotated).
  2. `git-cliff` 로 CHANGELOG 갱신 PR.
  3. `git push origin v0.X.Y`.
  4. downstream consumer 측 별도 PR 로 require 버전 bump.

## 보안 취약점

[SECURITY.ko.md](SECURITY.ko.md) 절차. 공개 issue 금지.

---

<p align="center">© 2026 keiailab · Apache-2.0 · <a href="https://keiailab.com">keiailab.com</a></p>
