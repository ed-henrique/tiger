package objects

type GitCommit struct{}

func newGitCommit([]byte) GitCommit            { return GitCommit{} }
func (g GitCommit) GetFormat() string         { return "" }
func (g GitCommit) Serialize() ([]byte, error) { return nil, nil }
func (g GitCommit) Deserialize([]byte) error   { return nil }
