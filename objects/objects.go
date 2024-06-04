package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"strconv"
	"tiger/repo"
)

type GitObject interface {
	GetFormat() string
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

func Read(r *repo.Repo, sha string) (GitObject, error) {
	path, err := repo.File(r, false, "objects", sha[:2], sha[2:])
	if err != nil {
		return nil, err
	}

	fiObject, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fiObject.IsDir() {
		return nil, err
	}

	object, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	reader, err := zlib.NewReader(object)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	rawBytes := make([]byte, 0)
	_, err = reader.Read(rawBytes)
	if err != nil {
		return nil, err
	}

	x := bytes.IndexByte(rawBytes, ' ')
	format := string(rawBytes[0:x])

	y := bytes.IndexByte(rawBytes, '\x00')
	size, err := strconv.Atoi(string(rawBytes[x:y]))
	if err != nil {
		return nil, err
	}

	if size != len(rawBytes)-y-1 {
		return nil, errors.New(fmt.Sprintf("Malformed object %s: bad length", sha))
	}

	switch format {
	case "commit":
		return newGitCommit(rawBytes[y+1:]), nil
	case "tree":
		return newGitTree(rawBytes[y+1:]), nil
	case "tag":
		return newGitTag(rawBytes[y+1:]), nil
	case "blob":
		return newGitBlob(rawBytes[y+1:]), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown type %s for object %s", format, sha))
	}
}

func Write(o GitObject, r *repo.Repo) (string, error) {
	data, err := o.Serialize()
	if err != nil {
		return "", err
	}

	size := strconv.Itoa(len(data))

	buf := bytes.Buffer{}
	buf.WriteString(o.GetFormat())
	buf.WriteByte(' ')
	buf.WriteString(size)
	buf.WriteByte('\x00')
	buf.Write(data)

	result := make([]byte, buf.Len())
	buf.Read(result)

	shaBytes := sha1.Sum(result)
	sha := string(shaBytes[:])

	if r != nil {
		path, err := repo.File(r, true, "objects", sha[:2], sha[2:])
		if err != nil {
			return "", err
		}

		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return "", errors.New(fmt.Sprintf("object does not exist in %s", path))
			}

			return "", err
		}

		if fi.IsDir() {
			return "", errors.New(fmt.Sprintf("object %s is dir", path))
		}

		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		writer := zlib.NewWriter(f)
		writer.Write(result)
		writer.Close()
	}

	return sha, nil
}

func Find(r *repo.Repo, object, format string, follow bool) (string, error) {
	return object, nil
}
