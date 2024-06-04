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
		//cmd_add(args)
	case "cat-file":
		err = cmdCatFile(args)
	case "check-ignore":
		//cmd_check_ignore(args)
	case "checkout":
		//cmd_checkout(args)
	case "commit":
		//cmd_commit(args)
	case "hash-object":
		//cmd_hash_object(args)
	case "init":
		err = cmdInit(args)
	case "log":
		//cmd_log(args)
	case "ls-files":
		//cmd_ls_files(args)
	case "ls-tree":
		//cmd_ls_tree(args)
	case "rev-parse":
		//cmd_rev_parse(args)
	case "rm":
		//cmd_rm(args)
	case "show-ref":
		//cmd_show_ref(args)
	case "status":
		//cmd_status(args)
	case "tag":
		//cmd_tag(args)
	default:
		fmt.Println("Bad command")
	}

	if err != nil {
		panic(err)
	}
}

// cmdInit behaves the same as git init.
func cmdInit(args []string) error {
	if len(args) != 1 {
		return errors.New("Only one path accepted as argument for init")
	}

	_, err := repo.Create(args[0])
	return err
}

// cmdCatFile behaves the same as git cat-file.
func cmdCatFile(args []string) error {
	if len(args) != 2 {
		return errors.New("Please provide cat-file TYPE OBJECT")
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
