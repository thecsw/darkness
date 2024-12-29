package alpha

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

const (
	sshRemotePattern = "git@([^:]+):([^:]+)"
)

var (
	sshRemoteRegexp = regexp.MustCompile(sshRemotePattern)
)

// extractGitRemote gets the remote and path, like github.com and thecsw/repo.
func extractGitRemote(workDir string) (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = string(workDir)

	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("getting git origin url: %v", err)
	}

	outString := strings.TrimSuffix(strings.TrimSpace(string(out)), ".git")

	if sshRemoteRegexp.MatchString(outString) {
		matches := sshRemoteRegexp.FindAllStringSubmatch(outString, 1)
		return matches[0][1], matches[0][2], nil
	}

	url, err := url.Parse(outString)
	if err == nil {
		return url.Hostname(), url.Path, nil
	}

	return "", "", errors.New("couldn't find a valid git remote")
}

// extractGitBranch extracts the current working git branch.
func extractGitBranch(workDir string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = string(workDir)

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting current git branch: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}
