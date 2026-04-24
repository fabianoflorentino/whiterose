package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type VersionChecker struct{}

func NewVersionChecker() *VersionChecker {
	return &VersionChecker{}
}

func (vc *VersionChecker) ListGoVersions() error {
	fmt.Println("Fetching Go versions...")

	resp, err := http.Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return fmt.Errorf("failed to fetch Go versions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var versions []struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}

	if err := json.Unmarshal(body, &versions); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Println("\nAvailable Go versions:")
	var stable []string
	var unstable []string

	for _, v := range versions {
		if v.Stable {
			stable = append(stable, v.Version)
		} else if strings.Contains(v.Version, "beta") || strings.Contains(v.Version, "rc") {
			unstable = append(unstable, v.Version)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(stable)))
	for _, v := range stable[:min(10, len(stable))] {
		fmt.Printf("  %s (stable)\n", v)
	}

	if len(unstable) > 0 {
		fmt.Println("\nUnstable versions:")
		sort.Sort(sort.Reverse(sort.StringSlice(unstable)))
		for _, v := range unstable[:min(5, len(unstable))] {
			fmt.Printf("  %s\n", v)
		}
	}

	return nil
}

func (vc *VersionChecker) GetCurrentGoVersion(goModPath string) string {
	content, _ := os.ReadFile(filepath.Join(goModPath, "go.mod"))
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "go ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}
	return ""
}

func (vc *VersionChecker) ListGoLibUpdates(goModPath string) error {
	fmt.Println("Checking for library updates...")

	cmd := exec.Command("go", "list", "-m", "-u", "all")
	cmd.Dir = goModPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list updates: %w", err)
	}

	fmt.Println("\nAvailable updates:")
	lines := strings.Split(string(out), "\n")
	hasUpdates := false

	for _, line := range lines {
		if strings.HasSuffix(line, "(latest)") || strings.Contains(line, "->") {
			if !strings.HasPrefix(line, "go:") {
				fmt.Printf("  %s\n", line)
				hasUpdates = true
			}
		}
	}

	if !hasUpdates {
		fmt.Println("  All dependencies are up to date")
	}

	return nil
}

func (vc *VersionChecker) UpdatePackages(goModPath string, strategy string, dryRun bool) error {
	fmt.Println("Updating packages...")

	var cmd *exec.Cmd
	if dryRun {
		cmd = exec.Command("go", "get", "-u", "./...", "-dry-run")
	} else {
		cmd = exec.Command("go", "get", "-u", "./...")
	}
	cmd.Dir = goModPath

	out, err := cmd.CombinedOutput()
	if err != nil {
		if dryRun {
			return fmt.Errorf("dry-run failed: %w\n%s", err, out)
		}
		return fmt.Errorf("failed to update packages: %w\n%s", err, out)
	}

	if dryRun {
		fmt.Println("\nPackages that would be updated (dry-run):")
	} else {
		fmt.Println("\nPackages updated:")
	}

	lines := strings.Split(string(out), "\n")
	updated := false
	for _, line := range lines {
		if strings.HasPrefix(line, "go get:") || strings.Contains(line, "->") {
			fmt.Printf("  %s\n", line)
			updated = true
		}
	}

	if !updated {
		fmt.Println("  No packages to update")
	}

	if !dryRun {
		fmt.Println("\nRunning go mod tidy...")
		tidy := exec.Command("go", "mod", "tidy")
		tidy.Dir = goModPath
		if out, err := tidy.CombinedOutput(); err != nil {
			return fmt.Errorf("go mod tidy failed: %w\n%s", err, out)
		}
		fmt.Println("  Done.")
	}

	return nil
}

func (vc *VersionChecker) ListDockerUpdates(imageName string) error {
	fmt.Printf("Fetching Docker image versions for %s...\n", imageName)

	parts := strings.Split(imageName, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid image format, use name:tag (e.g., golang:1.25)")
	}

	image := parts[0]
	tag := parts[1]

	tags, err := vc.fetchDockerHubTags(image)
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %w", err)
	}

	fmt.Printf("\nAvailable versions for %s:\n", image)
	currentMajor := vc.extractMajorVersion(tag)

	var majorTags []string
	for _, t := range tags {
		if strings.HasPrefix(t, currentMajor+".") {
			majorTags = append(majorTags, t)
		}
	}

	if len(majorTags) == 0 {
		majorTags = tags[:min(10, len(tags))]
	}

	sort.Sort(sort.Reverse(sort.StringSlice(majorTags)))
	for _, t := range majorTags[:min(10, len(majorTags))] {
		marker := ""
		if t == tag {
			marker = " (current)"
		}
		fmt.Printf("  %s%s\n", t, marker)
	}

	return nil
}

func (vc *VersionChecker) fetchDockerHubTags(image string) ([]string, error) {
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags?page_size=25", image)
	if strings.Contains(image, "/") {
		image = strings.ReplaceAll(image, "/", "%2F")
		url = fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags?page_size=25", image)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return vc.fetchGHCRTags(image)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var tags []string
	for _, r := range result.Results {
		tags = append(tags, r.Name)
	}

	return tags, nil
}

func (vc *VersionChecker) fetchGHCRTags(image string) ([]string, error) {
	url := fmt.Sprintf("https://ghcr.io/v2/%s/tags/list", image)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tags []string `json:"tags"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Tags, nil
}

func (vc *VersionChecker) extractMajorVersion(tag string) string {
	re := regexp.MustCompile(`^(\d+)`)
	m := re.FindStringSubmatch(tag)
	if m != nil {
		return m[1]
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
