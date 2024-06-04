package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type repository struct {
	worktree string
	gitDir   string
	config   *ini.File
}

func newRepository(path string, force bool) (repository, error) {
	repo := repository{
		worktree: path,
		gitDir: filepath.Join(path, ".git"),
	}

	fi, err := os.Stat(repo.gitDir)
	if os.IsNotExist(err) {
		return repository{}, errors.New(fmt.Sprintf("Not a Git repository %s", path))
	}

	if !(force || fi.IsDir()) {
		return repository{}, errors.New(fmt.Sprintf("Not a Git repository %s", path))
	}

	conf, err := repoFile(repo, false, "config")

	if conf != "" && err == nil {
		if _, err := os.Stat(conf); !os.IsNotExist(err) {
			repo.config, err = ini.Load(conf)

			if err != nil {
				return repository{}, err
			}
		}
	} else if !force {
		return repository{}, errors.New("Configuration file missing")
	}

	if !force {
		core, err := repo.config.GetSection("core")

		if err != nil {
			return repository{}, err
		}

		vers := core.Key("repositoryformatversion")
		value, err := vers.Int()

		if err != nil {
			return repository{}, err
		}

		if value != 0 {
			return repository{}, errors.New(fmt.Sprintf("Unsupported repositoryformatversion %d", value))
		}
	}

	return repo, nil
}

func repoPath(repo repository, path ...string) string {
	return filepath.Join(append([]string{repo.gitDir}, path...)...)
}

func repoFile(repo repository, mkdir bool, path ...string) (string, error) {
	if _, err := repoDir(repo, mkdir, path[:len(path)-1]...); err != nil {
		return "", err
	}

	return repoPath(repo, path...), nil
}

func repoDir(repo repository, mkdir bool, path ...string) (string, error) {
	newPath := repoPath(repo, path...)

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

func repoDefaultConfig() *ini.File {
	config := ini.Empty()
	core, _ := config.NewSection("core")
	core.NewKey("repositoryformatversion", "0")
	core.NewKey("filemode", "false")
	core.NewKey("bare", "false")
	return config
}

func repoCreate(path string) (repository, error) {
	repo, err := newRepository(path, true)

	if err != nil {
		return repository{}, err
	}

	fiWorktree, err := os.Stat(repo.worktree)
	if !os.IsNotExist(err) {
		if !fiWorktree.IsDir() {
			return repository{}, errors.New(fmt.Sprintf("%s is not a directory", path))
		}

		gitDir, err := os.Open(repo.gitDir)
		defer gitDir.Close()
		if !os.IsNotExist(err) {
			if _, err := gitDir.Readdir(1); err == io.EOF {
				return repository{}, errors.New(fmt.Sprintf("%s is not empty", path))
			}
		}
	} else {
		err := os.MkdirAll(repo.worktree, 0777)

		if err != nil {
			return repository{}, err
		}
	}

	if _, err := repoDir(repo, true, "branches"); err != nil {
		return repository{}, err
	}
	if _, err := repoDir(repo, true, "objects"); err != nil {
		return repository{}, err
	}
	if _, err := repoDir(repo, true, "refs", "tags"); err != nil {
		return repository{}, err
	}
	if _, err := repoDir(repo, true, "refs", "heads"); err != nil {
		return repository{}, err
	}

	description, err := repoFile(repo, false, "description")
	if err != nil {
		return repository{}, err
	}

	fDescription, err := os.Create(description)
	if err != nil {
		return repository{}, err
	}
	defer fDescription.Close()
	_, err = fDescription.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")
	if err != nil {
		return repository{}, err
	}

	head, err := repoFile(repo, false, "HEAD")
	if err != nil {
		return repository{}, err
	}

	fHead, err := os.Create(head)
	if err != nil {
		return repository{}, err
	}
	defer fHead.Close()
	_, err = fHead.WriteString("ref: refs/heads/master\n")
	if err != nil {
		return repository{}, err
	}

	config, err := repoFile(repo, false, "config")
	if err != nil {
		return repository{}, err
	}

	conf := repoDefaultConfig()
	err = conf.SaveTo(config)
	if err != nil {
		return repository{}, err
	}

	return repo, nil
}

func repoFind(path string, required bool) (repository, error) {
	if path == "" {
		path = "."
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return repository{}, err
	}

	fiGit, err := os.Stat(filepath.Join(path, ".git"))
	if !os.IsNotExist(err) && fiGit.IsDir() {
		repo, err := newRepository(path, false)
		if err != nil {
			return repository{}, err
		}

		return repo, nil
	}

	parent, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return repository{}, err
	}

	if parent == path {
		if required {
			return repository{}, errors.New("No git directory.")
		}

		return repository{}, nil
	}

	return repoFind(parent, required)
}
