package version

// Version components
const (
	Maj = "0"
	Min = "7"
	Fix = "5"
)

var (
	// Must be a string because scripts like dist.sh read this file.
	Version = "0.7.5"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
