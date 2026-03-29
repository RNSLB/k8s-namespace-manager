package quota

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// QuotaConfig holds resource quota configuration
type QuotaConfig struct {
	Name             string
	Namespace        string
	RequestsCPU      string // e.g., "10"
	RequestsMemory   string // e.g., "20Gi"
	LimitsCPU        string // e.g., "20"
	LimitsMemory     string // e.g., "40Gi"
	MaxPods          string // e.g., "50"
}

// DefaultQuotaConfig returns sensible defaults for resource quotas
func DefaultQuotaConfig(namespace string) QuotaConfig {
	return QuotaConfig{
		Name:           "default-quota",
		Namespace:      namespace,
		RequestsCPU:    "10",
		RequestsMemory: "20Gi",
		LimitsCPU:      "20",
		LimitsMemory:   "40Gi",
		MaxPods:        "50",
	}
}

// BuildResourceQuota creates a ResourceQuota object from config
func BuildResourceQuota(config QuotaConfig) (*corev1.ResourceQuota, error) {
	// Create resource list for hard limits
	hard := corev1.ResourceList{}

	// Add CPU requests
	if config.RequestsCPU != "" {
		cpuQuantity, err := resource.ParseQuantity(config.RequestsCPU)
		if err != nil {
			return nil, fmt.Errorf("invalid requests.cpu '%s': %w", config.RequestsCPU, err)
		}
		hard[corev1.ResourceRequestsCPU] = cpuQuantity
	}

	// Add Memory requests
	if config.RequestsMemory != "" {
		memQuantity, err := resource.ParseQuantity(config.RequestsMemory)
		if err != nil {
			return nil, fmt.Errorf("invalid requests.memory '%s': %w", config.RequestsMemory, err)
		}
		hard[corev1.ResourceRequestsMemory] = memQuantity
	}

	// Add CPU limits
	if config.LimitsCPU != "" {
		cpuQuantity, err := resource.ParseQuantity(config.LimitsCPU)
		if err != nil {
			return nil, fmt.Errorf("invalid limits.cpu '%s': %w", config.LimitsCPU, err)
		}
		hard[corev1.ResourceLimitsCPU] = cpuQuantity
	}

	// Add Memory limits
	if config.LimitsMemory != "" {
		memQuantity, err := resource.ParseQuantity(config.LimitsMemory)
		if err != nil {
			return nil, fmt.Errorf("invalid limits.memory '%s': %w", config.LimitsMemory, err)
		}
		hard[corev1.ResourceLimitsMemory] = memQuantity
	}

	// Add pod limit
	if config.MaxPods != "" {
		podsQuantity, err := resource.ParseQuantity(config.MaxPods)
		if err != nil {
			return nil, fmt.Errorf("invalid pods '%s': %w", config.MaxPods, err)
		}
		hard[corev1.ResourcePods] = podsQuantity
	}

	// Build ResourceQuota object
	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: hard,
		},
	}

	return quota, nil
}

// ValidateQuotaConfig checks if quota configuration is valid
func ValidateQuotaConfig(config QuotaConfig) error {
	if config.Name == "" {
		return fmt.Errorf("quota name cannot be empty")
	}

	if config.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	// Validate that at least one resource is specified
	if config.RequestsCPU == "" && config.RequestsMemory == "" &&
		config.LimitsCPU == "" && config.LimitsMemory == "" &&
		config.MaxPods == "" {
		return fmt.Errorf("at least one resource limit must be specified")
	}

	return nil
}

// FormatQuota returns a human-readable representation of quota
func FormatQuota(quota *corev1.ResourceQuota) string {
	output := fmt.Sprintf("ResourceQuota: %s (Namespace: %s)\n", quota.Name, quota.Namespace)
	output += "Hard Limits:\n"

	for resourceName, quantity := range quota.Spec.Hard {
		output += fmt.Sprintf("  %s: %s\n", resourceName, quantity.String())
	}

	return output
}
