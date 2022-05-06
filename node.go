package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strconv"
)

func (session *CloudCmsSession) ReadNode(repositoryId string, branchId string, nodeId string) (JsonObject, error) {
	return session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s", repositoryId, branchId, nodeId), nil)
}

func (session *CloudCmsSession) QueryNodes(repositoryId string, branchId string, query JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	res, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/query", repositoryId, branchId), ToParams(pagination), MapToReader(query))
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) QueryOneNode(repositoryId string, branchId string, query JsonObject) (JsonObject, error) {

	res, err := session.QueryNodes(repositoryId, branchId, query, JsonObject{"limit": 1})
	if err != nil {
		return nil, err
	}

	if res == nil || len(res.rows) == 0 {
		return nil, nil
	}

	return res.rows[0], nil
}

func (session *CloudCmsSession) SearchNodes(repositoryId string, branchId string, text string, pagination JsonObject) (*ResultMap[JsonObject], error) {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/nodes/search", repositoryId, branchId)
	params := ToParams(pagination)
	params.Add("text", text)

	res, err := session.Get(uri, params)
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) FindNodes(repositoryId string, branchId string, config JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/nodes/find", repositoryId, branchId)
	params := ToParams(pagination)

	res, err := session.Post(uri, params, MapToReader(config))
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) CreateNode(repositoryId string, branchId string, obj JsonObject, opts map[string]string) (string, error) {
	params := url.Values{}
	for key, val := range opts {
		switch key {
		case "rootNodeId":
			params.Add("rootNodeId", val)
		case "parentFolderPath":
			params.Add("parentFolderPath", val)
		case "filePath":
			params.Add("filePath", val)
		case "fileName":
			params.Add("fileName", val)
		case "associationTypeString":
			params.Add("associationTypeString", val)
		case "filename":
			params.Add("fileName", val)
		case "filepath":
			params.Add("filePath", val)
		case "parentPath":
			params.Add("parentFolderPath", val)
		}
	}

	res, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes", repositoryId, branchId), params, MapToReader(obj))
	if err != nil {
		return "", err
	}

	return ExtractId(&res), nil
}

func (session *CloudCmsSession) QueryNodeRelatives(repositoryId string, branchId string, nodeId string, associationTypeQName string, associationDirection string, query JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	params := ToParams(pagination)
	params.Add("type", associationTypeQName)
	params.Add("direction", associationDirection)

	res, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/relatives/query", repositoryId, branchId, nodeId), params, MapToReader(query))
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) QueryNodeChildren(repositoryId string, branchId string, nodeId string, query JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	return session.QueryNodeRelatives(repositoryId, branchId, nodeId, "a:child", "OUTGOING", query, pagination)
}

func (session *CloudCmsSession) ListNodeAssociations(repositoryId string, branchId string, nodeId string, associationTypeQName string, associationDirection string, pagination JsonObject) (*ResultMap[JsonObject], error) {
	params := ToParams(pagination)

	if associationTypeQName != "" {
		params.Add("type", associationTypeQName)
	}
	if associationDirection != "" {
		params.Add("direction", associationDirection)
	}

	res, err := session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/associations", repositoryId, branchId, nodeId), params)
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) ListOutgoingAssociations(repositoryId string, branchId string, nodeId string, associationTypeQName string, pagination JsonObject) (*ResultMap[JsonObject], error) {
	return session.ListNodeAssociations(repositoryId, branchId, nodeId, associationTypeQName, "OUTGOING", pagination)
}

func (session *CloudCmsSession) ListIncomingAssociations(repositoryId string, branchId string, nodeId string, associationTypeQName string, pagination JsonObject) (*ResultMap[JsonObject], error) {
	return session.ListNodeAssociations(repositoryId, branchId, nodeId, associationTypeQName, "INCOMING", pagination)
}

func (session *CloudCmsSession) Associate(repositoryId string, branchId string, nodeId string, otherNodeId string, associationTypeQName string, associationDirection string, obj JsonObject) (JsonObject, error) {
	params := url.Values{"node": []string{otherNodeId}}
	if associationTypeQName != "" {
		params.Add("type", associationTypeQName)
	}
	if associationDirection != "" {
		params.Add("directionality", associationDirection)
	}

	return session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/associate", repositoryId, branchId, nodeId), params, MapToReader(obj))
}

