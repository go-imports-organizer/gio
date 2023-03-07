package version

import (
	"runtime/debug"
)

var Version string

func Get() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return "unknown"
}
