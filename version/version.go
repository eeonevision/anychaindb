package version

const Major = "1"
const Minor = "0"
const Fix = "0"

var (
	// Version is the current version of Leadschain platform
	Version = "1.0.0"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
