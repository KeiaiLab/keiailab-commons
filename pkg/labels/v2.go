package labels

// V2 추가 라벨 — Kubernetes 1.30+ Recommended labels v2 mapping.
//
// K8s 1.30 부터 도입된 추가 라벨 (`app.kubernetes.io/created-by`,
// `app.kubernetes.io/managed-by` 외) 를 v1 Set 위에 일관 적용.
//
// Refs: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
//
//	docs/ROADMAP.md "Recommended labels v2 매핑 (K8s 1.30+)"
const (
	// LabelCreatedBy — 리소스를 생성한 controller / tool. RFC: K8s 1.30+
	LabelCreatedBy = "app.kubernetes.io/created-by"
	// LabelTier — 컴포넌트 tier (e.g. "frontend" / "backend" / "cache").
	LabelTier = "app.kubernetes.io/tier"
	// LabelOwner — 책임 팀 / 사용자.
	LabelOwner = "app.kubernetes.io/owner"
)

// AllV2 는 v1 All() 위에 v2 라벨 (createdBy / tier / owner) 을 덧붙인다.
// v2 라벨이 빈 문자열인 경우 생략 (선택적).
//
// 사용 예:
//
//	labels := labels.AllV2(set, labels.V2{
//	    CreatedBy: "downstream-operator",
//	    Tier:      "cache",
//	    Owner:     "platform-team",
//	})
func (s Set) AllV2(v2 V2) map[string]string {
	out := s.All()
	if v2.CreatedBy != "" {
		out[LabelCreatedBy] = v2.CreatedBy
	}
	if v2.Tier != "" {
		out[LabelTier] = v2.Tier
	}
	if v2.Owner != "" {
		out[LabelOwner] = v2.Owner
	}
	return out
}

// V2 는 K8s 1.30+ v2 라벨 매핑.
type V2 struct {
	// CreatedBy — controller / tool 이름 (e.g. "downstream-operator").
	CreatedBy string
	// Tier — 컴포넌트 tier (e.g. "cache").
	Tier string
	// Owner — 책임 팀.
	Owner string
}
