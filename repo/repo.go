package repo

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Repo is a git repository
type Repo struct {
	Worktree string
	GitDir   string
	Config   *ini.File
}

// New creates a new Repo object. If force is set, it will ignore missing directories and
// configurations.
func New(path string, force bool) (*Repo, error) {
	repo := &Repo{
		Worktree: path,
		GitDir:   filepath.Join(path, ".git"),
	}

	if force {
		return repo, nil
	}

	fiGitDir, err := os.Stat(repo.GitDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf(".git directory does not exist in %s", path))
		}

		return nil, err
	}

	if !fiGitDir.IsDir() {
		return nil, errors.New(fmt.Sprintf("not a git repository %s", path))
	}

	conf, err := File(repo, false, "config")

	if conf != "" && err == nil {
		fiConf, err := os.Stat(repo.GitDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, errors.New(fmt.Sprintf("config file missing in %s", repo.GitDir))
			}

			return nil, err
		}

		if fiConf.IsDir() {
			return nil, errors.New(fmt.Sprintf("config is a directory in %s", repo.GitDir))
		}

		if repo.Config, err = ini.Load(conf); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if conf == "" {
		return nil, errors.New(fmt.Sprintf("config file missing in %s", repo.GitDir))
	}

	core, err := repo.Config.GetSection("core")

	if err != nil {
		return nil, err
	}

	vers := core.Key("repositoryformatversion")
	value, err := vers.Int()

	if err != nil {
		return nil, err
	}

	if value != 0 {
		return nil, errors.New(fmt.Sprintf("Unsupported repositoryformatversion %d", value))
	}

	return repo, nil
}

// Path concatenates the given paths to the repository's .git dir.
func Path(repo *Repo, path ...string) string {
	return filepath.Join(append([]string{repo.GitDir}, path...)...)
}

// File checks if the path to a file in .git exists, and if so, returns the
// file path. Be aware, this doesn't mean that the file itself exists, only
// that the directories needed to reach it are available.
//
// If mkdir is passed, it will create the non-existent directories as needed.
func File(repo *Repo, mkdir bool, path ...string) (string, error) {
	if _, err := Dir(repo, mkdir, path[:len(path)-1]...); err != nil {
		return "", err
	}

	return Path(repo, path...), nil
}

// Dir checks whether or not a path inside .git is valid, and if so, returns
// the path. If mkdir is passed, it will create the non-existent directories
// as needed.
func Dir(repo *Repo, mkdir bool, path ...string) (string, error) {
	newPath := Path(repo, path...)

	fi, err := os.Stat(newPath)

	if !os.IsNotExist(err) {
		if fi.IsDir() {
			return newPath, nil
		}

		return "", errors.New(fmt.Sprintf("Not a directory %s", newPath))
	}

	if mkdir {
		if err := os.MkdirAll(newPath, 0777); err != nil {
			return "", err
		}

		return newPath, nil
	}

	return "", nil
}

// DefaultConfig returns the default .git/config configurations.
func DefaultConfig() *ini.File {
	config := ini.Empty()
	core, _ := config.NewSection("core")
	core.NewKey("repositoryformatversion", "0")
	core.NewKey("filemode", "false")
	core.NewKey("bare", "false")
	return config
}

// Create builds the whole .git structure inside the given path.
func Create(path string) (*Repo, error) {
	repo, err := New(path, true)

	if err != nil {
		return nil, err
	}

	fiWorktree, err := os.Stat(repo.Worktree)
	if !os.IsNotExist(err) {
		if !fiWorktree.IsDir() {
			return nil, errors.New(fmt.Sprintf("%s is not a directory", path))
		}

		GitDir, err := os.Open(repo.GitDir)
		defer GitDir.Close()
		if !os.IsNotExist(err) {
			if _, err := GitDir.Readdir(1); err == io.EOF {
				return nil, errors.New(fmt.Sprintf("%s is not empty", path))
			}
		}
	} else {
		err := os.MkdirAll(repo.Worktree, 0777)

		if err != nil {
			return nil, err
		}
	}

	if _, err := Dir(repo, true, "branches"); err != nil {
		return nil, err
	}
	if _, err := Dir(repo, true, "objects"); err != nil {
		return nil, err
	}
	if _, err := Dir(repo, true, "refs", "tags"); err != nil {
		return nil, err
	}
	if _, err := Dir(repo, true, "refs", "heads"); err != nil {
		return nil, err
	}

	description, err := File(repo, false, "description")
	if err != nil {
		return nil, err
	}

	fDescription, err := os.Create(description)
	if err != nil {
		return nil, err
	}
	defer fDescription.Close()
	_, err = fDescription.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")
	if err != nil {
		return nil, err
	}

	head, err := File(repo, false, "HEAD")
	if err != nil {
		return nil, err
	}

	fHead, err := os.Create(head)
	if err != nil {
		return nil, err
	}
	defer fHead.Close()
	_, err = fHead.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return nil, err
	}

	config, err := File(repo, false, "config")
	if err != nil {
		return nil, err
	}

	conf := DefaultConfig()
	err = conf.SaveTo(config)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// FindRoot searches for a git repo in the given path, and if not found,
// in its parents, until a valid one or none is found.
func FindRoot(path string, required bool) (*Repo, error) {
	if path == "" {
		path = "."
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	fiGit, err := os.Stat(filepath.Join(path, ".git"))
	if !os.IsNotExist(err) && fiGit.IsDir() {
		repo, err := New(path, false)
		if err != nil {
			return nil, err
		}

		return repo, nil
	}

	parent, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return nil, err
	}

	if parent == path {
		if required {
			return nil, errors.New("No git directory.")
		}

		return nil, nil
	}

	return FindRoot(parent, required)
}
