# 한국어 용어 사전 (Glossary)

> **English** | [한국어](glossary-ko.md) | [日本語](glossary-ja.md) (예정) | [中文](glossary-zh.md) (예정)
>
> 본 사전은 keiailab operator family 4 repo (operator-commons + postgres-operator + mongodb-operator + valkey-operator) 의 한국어 번역 시 *반드시 참조* 하는 표준 용어 사전입니다. 사용자 supercycle plan `~/.claude/plans/2026-05-21-keiailab-operator-supercycle.md` §6 Wave 4 §4.4 정합.

## §1 일관성 규칙

1. **코드 식별자는 영문 그대로** — 예: `ValkeyCluster`, `kubectl`, `Helm`, `pkg/probes`. 한글 번역 금지.
2. **표준 K8s 용어는 영문 우선 + 한국어 부기** — 예: `Pod (파드)`, `Deployment (디플로이먼트)`. 본문 첫 등장 시 영문 + 괄호 한국어, 이후 한국어 단독 허용.
3. **operator-commons API 명칭은 영문 그대로** — 예: `Reconciler`, `Finalizer`, `EventRecorder`. 한국어 부기 가능 (`Finalizer (파이널라이저)`).
4. **외부 사용자 가시 문서 (README/CONTRIBUTING/SECURITY 등)** = 격식체 (`-습니다`, `-입니다`). 내부 문서 (HANDOFF/AGENTS 등) = 평어 또는 자유.
5. **존댓말 평어 혼용 금지** — 한 문서 안 일관.

## §2 Kubernetes 표준 용어

| English (canonical) | 한국어 |
|---|---|
| CustomResourceDefinition (CRD) | 커스텀 리소스 정의 (CRD) |
| Custom Resource (CR) | 커스텀 리소스 (CR) |
| Reconciler | 리컨실러 (또는 영문 그대로) |
| Reconcile Loop | 조정 루프 (또는 reconcile 루프) |
| Controller | 컨트롤러 |
| Operator | 오퍼레이터 (또는 영문 그대로) |
| Operator Pattern | 오퍼레이터 패턴 |
| Finalizer | 파이널라이저 (또는 영문 그대로) |
| Webhook | 웹훅 (또는 영문 그대로) |
| Validating Admission Webhook | 검증 어드미션 웹훅 |
| Mutating Admission Webhook | 변경 어드미션 웹훅 |
| Conversion Webhook | 변환 웹훅 |
| Pod | 파드 (Pod) |
| Deployment | 디플로이먼트 (Deployment) |
| StatefulSet | 스테이트풀셋 (StatefulSet) |
| DaemonSet | 데몬셋 (DaemonSet) |
| Job | 잡 (Job) |
| CronJob | 크론잡 (CronJob) |
| Service | 서비스 (Service) |
| Ingress | 인그레스 (Ingress) |
| ConfigMap | 컨피그맵 (ConfigMap) |
| Secret | 시크릿 (Secret) |
| PersistentVolume (PV) | 영구 볼륨 (PV) |
| PersistentVolumeClaim (PVC) | 영구 볼륨 클레임 (PVC) |
| StorageClass | 스토리지 클래스 (StorageClass) |
| Namespace | 네임스페이스 |
| Node | 노드 |
| Cluster | 클러스터 |
| Control Plane | 컨트롤 플레인 |
| Worker Node | 워커 노드 |
| API Server | API 서버 |
| Scheduler | 스케줄러 |
| kubelet | kubelet (영문 그대로) |
| kube-proxy | kube-proxy (영문 그대로) |
| etcd | etcd (영문 그대로) |
| Helm Chart | 헬름 차트 (Helm Chart) |
| Helm Release | 헬름 릴리스 |
| Kubernetes RBAC | 쿠버네티스 RBAC |
| ServiceAccount | 서비스 어카운트 (ServiceAccount) |
| Role / ClusterRole | 롤 / 클러스터 롤 |
| RoleBinding / ClusterRoleBinding | 롤 바인딩 / 클러스터 롤 바인딩 |
| NetworkPolicy | 네트워크 정책 (NetworkPolicy) |
| PodSecurityAdmission (PSA) | 파드 보안 어드미션 (PSA) |
| restricted profile | restricted 프로파일 (영문 그대로) |
| baseline profile | baseline 프로파일 (영문 그대로) |
| privileged profile | privileged 프로파일 (영문 그대로) |
| Probe (liveness/readiness/startup) | 프로브 (liveness/readiness/startup) |
| HTTPGetAction | HTTP Get 액션 (HTTPGetAction) |
| TCPSocketAction | TCP 소켓 액션 (TCPSocketAction) |
| ExecAction | 실행 액션 (ExecAction) |
| InitContainer | 초기화 컨테이너 (InitContainer) |
| Sidecar Container | 사이드카 컨테이너 |

