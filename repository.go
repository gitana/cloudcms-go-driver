package main

import "fmt"

func (session *CloudCmsSession) CreateRepository(obj JsonObject) (JsonObject, error) {
	return session.Post("/repositories", nil, MapToReader(obj))
}

func (session *CloudCmsSession) DeleteRepository(repositoryId string) error {
	_, err := session.Delete(fmt.Sprintf("/repositories/%s", repositoryId), nil)
	return err
}

func (session *CloudCmsSession) ReadRepository(repositoryId string) (JsonObject, error) {
	return session.Get(fmt.Sprintf("/repositories/%s", repositoryId), nil)
}

func (session *CloudCmsSession) QueryRepositories(query JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	res, err := session.Post("/repositories/query", ToParams(pagination), MapToReader(query))
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}
