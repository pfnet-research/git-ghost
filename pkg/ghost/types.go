package ghost

type Commit struct {
	BaseCommitHash string
	Commits        []string
	Diff           string
}