## §3 operator-commons 라이브러리 용어

| English (canonical) | 한국어 |
|---|---|
| `pkg/finalizer` | 파이널라이저 패키지 (`pkg/finalizer`) |
| `pkg/labels` | 라벨 패키지 (`pkg/labels`) |
| `pkg/status` | 상태 조건 패키지 (`pkg/status`) |
| `pkg/version` | 버전 호환성 패키지 (`pkg/version`) |
| `pkg/monitoring` | 모니터링 패키지 (`pkg/monitoring`) |
| `pkg/networkpolicy` | 네트워크 정책 패키지 (`pkg/networkpolicy`) |
| `pkg/security` | 보안 패키지 (`pkg/security`) |
| `pkg/webhook` | 웹훅 패키지 (`pkg/webhook`) |
| `pkg/probes` (v0.8.0 신규) | 프로브 빌더 패키지 (`pkg/probes`) |
| `pkg/storageclass` (v0.8.0 신규) | 스토리지 클래스 패키지 (`pkg/storageclass`) |
| `pkg/events` (v0.8.0 신규) | 이벤트 레코더 패키지 (`pkg/events`) |
| Recorder interface | Recorder 인터페이스 (영문 그대로) |
| EventType (Normal / Warning) | 이벤트 타입 (Normal / Warning) |
| Reason (event reason) | 이벤트 사유 (Reason) |
| Builder pattern | 빌더 패턴 |
| Fluent API | 플루언트 API (Fluent API) |
| API Stability Tier | API 안정성 등급 |
| Stable / Beta / Experimental | 안정 / 베타 / 실험 (영문 권장 — Stable/Beta/Experimental) |
| Tier 격상 (promotion) | 등급 격상 (Tier promotion) |
| Breaking change | 호환성 깨지는 변경 (Breaking change) |
| Semver (Semantic Versioning) | 시맨틱 버저닝 (Semver) |
| Deprecated | 사용 중단 예정 (Deprecated) |

## §4 reconciler 패턴 용어

| English (canonical) | 한국어 |
|---|---|
| Desired state | 의도한 상태 |
| Actual state | 실제 상태 |
| Reconcile | 조정 (reconcile) |
| Drift | 드리프트 (영문 그대로) |
| Idempotency | 멱등성 |
| Eventual consistency | 결과적 일관성 |
| Owner reference | 소유자 참조 (OwnerReference) |
| Garbage collection | 가비지 컬렉션 (GC) |
| Status condition | 상태 조건 (status condition) |
| Available / Ready / Progressing / Degraded | Available / Ready / Progressing / Degraded (영문 그대로) |
| Failover | 페일오버 (failover) |
| Provisioning | 프로비저닝 |
| Rolling update | 롤링 업데이트 |
| Blue-green deployment | 블루-그린 배포 |
| Canary deployment | 카나리 배포 |
| Backup / Restore | 백업 / 복원 |
| Point-in-time recovery (PITR) | 시점 복구 (PITR) |
| Horizontal scaling | 수평 확장 |
| Vertical scaling | 수직 확장 |
| Sharding | 샤딩 |
| Replication | 복제 (replication) |
| Primary / Replica / Secondary | Primary / Replica / Secondary (영문 그대로) |

