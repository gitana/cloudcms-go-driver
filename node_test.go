package main

import (
	"testing"
	"time"
)

func setupTestRepository(t *testing.T) (*CloudCmsSession, JsonObject) {
	session, err := ConnectDefault()
	if err != nil {
		t.Fatal(err)
	}

	repository, err := session.CreateRepository(nil)
	if err != nil {
		t.Fatal(err)
	}

	return session, repository
}

func contains(id string, arr []JsonObject) bool {
	for _, v := range arr {
		if ExtractId(&v) == id {
			return true
		}
	}

	return false
}
func TestNodeCrud(t *testing.T) {
	session, repository := setupTestRepository(t)

	repositoryId := ExtractId(&repository)
	branchId := "master"
	defer session.DeleteRepository(repositoryId)

	branch, err := session.ReadBranch(repositoryId, "master")
	if err != nil {
		t.Fatal(err)
	}

	nodeId, err := session.CreateNode(repositoryId, branchId, JsonObject{"title": "MyNode"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	node, err := session.ReadNode(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}

	if nodeId != ExtractId(&node) {
		t.Fatal("read node not equivalent")
	}

	node["title"] = "blah"
	nodeUpdated, err := session.UpdateNode(repositoryId, branchId, node)
	if err != nil {
		t.Fatal(err)
	}

	node, err = session.ReadNode(repositoryId, branchId, ExtractId(&nodeUpdated))
	if err != nil {
		t.Fatal(err)
	}
	if node["title"] != "blah" {
		t.Fatal("update failed")
	}

	// Test change qname
	err = session.ChangeNodeQName(repositoryId, ExtractId(&branch), nodeId, "my:specialNode")
	if err != nil {
		t.Fatal(err)
	}
	node, err = session.ReadNode(repositoryId, branchId, ExtractId(&nodeUpdated))
	if err != nil {
		t.Fatal(err)
	}
	if node["_qname"] != "my:specialNode" {
		t.Fatal("change qname failed")
	}

	err = session.DeleteNode(repositoryId, branchId, ExtractId(&node))
	if err != nil {
		t.Fatal(err)
	}

	node, err = session.ReadNode(repositoryId, branchId, ExtractId(&node))
	if err == nil {
		t.Fatal("deleted node should 404")
	}
}

func TestNodeSearchFind(t *testing.T) {
	session, repository := setupTestRepository(t)

	repositoryId := ExtractId(&repository)
	branchId := "master"
	defer session.DeleteRepository(repositoryId)

	node1Obj := JsonObject{
		"title": "Cheese burger",
		"meal":  "lunch",
	}
	node2Obj := JsonObject{
		"title": "Ham burger",
		"meal":  "lunch",
	}
	node3Obj := JsonObject{
		"title": "Turkey sandwich",
		"meal":  "lunch",
	}
	node4Obj := JsonObject{
		"title": "Oatmeal",
		"meal":  "breakfast",
	}

	node1Id, err := session.CreateNode(repositoryId, branchId, node1Obj, nil)
	if err != nil {
		t.Fatal(err)
	}
	node2Id, err := session.CreateNode(repositoryId, branchId, node2Obj, nil)
	if err != nil {
		t.Fatal(err)
	}
	node3Id, err := session.CreateNode(repositoryId, branchId, node3Obj, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = session.CreateNode(repositoryId, branchId, node4Obj, nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 10)

	queryNodes, err := session.QueryNodes(repositoryId, branchId, JsonObject{"meal": "lunch"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(queryNodes.rows) != 3 {
		t.Fatal("wrong number of query results")
	}
	if !contains(node1Id, queryNodes.rows) {
		t.Fatal("query didn't return node1")
	}
	if !contains(node2Id, queryNodes.rows) {
		t.Fatal("query didn't return node2")
	}
	if !contains(node3Id, queryNodes.rows) {
		t.Fatal("query didn't return node3")
	}

	find := JsonObject{"search": "burger"}
	findNodes, err := session.FindNodes(repositoryId, branchId, find, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(findNodes.rows) != 2 {
		t.Fatal("wrong number of find results")
	}
	if !contains(node1Id, findNodes.rows) {
		t.Fatal("find didn't return node1")
	}
	if !contains(node2Id, findNodes.rows) {
		t.Fatal("find didn't return node2")
	}

	searchNodes, err := session.SearchNodes(repositoryId, branchId, "burger", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(searchNodes.rows) != 2 {
		t.Fatal("wrong number of search results")
	}
	if !contains(node1Id, searchNodes.rows) {
		t.Fatal("search didn't return node1")
	}
	if !contains(node2Id, searchNodes.rows) {
		t.Fatal("search didn't return node2")
	}
}

func TestFeatures(t *testing.T) {
	session, repository := setupTestRepository(t)

	repositoryId := ExtractId(&repository)
	branchId := "master"
	defer session.DeleteRepository(repositoryId)

	nodeId, err := session.CreateNode(repositoryId, branchId, JsonObject{"title": "node"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = session.AddNodeFeature(repositoryId, branchId, nodeId, "f:filename", JsonObject{"filename": "node"})
	if err != nil {
		t.Fatal(err)
	}

	node, err := session.ReadNode(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}
	nodeFeatures := node.GetObject("_features")
	if nodeFeatures.GetObject("f:filename") == nil {
		t.Fatal("filename failed to add")
	}

	err = session.RemoveNodeFeature(repositoryId, branchId, nodeId, "f:filename")
	if err != nil {
		t.Fatal(err)
	}

	node, err = session.ReadNode(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}
	nodeFeatures = node.GetObject("_features")
	if nodeFeatures.GetObject("f:filename") != nil {
		t.Fatal("filename failed to remove")
	}
}

func TestVersions(t *testing.T) {
	session, repository := setupTestRepository(t)

	repositoryId := ExtractId(&repository)
	branchId := "master"
	defer session.DeleteRepository(repositoryId)

	nodeId, err := session.CreateNode(repositoryId, branchId, JsonObject{"title": "node"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	node, err := session.ReadNode(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}

	systemObj := node.GetObject("_system")
	firstChangeset := systemObj.GetString("changeset")
	if firstChangeset == "" {
		t.Fatal("missing initial changeset")
	}

	node["title"] = "new stuff"
	_, err = session.UpdateNode(repositoryId, branchId, node)
	if err != nil {
		t.Fatal(err)
	}
	node, err = session.ReadNode(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}

	versions, err := session.ListVersions(repositoryId, branchId, nodeId, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(versions.rows) != 2 {
		t.Fatal("invalid version count")
	}

	version1, err := session.ReadVersion(repositoryId, branchId, nodeId, firstChangeset, nil)
	if err != nil {
		t.Fatal(err)
	}
	if version1.GetString("title") != "node" {
		t.Fatal("wrong intial version title")
	}

	restoredVersion, err := session.RestoreVersion(repositoryId, branchId, nodeId, firstChangeset)
	if err != nil {
		t.Fatal(err)
	}
	if restoredVersion.GetString("title") != "node" {
		t.Fatal("wrong restored version title")
	}
}

func createFile(t *testing.T, session *CloudCmsSession, repositoryId string, branchId string, obj JsonObject, parentPath string, isFolder bool) string {
	nodeId, err := session.CreateNode(repositoryId, branchId, obj, map[string]string{"parentPath": parentPath})
	if err != nil {
		t.Fatal(err)
	}

	if isFolder {
		session.AddNodeFeature(repositoryId, branchId, nodeId, "f:container", JsonObject{})
	}

	return nodeId
}

func TestTraverse(t *testing.T) {
	session, repository := setupTestRepository(t)

	repositoryId := ExtractId(&repository)
	branchId := "master"
	defer session.DeleteRepository(repositoryId)

	_ = createFile(t, session, repositoryId, branchId, JsonObject{"title": "folder1"}, "/", true)
	_ = createFile(t, session, repositoryId, branchId, JsonObject{"title": "file1"}, "/", false)
	_ = createFile(t, session, repositoryId, branchId, JsonObject{"title": "folder2"}, "/folder1", true)
	_ = createFile(t, session, repositoryId, branchId, JsonObject{"title": "file2"}, "/folder1", false)
	_ = createFile(t, session, repositoryId, branchId, JsonObject{"title": "file3"}, "/folder1/folder2", false)
	_ = createFile(t, session, repositoryId, branchId, JsonObject{"title": "file4"}, "/folder1", false)
	file5Id := createFile(t, session, repositoryId, branchId, JsonObject{"title": "file5"}, "/folder1/folder2", false)

	// test path resolves
	path, err := session.ResolveNodePath(repositoryId, branchId, file5Id)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/folder1/folder2/file5" {
		t.Fatal("incorrect resolved path")
	}

	paths, err := session.ResolveNodePaths(repositoryId, branchId, file5Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) < 1 {
		t.Fatal("failed to resolve paths")
	}

	traverse := JsonObject{
		"depth":  1,
		"filter": "ALL_BUT_START_NODE",
		"associations": JsonObject{
			"a:child": "ANY",
		},
	}
	time.Sleep(time.Second * 5)

	rootNode, _ := session.ReadNode(repositoryId, branchId, "root")

	results, err := session.TraverseNode(repositoryId, branchId, ExtractId(&rootNode), traverse)
	if err != nil {
		t.Fatal(err)
	}
	nodes := results.GetObject("nodes")
	if len(nodes) != 2 {
		t.Fatal("failed to traverse nodes")
	}
	associations := results.GetObject("associations")
	if len(associations) != 2 {
		t.Fatal("failed to traverse associations")
	}
}
