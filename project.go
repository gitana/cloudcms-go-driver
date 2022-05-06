package main

import "fmt"

// readProject(project: TypedID|string, callback?: ResultCb<PlatformObject>): Promise<PlatformObject>
// createProject(obj: Object, callback?: ResultCb<StartJobResult>):  Promise<StartJobResult>

func (session *CloudCmsSession) ReadProject(projectId string) (JsonObject, error) {
	return session.Get(fmt.Sprintf("/projects/%s", projectId), nil)
}

func (session *CloudCmsSession) StartCreateProject(obj JsonObject) (string, error) {
	res, err := session.Post("/projects/start", nil, MapToReader(obj))

	if err != nil {
		return "", err
	}

	return ExtractId(&res), nil
}