func (session *CloudCmsSession) Unassociate(repositoryId string, branchId string, nodeId string, otherNodeId string, associationTypeQName string, associationDirection string) error {
	params := url.Values{"node": []string{otherNodeId}}
	if associationTypeQName != "" {
		params.Add("type", associationTypeQName)
	}
	if associationDirection != "" {
		params.Add("directionality", associationDirection)
	}

	_, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/unassociate", repositoryId, branchId, nodeId), params, MapToReader(JsonObject{}))
	return err
}

func (session *CloudCmsSession) AssociateChild(repositoryId string, branchId string, nodeId string, otherNodeId string, obj JsonObject) (JsonObject, error) {
	return session.Associate(repositoryId, branchId, nodeId, otherNodeId, "a:child", "DIRECTED", obj)
}

func (session *CloudCmsSession) UnassociateChild(repositoryId string, branchId string, nodeId string, otherNodeId string) error {
	return session.Unassociate(repositoryId, branchId, nodeId, otherNodeId, "a:child", "DIRECTED")
}

func (session *CloudCmsSession) DeleteNode(repositoryId string, branchId string, nodeId string) error {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s", repositoryId, branchId, nodeId)
	_, err := session.Delete(uri, nil)
	return err
}

func (session *CloudCmsSession) DeleteNodes(repositoryId string, branchId string, nodeIds []string) ([]string, error) {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/nodes/delete", repositoryId, branchId)
	body := JsonObject{
		"_docs": nodeIds,
	}

	res, err := session.Post(uri, nil, MapToReader(body))
	if err != nil {
		return nil, err
	}

	return res["_docs"].([]string), nil
}

