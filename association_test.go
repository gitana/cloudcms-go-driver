package cloudcms

import "testing"

func TestAssociations(t *testing.T) {
	session, err := ConnectDefault()
	if err != nil {
		t.Fatal(err)
	}

	repository, err := session.CreateRepository(nil)
	if err != nil {
		t.Fatal(err)
	}

	repositoryId := ExtractId(&repository)
	branchId := "master"
	defer session.DeleteRepository(repositoryId)

	node1Id, _ := session.CreateNode(repositoryId, branchId, JsonObject{"title": "node 1"}, nil)
	node2Id, _ := session.CreateNode(repositoryId, branchId, JsonObject{"title": "node 2"}, nil)
	node3Id, _ := session.CreateNode(repositoryId, branchId, JsonObject{"title": "node 3"}, nil)

	association1, err := session.Associate(repositoryId, branchId, node1Id, node2Id, "a:child", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if association1.GetString("directionality") != "DIRECTED" {
		t.Fatal("wrong directionality")
	}
	if association1.GetString("source") != node1Id {
		t.Fatal("wrong source")
	}
	if association1.GetString("target") != node2Id {
		t.Fatal("wrong target")
	}

	// association2 = node1.associate(node3, QName.create('a:linked'), directionality=Directionality.UNDIRECTED, data={'test': 'field'})
	association2, err := session.Associate(repositoryId, branchId, node1Id, node3Id, "a:linked", "UNDIRECTED", JsonObject{"test": "field"})
	if err != nil {
		t.Fatal(err)
	}

	if association2.GetString("source") != node1Id {
		t.Fatal("wrong source")
	}
	if association2.GetString("target") != node3Id {
		t.Fatal("wrong target")
	}
	if association2.GetString("source") != node1Id {
		t.Fatal("wrong source")
	}

	node1Associations, err := session.ListNodeAssociations(repositoryId, branchId, node1Id, "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if node1Associations.size != 3 {
		t.Fatal("node1 has wrong number of associations")
	}

	node1Outgoing, err := session.ListOutgoingAssociations(repositoryId, branchId, node1Id, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if node1Outgoing.size != 2 {
		t.Fatal("node1 has wrong number of outgoing associations")
	}

	node1Incoming, err := session.ListIncomingAssociations(repositoryId, branchId, node1Id, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if node1Incoming.size != 2 {
		t.Fatal("node1 has wrong number of incoming associations")
	}

	node1Children, err := session.QueryNodeChildren(repositoryId, branchId, node1Id, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if node1Children.size != 1 {
		t.Fatal("node1 has wrong number of children")
	}

	err = session.UnassociateChild(repositoryId, branchId, node1Id, node2Id)
	if err != nil {
		t.Fatal(err)
	}
	err = session.Unassociate(repositoryId, branchId, node1Id, node3Id, "a:linked", "UNDIRECTED")
	if err != nil {
		t.Fatal(err)
	}

	node1Associations, err = session.ListNodeAssociations(repositoryId, branchId, node1Id, "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if node1Associations.size != 1 {
		t.Fatal("node1 has wrong number of associations")
	}

}
