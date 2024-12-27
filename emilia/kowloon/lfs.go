package kowloon

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/alpha"
)

const (
	// The parameters are in the order of 1) git repo (thecsw/repo) 2) branch name 3) file path
	githubLink         = "github.com"
	githubLfsMediaPath = "https://media.githubusercontent.com/media/%s/refs/heads/%s/%s"

	// This should be absolute against the entire project
	lfsLinkPrefix = "lfs:"
)

// GetLfsMediaPath returns the path to the LFS media file.
func GetLfsMediaPath(conf *alpha.DarknessConfig, path string) string {
	if conf.External.GitRemoteService == githubLink {
		return fmt.Sprintf(githubLfsMediaPath,
			strings.Trim(conf.External.GitRemotePath, "/"),
			conf.External.GitBranch, path)
	}
	conf.Runtime.Logger.Fatalf("unsupported service for linking LFS: %s", conf.External.GitRemoteService)
	return "" // unreachable
}

// ConvertImageToLfsMediaLink takes a link and returns full remote path if it starts with lfs:
func ConvertImageToLfsMediaLink(conf *alpha.DarknessConfig, link string) string {
	if !strings.HasPrefix(link, lfsLinkPrefix) {
		return link
	}
	cleanLink := strings.TrimPrefix(link, lfsLinkPrefix)
	return GetLfsMediaPath(conf, strings.TrimLeft(cleanLink, "/"))
}
