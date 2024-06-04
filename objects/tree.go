package objects

type GitTree struct{}

func newGitTree([]byte) GitTree              { return GitTree{} }
func (g GitTree) GetFormat() string         { return "" }
func (g GitTree) Serialize() ([]byte, error) { return nil, nil }
func (g GitTree) Deserialize([]byte) error   { return nil }
