package main

import (
	"fmt"
	"net/url"
)

func (session *CloudCmsSession) QueryBranches(repositoryId string, query JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	res, err := session.Post(fmt.Sprintf("/repositories/%s/branches/query", repositoryId), ToParams(pagination), MapToReader(query))
	if err != nil {
		return nil, err
	}

	return ToResultMap[JsonObject](res), nil
}

func (session *CloudCmsSession) ReadBranch(repositoryId string, branchId string) (JsonObject, error) {
	return session.Get(fmt.Sprintf("/repositories/%s/branches/%s", repositoryId, branchId), nil)
}

func (session *CloudCmsSession) CreateBranch(repositoryId string, parentBranchId string, changesetId string, obj JsonObject) (JsonObject, error) {
	params := url.Values{}
	if changesetId != "" {
		params.Add("changeset", changesetId)
	}
	if parentBranchId != "" {
		params.Add("branch", parentBranchId)
	}

	return session.Post(fmt.Sprintf("/repositories/%s/branches", repositoryId), params, MapToReader(obj))
}

func (session *CloudCmsSession) ListBranches(repositoryId string, pagination JsonObject) (*ResultMap[JsonObject], error) {
	params := ToParams(pagination)
	params.Add("full", "true")
	res, err := session.Get(fmt.Sprintf("/repositories/%s/branches", repositoryId), params)
	if err != nil {
		return nil, err
	}

	return ToResultMap[JsonObject](res), nil
}

func (session *CloudCmsSession) DeleteBranch(repositoryId string, branchId string) error {
	_, err := session.Delete(fmt.Sprintf("/repositories/%s/branches/%s", repositoryId, branchId), nil)
	return err
}

func (session *CloudCmsSession) UpdateBranch(repositoryId string, branchObj JsonObject) (JsonObject, error) {
	doc := branchObj["_doc"]

	// Ensure branch id is a string
	switch doc.(type) {
	case string:
		break
	default:
		return nil, fmt.Errorf("failed to determine branch ID: %v", branchObj)
	}

	return session.Put(fmt.Sprintf("/repositories/%s/branches/%s", repositoryId, doc.(string)), nil, MapToReader(branchObj))
}

func (session *CloudCmsSession) StartResetBranch(repositoryId string, branchId string, changesetId string) (string, error) {
	params := url.Values{"id": []string{changesetId}}
	res, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/reset/start", repositoryId, branchId), params, nil)
	if err != nil {
		return "", err
	}

	return ExtractId(&res), nil
}

func (session *CloudCmsSession) StartChangesetHistory(repositoryId string, branchId string, config JsonObject) (string, error) {
	res, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/history/start", repositoryId, branchId), ToParams(config), nil)
	if err != nil {
		return "", err
	}

	return ExtractId(&res), nil
}
