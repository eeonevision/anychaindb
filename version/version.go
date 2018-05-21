package version

// Version related constants.
const (
	Major = "1"
	Minor = "2"
	Fix   = "3"
)

var (
	// Version is the current version of Leadschain platform.
	Version = "1.2.3"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
