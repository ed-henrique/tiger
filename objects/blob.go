package objects

type GitBlob struct {
	Format string
	Data   []byte
}

func newGitBlob(data []byte) *GitBlob {
	gb := &GitBlob{
		Format: "blob",
		Data:   data,
	}

	return gb
}

func (g *GitBlob) GetFormat() string {
	return g.Format
}

func (g *GitBlob) Serialize() ([]byte, error) {
	return g.Data, nil
}

func (g *GitBlob) Deserialize(data []byte) error {
	g.Data = data
	return nil
}
