package main

import "testing"

func TestRepoPath(t *testing.T) {
	testRepo := repository{
		gitDir: "a",
	}

	got := repoPath(testRepo, "test", "here")
	expected := "a/test/here"

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}

func TestRepoFile(t *testing.T) {
	testRepo := repository{
		gitDir: ".",
	}

	got, err := repoFile(testRepo, false, "test", "here", "a")
	expected := "test/here/a"

	if err != nil {
		t.Fatalf("got errors %v", err)
	}

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}
