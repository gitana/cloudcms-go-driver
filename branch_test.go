package main

import (
	"fmt"
	"testing"
)

func TestBranches(t *testing.T) {
	session, err := ConnectDefault()
	if err != nil {
		t.Fatal(err)
	}

	repository, err := session.CreateRepository(nil)
	if err != nil {
		t.Fatal(err)
	}

	repositoryId := ExtractId(&repository)
	defer session.DeleteRepository(repositoryId)

	branches, err := session.ListBranches(repositoryId, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, branch := range branches.rows {
		fmt.Printf("Repository: %s, Branch: %s, Title: %s\n", repositoryId, ExtractId(&branch), branch["title"])
	}

	master, err := session.ReadBranch(repositoryId, "master")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Repository: %s, Branch: %s, Title: %s\n", repositoryId, "master", master["title"])

	tip := master.GetString("tip")
	newBranch, err := session.CreateBranch(repositoryId, "master", tip, JsonObject{"title": "new branch 1"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Repository: %s, Branch: %s, Title: %s\n", repositoryId, ExtractId(&newBranch), newBranch["title"])

}
