package main

import (
	"fmt"
	"testing"
)

func TestRepositories(t *testing.T) {
	session, err := ConnectDefault()
	if err != nil {
		t.Fatal(err)
	}

	repository, err := session.CreateRepository(nil)
	if err != nil {
		t.Fatal(err)
	}

	repository, err = session.ReadRepository(ExtractId(&repository))
	if err != nil {
		t.Fatal(err)
	}
	if repository == nil {
		t.Fatal("No repository")
	}

	branches, err := session.QueryBranches(ExtractId(&repository), make(JsonObject), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, branch := range branches.rows {
		node, err := session.QueryOneNode(ExtractId(&repository), ExtractId(&branch), nil)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ExtractId(&node))

		nodes, err := session.QueryNodes(ExtractId(&repository), ExtractId(&branch), nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		for _, node := range nodes.rows {
			fmt.Printf("Repository: %s, Branch: %s, Node: %s\n", ExtractId(&repository), ExtractId(&branch), ExtractId(&node))
		}
	}

	err = session.DeleteRepository(ExtractId(&repository))
	if err != nil {
		t.Fatal(err)
	}

}
