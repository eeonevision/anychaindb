package version

const Major = "1"
const Minor = "2"
const Fix = "2"

var (
	// Version is the current version of Leadschain platform
	Version = "1.2.2"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
