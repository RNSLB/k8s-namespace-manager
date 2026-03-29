package quota

import (
	"testing"
)

func TestDefaultQuotaConfig(t *testing.T) {
	config := DefaultQuotaConfig("test-namespace")

	if config.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", config.Namespace)
	}

	if config.Name != "default-quota" {
		t.Errorf("Expected name 'default-quota', got '%s'", config.Name)
	}

	if config.RequestsCPU != "10" {
		t.Errorf("Expected RequestsCPU '10', got '%s'", config.RequestsCPU)
	}
}

func TestValidateQuotaConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    QuotaConfig
		wantError bool
	}{
		{
			name: "valid config",
			config: QuotaConfig{
				Name:         "test-quota",
				Namespace:    "test-ns",
				RequestsCPU:  "10",
			},
			wantError: false,
		},
		{
			name: "empty name",
			config: QuotaConfig{
				Name:         "",
				Namespace:    "test-ns",
				RequestsCPU:  "10",
			},
			wantError: true,
		},
		{
			name: "empty namespace",
			config: QuotaConfig{
				Name:         "test-quota",
				Namespace:    "",
				RequestsCPU:  "10",
			},
			wantError: true,
		},
		{
			name: "no resources",
			config: QuotaConfig{
				Name:      "test-quota",
				Namespace: "test-ns",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuotaConfig(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateQuotaConfig() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestBuildResourceQuota(t *testing.T) {
	config := QuotaConfig{
		Name:           "test-quota",
		Namespace:      "test-ns",
		RequestsCPU:    "10",
		RequestsMemory: "20Gi",
		LimitsCPU:      "20",
		LimitsMemory:   "40Gi",
		MaxPods:        "50",
	}

	quota, err := BuildResourceQuota(config)
	if err != nil {
		t.Fatalf("BuildResourceQuota() failed: %v", err)
	}

	if quota.Name != "test-quota" {
		t.Errorf("Expected quota name 'test-quota', got '%s'", quota.Name)
	}

	if quota.Namespace != "test-ns" {
		t.Errorf("Expected namespace 'test-ns', got '%s'", quota.Namespace)
	}

	if len(quota.Spec.Hard) != 5 {
		t.Errorf("Expected 5 hard limits, got %d", len(quota.Spec.Hard))
	}
}
