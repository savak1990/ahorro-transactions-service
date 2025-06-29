package buildinfo

import (
	_ "embed"
	"encoding/json"
	"time"
)

// BuildInfo represents build metadata
type BuildInfo struct {
	Version     string    `json:"version"`
	GitBranch   string    `json:"gitBranch"`
	GitCommit   string    `json:"gitCommit"`
	GitShort    string    `json:"gitShort"`
	BuildTime   time.Time `json:"buildTime"`
	BuildUser   string    `json:"buildUser"`
	GoVersion   string    `json:"goVersion"`
	ServiceName string    `json:"serviceName"`
}

//go:embed build-info.json
var buildInfoJSON []byte

var buildInfo *BuildInfo

// Get returns the current build information
func Get() *BuildInfo {
	if buildInfo == nil {
		buildInfo = &BuildInfo{
			Version:     "unknown",
			GitBranch:   "unknown",
			GitCommit:   "unknown",
			GitShort:    "unknown",
			BuildTime:   time.Time{},
			BuildUser:   "unknown",
			GoVersion:   "unknown",
			ServiceName: "ahorro-transactions-service",
		}

		// Try to parse embedded build info
		if len(buildInfoJSON) > 0 {
			var parsed BuildInfo
			if err := json.Unmarshal(buildInfoJSON, &parsed); err == nil {
				buildInfo = &parsed
			}
		}
	}
	return buildInfo
}

// GetVersion returns just the version string
func GetVersion() string {
	return Get().Version
}

// GetGitInfo returns git branch and commit
func GetGitInfo() (branch, commit string) {
	info := Get()
	return info.GitBranch, info.GitCommit
}