func (session *CloudCmsSession) UpdateNode(repositoryId string, branchId string, node JsonObject) (JsonObject, error) {

	doc := node["_doc"]

	// Ensure node id is a string
	switch doc.(type) {
	case string:
	default:
		return nil, fmt.Errorf("failed to determine node ID: %v", node)
	}

	return session.Put(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s", repositoryId, branchId, doc.(string)), url.Values{}, MapToReader(node))
}

func (session *CloudCmsSession) PatchNode(repositoryId string, branchId string, nodeId string, patchObj JsonObject) (JsonObject, error) {

	return session.Patch(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s", repositoryId, branchId, nodeId), url.Values{}, MapToReader(patchObj))
}

func (session *CloudCmsSession) AddNodeFeature(repositoryId string, branchId string, nodeId string, featureId string, config JsonObject) error {
	_, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/features/%s", repositoryId, branchId, nodeId, featureId), url.Values{}, MapToReader(config))
	return err
}

func (session *CloudCmsSession) RemoveNodeFeature(repositoryId string, branchId string, nodeId string, featureId string) error {
	_, err := session.Delete(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/features/%s", repositoryId, branchId, nodeId, featureId), url.Values{})
	return err
}

func (session *CloudCmsSession) RefreshNode(repositoryId string, branchId string, node JsonObject) error {
	doc := node["_doc"]

	// Ensure node id is a string
	switch doc.(type) {
	case string:
		break
	default:
		return fmt.Errorf("failed to determine node ID: %v", node)
	}

	_, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/refresh", repositoryId, branchId, doc.(string)), nil, MapToReader(node))
	return err
}

func (session *CloudCmsSession) ChangeNodeQName(repositoryId string, branchId string, nodeId string, newQName string) error {
	params := url.Values{"qname": []string{newQName}}
	_, err := session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/change_qname", repositoryId, branchId, nodeId), params, nil)
	return err
}

func (session *CloudCmsSession) NodeTree(repositoryId string, branchId string, nodeId string, config JsonObject) (JsonObject, error) {
	params := url.Values{}
	for key, val := range config {
		if valStr, ok := val.(string); ok {
			switch key {
			case "leafPath":
				params.Add("leaf", valStr)
			case "leaf":
				params.Add("leaf", valStr)
			case "basePath":
				params.Add("base", valStr)
			case "base":
				params.Add("base", valStr)
			}
		}

		if valBool, ok := val.(bool); ok {
			switch key {
			case "containers":
				params.Add("containers", strconv.FormatBool(valBool))
			case "properties":
				params.Add("properties", strconv.FormatBool(valBool))

			}
		}

		if valInt, ok := val.(int); ok {
			switch key {
			case "depth":
				params.Add("depth", strconv.Itoa(valInt))
			}
		}

	}

	payload := make(JsonObject)
	if val, ok := config["query"]; ok {
		payload["query"] = val
	}

	if val, ok := config["search"]; ok {
		payload["search"] = val
	}

	return session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/tree", repositoryId, branchId, nodeId), params, MapToReader(payload))
}

func (session *CloudCmsSession) ResolveNodePath(repositoryId string, branchId string, nodeId string) (string, error) {
	res, err := session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/path", repositoryId, branchId, nodeId), nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", res["path"]), nil
}

func (session *CloudCmsSession) ResolveNodePaths(repositoryId string, branchId string, nodeId string) (map[string]string, error) {
	res, err := session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/paths", repositoryId, branchId, nodeId), nil)
	if err != nil {
		return nil, err
	}

	pathsRaw := res.GetObject("paths")
	paths := make(map[string]string)
	for k, v := range pathsRaw {
		paths[k] = fmt.Sprintf("%v", v)
	}

	return paths, nil
}

func (session *CloudCmsSession) TraverseNode(repositoryId string, branchId string, nodeId string, config JsonObject) (JsonObject, error) {
	return session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/traverse", repositoryId, branchId, nodeId), nil, MapToReader(JsonObject{"traverse": config}))
}
func (session *CloudCmsSession) UploadAttachment(repositoryId string, branchId string, nodeId string, attachmentId string, file io.Reader, mimeType string, filename string) error {
	uri := fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/attachments/%s", repositoryId, branchId, nodeId, attachmentId)

	if attachmentId == "" {
		attachmentId = "default"
	}

	if filename == "" {
		filename = attachmentId
	}

	// Setup a multipart post request with one part corresponding to upload file
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, attachmentId, filename))
	partHeader.Set("Content-Type", mimeType)
	part, err := mw.CreatePart(partHeader)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	mw.Close()

	resp, err := session.MultipartPost(uri, nil, mw.FormDataContentType(), &b)
	if err == nil {
		defer resp.Close()
	}
	return err
}

func (session *CloudCmsSession) DownloadAttachment(repositoryId string, branchId string, nodeId string, attachmentId string) (io.ReadCloser, error) {
	return session.Download(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/attachments/%s", repositoryId, branchId, nodeId, attachmentId), nil)
}

func (session *CloudCmsSession) ListAttachments(repositoryId string, branchId string, nodeId string) (*ResultMap[JsonObject], error) {
	res, err := session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/attachments", repositoryId, branchId, nodeId), nil)
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) DeleteAttachment(repositoryId string, branchId string, nodeId string, attachmentId string) error {
	if attachmentId == "" {
		attachmentId = "default"
	}

	_, err := session.Delete(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/attachments/%s", repositoryId, branchId, nodeId, attachmentId), nil)
	return err
}

func (session *CloudCmsSession) ListVersions(repositoryId string, branchId string, nodeId string, options JsonObject, pagination JsonObject) (*ResultMap[JsonObject], error) {
	params := ToParams(options, pagination)

	res, err := session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/versions", repositoryId, branchId, nodeId), params)
	if err != nil {
		return nil, err
	}

	return ToResultMap(res), nil
}

func (session *CloudCmsSession) ReadVersion(repositoryId string, branchId string, nodeId string, changesetId string, options JsonObject) (JsonObject, error) {
	params := ToParams(options)
	return session.Get(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/versions/%s", repositoryId, branchId, nodeId, changesetId), params)
}

func (session *CloudCmsSession) RestoreVersion(repositoryId string, branchId string, nodeId string, changesetId string) (JsonObject, error) {
	return session.Post(fmt.Sprintf("/repositories/%s/branches/%s/nodes/%s/versions/%s/restore", repositoryId, branchId, nodeId, changesetId), nil, nil)
}
