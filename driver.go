package cloudcms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"

	"golang.org/x/oauth2"
)

type CloudcmsConfig struct {
	Client_id     string `json:"clientKey"`
	Client_secret string `json:"clientSecret"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	BaseURL       string `json:"baseURL"`
	Debug         bool   `json:"debug"`
}

type CloudCmsSession struct {
	oauthClient *http.Client
	config      *CloudcmsConfig
}

type JsonObject map[string]interface{}

type ResultMap struct {
	rows       []JsonObject
	size       int
	total_rows int
	offset     int
}

func (obj *JsonObject) GetString(key string) string {
	val, ok := (*obj)[key]
	if !ok {
		return ""
	}

	return fmt.Sprintf("%v", val)
}

func (obj *JsonObject) GetObject(key string) JsonObject {
	val, ok := (*obj)[key]
	if !ok {
		return nil
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}

	return JsonObject(m)
}

func (obj *JsonObject) GetArray(key string) []interface{} {

	val, ok := (*obj)[key]
	if !ok {
		return nil
	}

	a, ok := val.([]interface{})
	if !ok {
		return nil
	}

	return a
}

func (obj *JsonObject) GetObjectArray(key string) []JsonObject {
	arr := obj.GetArray(key)
	if arr == nil {
		return nil
	}

	res := make([]JsonObject, len(arr))
	for i, v := range arr {
		res[i] = v.(map[string]interface{})
	}

	return res
}

// This type implements the http.RoundTripper interface
type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	info := fmt.Sprintf("%v %v %v", req.Method, req.URL, req.Proto)
	fmt.Println(info)

	// Send the request, get the response (or the error)
	res, e = lrt.Proxied.RoundTrip(req)

	// Handle the result.
	if e != nil {
		fmt.Printf("Error: %v", e)
		debug.PrintStack()
	} else {
		fmt.Printf("Received %v response\n", res.Status)
	}

	return res, e
}

func buildOAuthClient(cloudcmsConfig *CloudcmsConfig) (*http.Client, error) {
	ctx := context.Background()
	httpClient := http.Client{}

	if cloudcmsConfig.Debug {
		httpClient.Transport = LoggingRoundTripper{http.DefaultTransport}
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, &httpClient)
	conf := &oauth2.Config{
		ClientID:     cloudcmsConfig.Client_id,
		ClientSecret: cloudcmsConfig.Client_secret,
		Endpoint: oauth2.Endpoint{
			TokenURL:  cloudcmsConfig.BaseURL + "/oauth/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes: []string{"api"},
	}

	token, err := conf.PasswordCredentialsToken(ctx, cloudcmsConfig.Username, cloudcmsConfig.Password)
	if err != nil {
		return nil, err
	}

	oauthClient := conf.Client(ctx, token)
	return oauthClient, nil
}

func ToParams(objs ...JsonObject) url.Values {
	params := url.Values{}

	for _, obj := range objs {
		if obj != nil {
			for key, val := range obj {
				bytes, _ := json.Marshal(val)
				params.Add(key, string(bytes))
			}
		}
	}

	return params
}

func ToResultMap(obj JsonObject) *ResultMap {
	// rows := []JsonObject{}
	// rowsInterface := obj["rows"]

	// rowsArr, ok := rowsInterface.([]interface{})
	// if !ok {
	// 	panic("Failed to retrieve rows")
	// }

	// for _, rowObj := range rowsArr {
	// 	rows = append(rows, rowObj.(map[string]interface{}))
	// }

	rows := obj.GetObjectArray("rows")

	return &ResultMap{
		rows:       rows,
		size:       int(obj["size"].(float64)),
		total_rows: int(obj["total_rows"].(float64)),
		offset:     int(obj["offset"].(float64)),
	}
}

func MapToReader(obj JsonObject) io.Reader {
	if obj == nil {
		return bytes.NewReader([]byte("{}"))
	}

	data, _ := json.Marshal(obj)
	return bytes.NewReader(data)
}

func LoadConfig() *CloudcmsConfig {
	var result *CloudcmsConfig
	result = nil
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}

	if result == nil {
		result = ReadConfig(filepath.Join(wd, "gitana.json"))
	}
	if result == nil {
		result = ReadConfig(filepath.Join(wd, "gitana-test.json"))
	}
	if result == nil {
		result = ReadConfig(filepath.Join(wd, "cloudcms.json"))
	}
	if result == nil {
		result = ReadConfig(filepath.Join(wd, "cloudcms-test.json"))
	}

	return result
}

func ReadConfig(path string) *CloudcmsConfig {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var res CloudcmsConfig
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil
	}

	return &res
}

func ConnectDefault() (*CloudCmsSession, error) {
	config := LoadConfig()
	if config == nil {
		return nil, fmt.Errorf("could not locate gitana.json")
	}

	return Connect(config)
}

func Connect(cloudcmsConfig *CloudcmsConfig) (*CloudCmsSession, error) {
	oauthClient, err := buildOAuthClient(cloudcmsConfig)
	if err != nil {
		return nil, err
	}

	client := &CloudCmsSession{
		oauthClient: oauthClient,
		config:      cloudcmsConfig,
	}

	return client, nil
}

func (session *CloudCmsSession) Request(req *http.Request) (*http.Response, error) {
	resp, err := session.oauthClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Throw error for non-2xx repsonses
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Try to get response message
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, b)
	}

	return resp, nil
}

func (session *CloudCmsSession) RequestJson(method string, uri string, params url.Values, body io.Reader) (JsonObject, error) {
	if params == nil {
		params = url.Values{}
	}

	if !params.Has("full") {
		params.Add("full", "true")
	}
	if !params.Has("metadata") {
		params.Add("metadata", "true")
	}

	uri += "?" + params.Encode()

	req, err := http.NewRequest(method, session.config.BaseURL+uri, body)
	if err != nil {
		return nil, err
	}

	resp, err := session.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	target := make(JsonObject)
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (session *CloudCmsSession) Get(url string, params url.Values) (JsonObject, error) {
	return session.RequestJson("GET", url, params, nil)
}

func (session *CloudCmsSession) Delete(url string, params url.Values) (JsonObject, error) {
	return session.RequestJson("DELETE", url, params, nil)
}

func (session *CloudCmsSession) Post(url string, params url.Values, body io.Reader) (JsonObject, error) {
	return session.RequestJson("POST", url, params, body)
}

func (session *CloudCmsSession) Put(url string, params url.Values, body io.Reader) (JsonObject, error) {
	return session.RequestJson("PUT", url, params, body)
}

func (session *CloudCmsSession) Patch(url string, params url.Values, body io.Reader) (JsonObject, error) {
	return session.RequestJson("PATCH", url, params, body)
}

func (session *CloudCmsSession) Download(url string, params url.Values) (io.ReadCloser, error) {
	if params != nil {
		url += "?" + params.Encode()
	}

	req, _ := http.NewRequest("GET", session.config.BaseURL+url, nil)
	resp, err := session.Request(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (session *CloudCmsSession) MultipartPost(url string, params url.Values, formContentType string, body io.Reader) (io.ReadCloser, error) {
	if params != nil {
		url += "?" + params.Encode()
	}

	req, _ := http.NewRequest("POST", session.config.BaseURL+url, body)
	req.Header.Set("Content-Type", formContentType)

	resp, err := session.Request(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
