package alpha

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	sshRemotePattern = "git@([^:]+):([^:]+)"

	// what we pass to --pretty
	gitPretty = "format:%cd"
	// this is what we pass to git --date to generate RFC3339
	rfc3339GitFormat = "format:%Y-%m-%dT%H:%M:%SZ%z"
	rfc3339Pattern   = "2006-01-02T15:04:05Z-0700"
)

var (
	sshRemoteRegexp = regexp.MustCompile(sshRemotePattern)
)

// ExtractGitRemote gets the remote and path, like github.com and thecsw/repo.
func ExtractGitRemote(conf *DarknessConfig) (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = string(conf.Runtime.WorkDir)

	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("getting git origin url: %v", err)
	}

	outString := strings.TrimSuffix(strings.TrimSpace(string(out)), ".git")

	if sshRemoteRegexp.MatchString(outString) {
		matches := sshRemoteRegexp.FindAllStringSubmatch(outString, 1)
		if len(matches) > 0 && len(matches[0]) > 2 {
			return matches[0][1], matches[0][2], nil
		}
	}

	url, err := url.Parse(outString)
	if err == nil {
		return url.Hostname(), url.Path, nil
	}

	return "", "", errors.New("couldn't find a valid git remote")
}

// ExtractGitBranch extracts the current working git branch.
func ExtractGitBranch(conf *DarknessConfig) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = string(conf.Runtime.WorkDir)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting current git branch: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// isPathSafe checks if a path is safe to use in exec.Command
// It validates the path doesn't contain suspicious patterns
func isPathSafe(path string) bool {
	// Clean the path to resolve any . or .. components
	cleanPath := filepath.Clean(path)
	
	// Don't allow absolute paths
	if filepath.IsAbs(cleanPath) {
		return false
	}
	
	// Don't allow paths that try to traverse up directories
	if strings.Contains(cleanPath, "../") {
		return false
	}
	
	// Ensure it doesn't have any suspicious characters
	// Only allow alphanumeric, '_', '-', '/', '.', and space
	safePattern := regexp.MustCompile(`^[a-zA-Z0-9_\-/\. ]+$`)
	return safePattern.MatchString(cleanPath)
}

// ExtractGitLastModified will give the git date of when the file was last modified.
func ExtractGitLastModified(conf *DarknessConfig, path yunyun.RelativePathFile) (time.Time, error) {
	pathStr := string(path)
	
	// Validate the path before using it in exec.Command
	if !isPathSafe(pathStr) {
		return time.Time{}, fmt.Errorf("invalid path provided: %s", path)
	}
	
	cmd := exec.Command("git", "log", "--date", rfc3339GitFormat, "-1", "--pretty="+gitPretty, "--", pathStr) // #nosec G204 - pathStr validated by isPathSafe
	cmd.Dir = string(conf.Runtime.WorkDir)

	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, fmt.Errorf("getting last modified for file %s: %v", path, err)
	}
	return time.Parse(rfc3339Pattern, string(out))
}
