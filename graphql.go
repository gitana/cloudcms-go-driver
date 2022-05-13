package cloudcms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

func (session *CloudCmsSession) GraphQLQuery(repositoryId string, branchId string, query string, operationName string, variables JsonObject) (JsonObject, error) {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/graphql", repositoryId, branchId)
	params := url.Values{"query": []string{query}}
	if operationName != "" {
		params.Add("operationName", operationName)
	}

	if variables != nil {
		variablesStr, _ := json.Marshal(variables)
		params.Add("variables", string(variablesStr))
	}

	res, err := session.Get(uri, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (session *CloudCmsSession) GraphQLSchema(repositoryId string, branchId string) (string, error) {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/graphql/schema", repositoryId, branchId)
	reader, err := session.Download(uri, nil)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
