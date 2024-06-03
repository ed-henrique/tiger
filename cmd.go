package main

import (
	"errors"
	"flag"
	"fmt"
)

func cli() {
	flag.Parse()
	args := flag.Args()[1:]

	var err error
	switch flag.Args()[0] {
	case "add":
		//cmd_add(args)
	case "cat-file":
		//cmd_cat_file(args)
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

func cmdInit(args []string) error {
	if len(args) != 1 {
		return errors.New("Only one path accepted as argument for init")
	}

	_, err := repoCreate(args[0])
	return err
}
