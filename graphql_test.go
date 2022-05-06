package main

import (
	"fmt"
	"testing"
)

func TestGraphQL(t *testing.T) {
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

	branchId := "master"

	bookTypeId, err := session.CreateNode(repositoryId, branchId, JsonObject{
		"_qname":      "custom:book",
		"_type":       "d:type",
		"type":        "object",
		"description": "Node List Test Book Type",
		"properties": JsonObject{
			"title": JsonObject{
				"type": "string",
			},
			"description": JsonObject{
				"type": "string",
			},
			"author": JsonObject{
				"type": "string",
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	bookType, _ := session.ReadNode(repositoryId, branchId, bookTypeId)
	fmt.Printf("%v\n", bookType)

	schema, err := session.GraphQLSchema(repositoryId, branchId)
	if err != nil {
		t.Fatal(err)
	}
	if schema == "" {
		t.Fatal("failed to find schema")
	}

	fmt.Println(schema)

	_, err = session.CreateNode(repositoryId, branchId, JsonObject{
		"title":       "hello",
		"description": "this is a book about salutations",
		"author":      "mr bean",
		"_type":       "custom:book",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = session.CreateNode(repositoryId, branchId, JsonObject{
		"title":       "goodbye",
		"description": "leave",
		"author":      "guillermo del toro",
		"_type":       "custom:book",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	query := `
		query {
		custom_books {
				title
				author
			}
		}
	`

	result, err := session.GraphQLQuery(repositoryId, branchId, query, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result["data"] == nil {
		t.Fatal("No graphql response data")
	}

	dataObj := result.GetObject("data")
	if dataObj == nil {
		t.Fatal("data not an object")
	}

	if dataObj["custom_books"] == nil {
		t.Fatal("No books in response data")
	}
	booksArr := dataObj["custom_books"].([]interface{})
	if booksArr == nil {
		t.Fatal("no custom_books in response")
	}

	if len(booksArr) != 2 {
		t.Fatal("wrong number of books in response data")
	}
}
