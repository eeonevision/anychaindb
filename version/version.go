package version

// Version related constants.
const (
	Major = "2"
	Minor = "0"
	Fix   = "1"
)

var (
	// Version is the current version of AnychainDB platform.
	Version = "2.0.1"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
