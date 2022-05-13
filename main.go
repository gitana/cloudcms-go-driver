package cloudcms

import (
	"fmt"
)

func main() {

	session, err := ConnectDefault()
	if err != nil {
		fmt.Println(err)
	}

	return
	cloudcmsConfig := &CloudcmsConfig{
		BaseURL:       "http://localhost:8080",
		Client_id:     "4bc95754-575e-4067-901b-a4c234e4a0e8",
		Client_secret: "SD6mE3zqnR2iec+51HoLdYWBgKp8BGTfaU1IIKLhV1lmrw+EpB9cj69iyr13r28Uci/j7A84F2c2mjPtNNX0sFvkLoIlDc7TP8OWELQ3WFg=",
		Username:      "bob",
		Password:      "password1",
	}

	session, err = Connect(cloudcmsConfig)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}

	platform, _ := session.ReadPlatform()

	fmt.Println(platform)
	fmt.Println(platform["_doc"])

	repositoryId := "7330339b45c50c4f53d3"
	nodeId := "35218c5222a72fbb5e73"

	// file, _ := os.Create("./test.jpeg")
	// att, _ := session.DownloadAttachment(repositoryId, "master", nodeId, "default")
	// defer att.Close()

	// io.Copy(file, att)
	// file.Close()

	node, _ := session.ReadNode(repositoryId, "master", nodeId)
	fmt.Println(node)
}
