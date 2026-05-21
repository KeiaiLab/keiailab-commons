package pvc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DefaultVCTName 는 StatefulSet 의 volumeClaimTemplate 표준 이름이다. K8s STS
// controller 는 각 replica 별로 `<vct-name>-<sts-name>-<ordinal>` 형태로 PVC 를
// 명명하므로, 본 상수 + STS 이름 + ordinal 로 PVC 이름이 재구성된다.
//
// downstream operator (mongodb / postgres / valkey) 모두 STS VCT 이름을 "data" 로 통일.
const DefaultVCTName = "data"

// ErrEmptyVCTName 은 옵션으로 빈 VCT 이름이 주어졌을 때 반환한다.
var ErrEmptyVCTName = errors.New("pvc: VCT name must not be empty")

// Option 은 ExpandDataPVCs 동작을 변경하는 함수형 옵션.
type Option func(*config)

type config struct {
	vctName string
}

func newConfig(opts ...Option) (*config, error) {
	c := &config{vctName: DefaultVCTName}
	for _, o := range opts {
		o(c)
	}
	if c.vctName == "" {
		return nil, ErrEmptyVCTName
	}
	return c, nil
}

// WithVCTName 은 STS volumeClaimTemplate 이름을 변경한다. 기본값은
// DefaultVCTName ("data"). downstream operator 모두 "data" 이므로 별도 지정 불요.
func WithVCTName(name string) Option {
	return func(c *config) { c.vctName = name }
}

// PVCNamePrefix — STS controller 가 명명하는 PVC 이름 prefix 를 반환한다.
// K8s STS controller 는 각 replica 별로 `<vct-name>-<sts-name>-<ordinal>` 형태로
// PVC 를 생성하므로, 본 함수는 그 prefix `<vct-name>-<sts-name>-` 를 반환한다.
//
// vctName 이 비어 있으면 DefaultVCTName 으로 대체한다.
func PVCNamePrefix(vctName, stsName string) string {
	if vctName == "" {
		vctName = DefaultVCTName
	}
	return vctName + "-" + stsName + "-"
}

// ExpandDataPVCs 는 desired 가 현재 PVC 크기보다 클 때, stsNames 에 속하는
// 모든 STS data PVC 의 spec.resources.requests.storage 를 patch 한다.
//
// 동작:
//   - desired.IsZero() 또는 len(stsNames) == 0 → noop.
//   - namespace 내 모든 PVC 를 List 후 prefix `<vct-name>-<stsName>-` 매칭.
//   - 각 PVC 별로 expandSinglePVC 실행 — best-effort: 1 PVC 실패가 차단 안 됨.
//
// 사전 조건 (위반 시 noop + 로그):
//   - StorageClass.AllowVolumeExpansion == true (false → skip).
//   - PVC 가 Bound 상태 (Pending / Lost → skip).
//   - desired > current (감소는 webhook 에서 reject — 도달 시 skip).
//
// CSI driver 가 online resize 미지원 시 PVC.status.conditions 가
// FileSystemResizePending 으로 남고 다음 pod restart 시 완료. 본 함수는
// *patch 만* 하고 폴링하지 않는다.
//
// 진단 로그는 controller-runtime/pkg/log 의 LoggerFromContext 를 사용한다.
// 호출자가 ctxlog.IntoContext 로 logger 를 주입하면 자동 구조화 로깅.
func ExpandDataPVCs(
	ctx context.Context,
	c client.Client,
	namespace string,
	stsNames []string,
	desired resource.Quantity,
	opts ...Option,
) error {
	if desired.IsZero() || len(stsNames) == 0 {
		return nil
	}
	cfg, err := newConfig(opts...)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx).WithName("pvc-expansion").WithValues("namespace", namespace)

	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := c.List(ctx, pvcList, client.InNamespace(namespace)); err != nil {
		return fmt.Errorf("pvc: list PVCs in %s: %w", namespace, err)
	}

	prefixes := make([]string, 0, len(stsNames))
	for _, n := range stsNames {
		prefixes = append(prefixes, PVCNamePrefix(cfg.vctName, n))
	}

	for i := range pvcList.Items {
		p := &pvcList.Items[i]
		if !matchAnyPrefix(p.Name, prefixes) {
			continue
		}
		if err := expandSinglePVC(ctx, c, p, desired); err != nil {
			logger.Error(err, "PVC expansion failed", "pvc", p.Name)
		}
	}
	return nil
}

func matchAnyPrefix(name string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func expandSinglePVC(
	ctx context.Context,
	c client.Client,
	p *corev1.PersistentVolumeClaim,
	desired resource.Quantity,
) error {
	logger := log.FromContext(ctx).WithName("pvc-expansion")

	if p.Status.Phase != corev1.ClaimBound {
		logger.V(1).Info("skip non-Bound PVC", "pvc", p.Name, "phase", p.Status.Phase)
		return nil
	}
	current, ok := p.Spec.Resources.Requests[corev1.ResourceStorage]
	if !ok {
		return fmt.Errorf("pvc: %s missing spec.resources.requests.storage", p.Name)
	}
	if desired.Cmp(current) <= 0 {
		return nil // 이미 desired 이상.
	}

	// StorageClass.AllowVolumeExpansion 검증.
	if p.Spec.StorageClassName != nil && *p.Spec.StorageClassName != "" {
		sc := &storagev1.StorageClass{}
		if err := c.Get(ctx, types.NamespacedName{Name: *p.Spec.StorageClassName}, sc); err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info("skip: StorageClass not found",
					"pvc", p.Name, "storageClass", *p.Spec.StorageClassName)
				return nil
			}
			return fmt.Errorf("pvc: get StorageClass %s: %w", *p.Spec.StorageClassName, err)
		}
		if sc.AllowVolumeExpansion == nil || !*sc.AllowVolumeExpansion {
			logger.Info("skip: StorageClass does not allow expansion",
				"pvc", p.Name, "storageClass", sc.Name)
			return nil
		}
	}

	patched := p.DeepCopy()
	if patched.Spec.Resources.Requests == nil {
		patched.Spec.Resources.Requests = corev1.ResourceList{}
	}
	patched.Spec.Resources.Requests[corev1.ResourceStorage] = desired
	if err := c.Patch(ctx, patched, client.MergeFrom(p)); err != nil {
		return fmt.Errorf("pvc: patch %s: %w", p.Name, err)
	}
	logger.Info("PVC expansion patched",
		"pvc", p.Name, "from", current.String(), "to", desired.String())
	return nil
}
