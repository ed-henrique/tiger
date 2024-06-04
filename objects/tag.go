package objects

type GitTag struct{}

func newGitTag([]byte) GitTag               { return GitTag{} }
func (g GitTag) GetFormat() string         { return "" }
func (g GitTag) Serialize() ([]byte, error) { return nil, nil }
func (g GitTag) Deserialize([]byte) error   { return nil }
