package cloudcms

import (
	"fmt"
	"testing"
	"time"
)

func TestProjects(t *testing.T) {
	session, err := ConnectDefault()
	if err != nil {
		t.Fatal(err)
	}

	title := fmt.Sprintf("Project-%d", time.Now().Unix())
	jobId, err := session.StartCreateProject(JsonObject{"title": title})
	if err != nil {
		t.Fatal(err)
	}

	err = session.WaitForJob(jobId)
	if err != nil {
		t.Fatal(err)
	}
	job, _ := session.ReadJob(jobId)
	projectId := job.GetString("created-project-id")

	project, err := session.ReadProject(projectId)
	if err != nil {
		t.Fatal(err)
	}
	if project.GetString("title") != title {
		t.Fatal("Project failed to create/read")
	}
}