## §5 보안 + 인증 용어

| English (canonical) | 한국어 |
|---|---|
| Authentication (AuthN) | 인증 (Authentication) |
| Authorization (AuthZ) | 권한 부여 (Authorization) |
| OIDC | OIDC (영문 그대로) |
| LDAP | LDAP (영문 그대로) |
| TLS / mTLS | TLS / mTLS (영문 그대로) |
| Certificate | 인증서 (certificate) |
| Service mesh | 서비스 메시 |
| Tenant / Tenancy | 테넌트 / 테넌시 |
| Multi-tenancy | 멀티 테넌시 |
| Encryption at rest | 저장 시 암호화 |
| Encryption in transit | 전송 시 암호화 |

## §6 운영 + 관찰 용어

| English (canonical) | 한국어 |
|---|---|
| Observability | 관찰 가능성 (observability) |
| Metric / Monitoring | 메트릭 / 모니터링 |
| Logging | 로깅 |
| Tracing | 트레이싱 |
| Prometheus | Prometheus (영문 그대로) |
| Grafana | Grafana (영문 그대로) |
| Alert / Alerting | 알림 / 알림 발생 |
| Service Level Objective (SLO) | 서비스 수준 목표 (SLO) |
| Service Level Indicator (SLI) | 서비스 수준 지표 (SLI) |
| Postmortem | 사후 분석 (postmortem) |
| Incident | 인시던트 (incident) |
| Severity (SEV-1 / SEV-2 / SEV-3) | 심각도 (SEV-1 / SEV-2 / SEV-3, 영문 그대로) |

## §7 거버넌스 + 협업 용어

| English (canonical) | 한국어 |
|---|---|
| ADR (Architecture Decision Record) | 아키텍처 결정 기록 (ADR) |
| RFC (Request For Comments) | RFC (영문 그대로) |
| Roadmap | 로드맵 |
| Adopter | 호출자 / 도입자 (Adopter) |
| Maintainer | 메인테이너 |
| Contributor | 기여자 |
| Pull Request (PR) | 풀 리퀘스트 (PR) |
| Merge Request (MR) | 머지 리퀘스트 (MR, GitLab) |
| Issue | 이슈 |
| Code review | 코드 리뷰 |
| Squash merge | 스쿼시 머지 |
| Rebase | 리베이스 |
| Cherry-pick | 체리픽 (cherry-pick) |
| CI/CD | CI/CD (영문 그대로) |
| Pipeline | 파이프라인 |
| Lint / Linter | 린트 / 린터 |
| Coverage (test coverage) | 커버리지 (테스트 커버리지) |

## §8 keiailab 운영 컨텍스트 (사내 용어)

| English (canonical) | 한국어 |
|---|---|
| keiailab operator family | keiailab 오퍼레이터 패밀리 |
| operator-commons | operator-commons (영문 그대로 — Go module 명) |
| supercycle | 슈퍼사이클 (supercycle) — 사용자 명명 진본 |
| Wave (Wave 0 ~ Wave 5) | Wave (영문 그대로) |
| Phase | 페이즈 (Phase) |
| Cadence checkpoint | 케이던스 체크포인트 |

## §9 참조

- 사용자 supercycle plan: `~/.claude/plans/2026-05-21-keiailab-operator-supercycle.md` §6 Wave 4 §4.4 "각 언어 별 *glossary* 1 파일"
- deep-petting-cookie plan: `~/.claude/plans/deep-petting-cookie.md` §Phase 0.8 (i18n-glossary 신규)
- Kubernetes 한국어 번역 가이드: <https://kubernetes.io/ko/docs/contribute/localization_ko/>
- 본 라이브러리 ROADMAP §API Stability Tier: `../../ROADMAP.md`

## §10 변경 이력

| Date | Change | Refs |
|---|---|---|
| 2026-05-21 | 신설 — 한국어 표준 용어 사전 v0.1 (10 categories, ~120 terms) | deep-petting-cookie §Phase 0.8 + supercycle Wave 4 §4.4 sister |
