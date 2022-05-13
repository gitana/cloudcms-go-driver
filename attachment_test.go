package cloudcms

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestAttachments(t *testing.T) {
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

	nodeId, err := session.CreateNode(repositoryId, branchId, JsonObject{"title": "nodule"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := os.ReadFile("res/cloudcms.png")
	if err != nil {
		t.Fatal(err)
	}

	err = session.UploadAttachment(repositoryId, branchId, nodeId, "default", bytes.NewReader(fileBytes), "image/png", "image.png")
	if err != nil {
		t.Fatal(err)
	}

	dl, err := session.DownloadAttachment(repositoryId, branchId, nodeId, "default")
	if err != nil {
		t.Fatal(err)
	}
	defer dl.Close()

	dlBytes, err := io.ReadAll(dl)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(dlBytes, fileBytes) {
		t.Fatal("Downloaded and uploaded files are not the same!")
	}

	atts, err := session.ListAttachments(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}
	if len(atts.rows) != 1 {
		t.Fatalf("Wrong number of attachments: %d", len(atts.rows))
	}

	err = session.DeleteAttachment(repositoryId, branchId, nodeId, "default")
	if err != nil {
		t.Fatal(err)
	}
	atts, err = session.ListAttachments(repositoryId, branchId, nodeId)
	if err != nil {
		t.Fatal(err)
	}
	if len(atts.rows) != 0 {
		t.Fatalf("Wrong number of attachments: %d", len(atts.rows))
	}

	dl, err = session.DownloadAttachment(repositoryId, branchId, nodeId, "default")
	if err == nil {
		// shouldn't have found attachment
		defer dl.Close()
		t.Fatal("Attachment not successfully deleted")
	}

}
