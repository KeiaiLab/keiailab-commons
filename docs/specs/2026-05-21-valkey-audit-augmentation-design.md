# Valkey-operator Audit ❌ 5건 보강 Sub-spec

| 메타 | 값 |
|---|---|
| 날짜 | 2026-05-21 |
| 상태 | Proposed (ralph-loop 또는 후속 cycle 진입 대기) |
| 범위 | valkey-operator 만 (postgres / mongodb / commons / forgewise 별 spec 분리) |
| 책임 | ralph-loop (사용자 결정 #2 "관조") 또는 사용자 명시 결정 후 본 thread / 별 subagent |

## 1. 배경

본 session 의 audit 진척 (38 → 7) 후 valkey-operator 만의 잔여 ❌ 5건:

| ID | 항목 | 현 상태 | 보강 방법 |
|---|---|---|---|
| P1-11 | kube-linter hook | `.lefthook.yml` 에 없음 | postgres-operator `.lefthook.yml` 의 kube-linter 블록 cp |
| P1-12 | go-licenses hook | 없음 | postgres-operator 의 go-licenses 블록 cp |
| P1-13 | markdown-link-check hook | 없음 | postgres-operator 의 markdown-link-check 블록 cp |
| OP-2 | scripts/helm-publish.sh | 없음 | postgres-operator/scripts/helm-publish.sh cp + valkey-specific 조정 |
| OP-10 | docs/UPGRADING.md | 없음 | postgres-operator/docs/UPGRADING.md 패턴 cp + valkey-specific 본문 |

## 2. Goals

| ID | 목표 | 검증 |
|---|---|---|
| G1 | valkey `.lefthook.yml` pre-push 에 kube-linter + go-licenses + markdown-link-check 3종 추가 | `lefthook run pre-push` 통과 (3 hook ✅ 또는 graceful skip) |
| G2 | valkey `scripts/helm-publish.sh` 신규 (postgres 패턴) | `bash -n scripts/helm-publish.sh` syntax OK |
| G3 | valkey `docs/UPGRADING.md` 신규 (postgres 패턴) | `ls docs/UPGRADING.md` 통과 |
| G4 | ADR (Status: Accepted) — audit 5 ❌ → ✅ 적용 결정 | `docs/kb/adr/0049-audit-augmentation.md` |

## 3. 아키텍처

```
[Phase 1] lefthook 3종 보강 — postgres .lefthook.yml 패턴 cp
   ↓
[Phase 2] scripts/helm-publish.sh — postgres 패턴 cp + valkey chart name
   ↓
[Phase 3] docs/UPGRADING.md — postgres 패턴 cp + valkey-specific 본문
   ↓
[Phase 4] ADR-0049 + INDEX 갱신
   ↓
[Phase 5] audit 재실행 — valkey 5 ❌ → ✅ 확인
```

## 4. 단계별 상세

### Phase 1 — lefthook 3종 보강

`postgres-operator/.lefthook.yml` 의 `pre-push.commands` 의 3 block 을 valkey-operator/.lefthook.yml 의 같은 위치에 cp. 정합 keyword:
- kube-linter (chart sanity)
- go-licenses (disallowed forbidden/restricted)
- markdown-link-check (변경 *.md only)

각 hook 의 graceful skip 정책 유지 (도구 미설치 시 skip + warn).

### Phase 2 — scripts/helm-publish.sh

`postgres-operator/scripts/helm-publish.sh` cp + 다음 변경:
- chart 이름: `valkey-operator` (postgres-operator 가 아님)
- 기본 빌더 (CLAUDE.md §2 RFC-0002): docker buildx default + amd64 단일
- gh-pages branch 동기 + helm repo index

### Phase 3 — docs/UPGRADING.md

`postgres-operator/docs/UPGRADING.md` 패턴 cp + valkey-specific 내용:
- §0 semver 정책 (v0.x → v1.0 valkey 자체 버전)
- §1 v0.x → 현재 (valkey ValkeyCluster CRD 변경 + Helm chart 변경)
- §2 후속 Sprint 1 + S5 채택 (commons v0.9.0)
- §3 v1.0 v3.x-stable 선언 시점
- §4 GHA dual-track 정책 (valkey ADR/0048 v2.0)

### Phase 4 — ADR-0049

`docs/kb/adr/0049-audit-augmentation.md` (Status: Accepted):
- Context: 5 repo audit 의 valkey 잔여 5 ❌ — postgres / mongodb 와 정합
- Decision: postgres 패턴 cp (lefthook + scripts + UPGRADING)
- Consequences: valkey audit P1-11/12/13 + OP-2 + OP-10 ✅ 적용

### Phase 5 — 검증

```bash
bash commons/scripts/audit-production-grade.sh /Users/phil/Workspace/keiailab
# valkey 의 5 항목 ✅ 확인
```

## 5. 리스크

| 리스크 | 영향 | 완화 |
|---|---|---|
| ralph-loop iteration 12 와 동시 작업 | main 변경 충돌 | branch 분리 + sequential push (ralph-loop 의 README PR 머지 후 본 spec 진입) |
| lefthook 3종 도구 (kube-linter / go-licenses / markdown-link-check) 미설치 | hook 실행 안 됨 | graceful skip 정책 (postgres 패턴) |
| valkey 의 chart 차이 (ValkeyCluster vs PostgresCluster) | helm-publish 의 chart name 불일치 | Phase 2 의 chart name 명시 조정 |

## 6. 성공 기준

```bash
# 1. lefthook 3종 hook 존재
grep -qE "kube-linter|go-licenses|markdown-link-check" valkey-operator/.lefthook.yml

# 2. helm-publish.sh 존재 + syntax
test -x valkey-operator/scripts/helm-publish.sh && bash -n valkey-operator/scripts/helm-publish.sh

# 3. UPGRADING.md 존재
test -f valkey-operator/docs/UPGRADING.md

# 4. ADR 0049 Accepted
grep -q "Status.*Accepted" valkey-operator/docs/kb/adr/0049-audit-augmentation.md

# 5. audit 5 항목 ✅
bash commons/scripts/audit-production-grade.sh /Users/phil/Workspace/keiailab \
  | grep -E "^\| (P1-11|P1-12|P1-13|OP-2|OP-10)" \
  | awk -F'|' '{print $5}' | grep -c "✅"
# 결과: 5
```

## 7. 본 spec 의 진입 정책

본 spec 의 실 실행은 다음 중 하나 시점:

1. **ralph-loop 가 본 spec 인지 후 자율 적용** (가장 자연스러움)
2. **사용자 명시 결정 변경** — "본 thread 가 valkey 직접 작업 가능" 결정 시 본 thread 가 진행
3. **별 subagent dispatch** — 사용자 결정 #2 의 "dispatch" 영역 — subagent 가 valkey 만지기

본 spec 은 진입 정책 *명시 안 함* — 위 3 옵션 중 사용자 명시 또는 ralph-loop 자율 결정.

## 8. Out-of-scope

- valkey-operator 의 *코드 변경* (Reconcile, CR spec, ...) — 본 spec 은 *audit 보강* 만
- README 4-lang 본문 확장 (ralph-loop PR #165/#166 이 진행)
- Sprint 1 채택 (이미 머지)
- ADR-0048 갱신 (ralph-loop 가 Accepted 승격 완료)
