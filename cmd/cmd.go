package cmd

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"tiger/objects"
	"tiger/repo"
)

func Cli() {
	flag.Parse()
	args := flag.Args()[1:]

	var err error
	switch flag.Args()[0] {
	case "add":
		//cmdAdd(args)
	case "cat-file":
		err = cmdCatFile(args)
	case "check-ignore":
		//cmdCheckIgnore(args)
	case "checkout":
		//cmdCheckout(args)
	case "commit":
		//cmdCommit(args)
	case "hash-object":
		//cmdHashObject(args)
	case "init":
		err = cmdInit(args)
	case "log":
		//cmdLog(args)
	case "ls-files":
		//cmdLsFiles(args)
	case "ls-tree":
		//cmdLsTree(args)
	case "rev-parse":
		//cmdRevParse(args)
	case "rm":
		//cmdRm(args)
	case "show-ref":
		//cmdShowRef(args)
	case "status":
		//cmdStatus(args)
	case "tag":
		//cmdTag(args)
	default:
		fmt.Println("Bad command")
	}

	if err != nil {
		panic(err)
	}
}

// cmdInit behaves the same as git init.
func cmdInit(args []string) error {
	switch len(args) {
	case 0:
		_, err := repo.Create(".")
		return err
	case 1:
		_, err := repo.Create(args[0])
		return err
	default:
		return errors.New("init [path]")
	}
}

// cmdCatFile behaves the same as git cat-file.
func cmdCatFile(args []string) error {
	if len(args) != 2 {
		return errors.New("cat-file TYPE OBJECT")
	}

	if !slices.Contains([]string{"blob", "commit", "tag", "tree"}, args[0]) {
		return errors.New("Invalid TYPE, use one of 'blob', 'commit', 'tag' or 'tree'")
	}
	
	objectType := args[0]
	objectSha := args[1]

	r, err := repo.FindRoot(".", true)
	if err != nil {
		return err
	}

	err = catFile(r, objectSha, objectType)
	if err != nil {
		return err
	}

	return nil
}

func catFile(r *repo.Repo, object string, format string) error {
	foundObject, err := objects.Find(r, object, format, true)
	newObject, err := objects.Read(r, foundObject)

	objectSha, err := newObject.Serialize()
	if err != nil {
		return err
	}

	fmt.Println(string(objectSha))
	return nil
}
