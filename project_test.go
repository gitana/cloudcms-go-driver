package main

import "testing"

func TestProjects(t *testing.T) {
	session, err := ConnectDefault()
	if err != nil {
		t.Fatal(err)
	}

	jobId, err := session.StartCreateProject(JsonObject{"title": "my project"})
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
	if project.GetString("title") != "my project" {
		t.Fatal("Project failed to create/read")
	}
}
