package version

import (
	"fmt"
	"runtime"
)

const versionFormat = "Papetier scraper version %s (commit: %s, date: %s, go version: %s, platform: %s/%s)\n"

var Version = "development"

var Arch string
var CommitShortHash string
var Os string
var Time string

func String() string {
	return fmt.Sprintf(versionFormat, Version, CommitShortHash, Time, runtime.Version(), Os, Arch)
}
