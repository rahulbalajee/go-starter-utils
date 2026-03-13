package create

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid simple name", "order", "order", false},
		{"valid with dashes", "order-status", "order-status", false},
		{"valid with numbers", "service1", "service1", false},
		{"valid mixed", "my-svc-2", "my-svc-2", false},
		{"empty string", "", "", true},
		{"whitespace only", "   ", "", true},
		{"uppercase letters", "Order", "", true},
		{"underscores", "order_service", "", true},
		{"spaces in name", "order service", "", true},
		{"special characters", "order@svc", "", true},
		{"leading whitespace trimmed", "  order  ", "order", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateServiceName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateServiceName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateServiceName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"order", "Order"},
		{"order-status", "OrderStatus"},
		{"a", "A"},
		{"my-long-name", "MyLongName"},
		{"api", "Api"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := toPascalCase(tt.input); got != tt.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestScaffoldService(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "order-service")

	data := serviceData{
		Name:       "order",
		PascalName: "Order",
		ModulePath: "order-service",
	}

	if err := scaffoldService(basePath, data); err != nil {
		t.Fatalf("scaffoldService() error = %v", err)
	}

	expectedDirs := []string{
		"cmd",
		"internal/domain",
		"internal/service",
		"internal/infrastructure/events",
		"internal/infrastructure/grpc",
		"internal/infrastructure/http",
		"internal/infrastructure/repository",
		"pkg/types",
	}

	for _, dir := range expectedDirs {
		path := filepath.Join(basePath, dir)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected directory %s to exist: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", dir)
		}
	}

	expectedFiles := []string{
		"cmd/main.go",
		"internal/domain/order.go",
		"internal/service/service.go",
		"pkg/types/types.go",
		"go.mod",
		"README.md",
		"internal/infrastructure/events/.gitkeep",
		"internal/infrastructure/grpc/.gitkeep",
		"internal/infrastructure/http/.gitkeep",
		"internal/infrastructure/repository/.gitkeep",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(basePath, file)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected file %s to exist: %v", file, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("expected %s to be a file, got directory", file)
		}
	}
}

func TestScaffoldServiceGeneratedContent(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "order-service")

	data := serviceData{
		Name:       "order",
		PascalName: "Order",
		ModulePath: "order-service",
	}

	if err := scaffoldService(basePath, data); err != nil {
		t.Fatalf("scaffoldService() error = %v", err)
	}

	domainContent, err := os.ReadFile(filepath.Join(basePath, "internal/domain/order.go"))
	if err != nil {
		t.Fatalf("failed to read domain file: %v", err)
	}

	checks := []string{
		"package domain",
		"type Order struct",
		"type OrderRepository interface",
		"type OrderService interface",
	}
	for _, check := range checks {
		if !contains(string(domainContent), check) {
			t.Errorf("domain file missing expected content: %q", check)
		}
	}

	svcContent, err := os.ReadFile(filepath.Join(basePath, "internal/service/service.go"))
	if err != nil {
		t.Fatalf("failed to read service file: %v", err)
	}

	svcChecks := []string{
		"package service",
		"type OrderService struct",
		"func NewOrderService",
		`"order-service/internal/domain"`,
	}
	for _, check := range svcChecks {
		if !contains(string(svcContent), check) {
			t.Errorf("service file missing expected content: %q", check)
		}
	}

	modContent, err := os.ReadFile(filepath.Join(basePath, "go.mod"))
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}
	if !contains(string(modContent), "module order-service") {
		t.Error("go.mod missing expected module path")
	}
}

func TestScaffoldServiceKebabCase(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "order-status-service")

	data := serviceData{
		Name:       "order-status",
		PascalName: "OrderStatus",
		ModulePath: "order-status-service",
	}

	if err := scaffoldService(basePath, data); err != nil {
		t.Fatalf("scaffoldService() error = %v", err)
	}

	domainContent, err := os.ReadFile(filepath.Join(basePath, "internal/domain/order-status.go"))
	if err != nil {
		t.Fatalf("failed to read domain file: %v", err)
	}

	if !contains(string(domainContent), "type OrderStatus struct") {
		t.Error("kebab-case name not properly converted to PascalCase in domain")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
