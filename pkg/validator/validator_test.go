package validator

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{"valid lowercase", "my-namespace", true},
		{"valid with numbers", "app-123", true},
		{"uppercase", "My-Namespace", false},
		{"starts with hyphen", "-invalid", false},
		{"ends with hyphen", "invalid-", false},
		{"underscore", "my_namespace", false},
		{"empty", "", false},
		{"too long", "this-namespace-name-is-way-too-long-and-exceeds-the-maximum-sixty-three-character-limit-for-sure", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := Validate(tt.input)
			if valid != tt.wantValid {
				t.Errorf("Validate(%q) = %v, want %v", tt.input, valid, tt.wantValid)
			}
		})
	}
}

func TestIsReserved(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		reserved bool
	}{
		{"default", "default", true},
		{"kube-system", "kube-system", true},
		{"normal name", "my-app", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsReserved(tt.input); got != tt.reserved {
				t.Errorf("IsReserved(%q) = %v, want %v", tt.input, got, tt.reserved)
			}
		})
	}
}
