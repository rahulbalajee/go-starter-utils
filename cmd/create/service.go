package create

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	validName = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	outputDir string
)

type serviceData struct {
	Name       string
	PascalName string
	ModulePath string
}

type scaffoldFile struct {
	tmplPath   string
	outputPath string
}

var serviceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Create a service with hexagonal architecture",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service, err := validateServiceName(args[0])
		if err != nil {
			return err
		}

		basePath := filepath.Join(outputDir, service+"-service")

		if _, err := os.Stat(basePath); err == nil {
			return fmt.Errorf("directory already exists at path %s", basePath)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check path %s: %w", basePath, err)
		}

		data := serviceData{
			Name:       service,
			PascalName: toPascalCase(service),
			ModulePath: service + "-service",
		}

		if err := scaffoldService(basePath, data); err != nil {
			_ = os.RemoveAll(basePath)
			return fmt.Errorf("scaffold failed (rolled back): %w", err)
		}

		printSuccess(service, basePath)
		return nil
	},
}

var serviceDirs = []string{
	"cmd",
	"internal/domain",
	"internal/service",
	"internal/infrastructure/events",
	"internal/infrastructure/grpc",
	"internal/infrastructure/http",
	"internal/infrastructure/repository",
	"pkg/types",
}

func scaffoldService(basePath string, data serviceData) error {
	for _, dir := range serviceDirs {
		if err := os.MkdirAll(filepath.Join(basePath, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	files := []scaffoldFile{
		{"templates/service/cmd_main.go.tmpl", "cmd/main.go"},
		{"templates/service/domain.go.tmpl", fmt.Sprintf("internal/domain/%s.go", data.Name)},
		{"templates/service/service_impl.go.tmpl", "internal/service/service.go"},
		{"templates/service/types.go.tmpl", "pkg/types/types.go"},
		{"templates/service/go_mod.tmpl", "go.mod"},
		{"templates/service/readme.md.tmpl", "README.md"},
	}

	for _, f := range files {
		if err := renderTemplate(basePath, f, data); err != nil {
			return err
		}
	}

	gitkeepDirs := []string{
		"internal/infrastructure/events",
		"internal/infrastructure/grpc",
		"internal/infrastructure/http",
		"internal/infrastructure/repository",
	}
	for _, dir := range gitkeepDirs {
		path := filepath.Join(basePath, dir, ".gitkeep")
		if err := os.WriteFile(path, nil, 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep in %s: %w", dir, err)
		}
	}

	return nil
}

func renderTemplate(basePath string, f scaffoldFile, data serviceData) error {
	content, err := serviceTemplates.ReadFile(f.tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", f.tmplPath, err)
	}

	tmpl, err := template.New(filepath.Base(f.tmplPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", f.tmplPath, err)
	}

	outPath := filepath.Join(basePath, f.outputPath)
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", f.outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to render %s: %w", f.tmplPath, err)
	}

	return nil
}

func printSuccess(service, basePath string) {
	fmt.Printf("Successfully created %s service at %s\n\n", service, basePath)
	fmt.Println("Generated files:")
	fmt.Printf("  %s/cmd/main.go\n", basePath)
	fmt.Printf("  %s/internal/domain/%s.go\n", basePath, service)
	fmt.Printf("  %s/internal/service/service.go\n", basePath)
	fmt.Printf("  %s/pkg/types/types.go\n", basePath)
	fmt.Printf("  %s/go.mod\n", basePath)
	fmt.Printf("  %s/README.md\n", basePath)
	fmt.Println("\nGet started:")
	fmt.Printf("  cd %s\n", basePath)
	fmt.Println("  go mod tidy")
}

func validateServiceName(name string) (string, error) {
	service := strings.TrimSpace(name)

	if service == "" {
		return "", fmt.Errorf("service name cannot be empty")
	}

	if !validName.MatchString(service) {
		return "", fmt.Errorf("service name must contain only lowercase letters, numbers, and dashes")
	}

	return service, nil
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func init() {
	serviceCmd.Flags().StringVarP(&outputDir, "output", "o", "services", "parent directory for the generated service")
	Cmd.AddCommand(serviceCmd)
}
