# cloudcms-go-driver

HTTP Driver for the for the [Cloud CMS](https://gitana.io) API

## Installation

In your module directory, run:

`go get github.com/gitana/cloudcms-go-driver`
## Examples

Below are some examples of how you might use this driver:

```go
package main

import (
	"github.com/gitana/cloudcms-go-driver"
)

func main() {
    // Connect to CloudCMS using gitana.json in working directory
    session, err := cloudcms.ConnectDefault()
	if err != nil {
		fmt.Println(err)
        return
	}

    var repositoryId string


    // List branches
    branches, _ := session.ListBranches(repositoryId, nil)

    // Read branch
    branchId := "master"
    branch, _ := session.ReadBranch(repositoryId, branchId)

    // Read Node
    node, _ := session.ReadNode(repositoryId, branchId, nodeId)

    // Create Node
    nodeObj := cloudcms.JsonObject{
        "title": "Twelfth Night",
        "description": "An old play",
    }
    nodeId, _ := session.createNode(repositoryId, branchId, nodeObj, nil)

    // Query Nodes
    query := cloudcms.JsonObject{
        "_type": "store:book",
    }
    pagination := cloudcms.JsonObject{
        "limit": 1,
    }
    queriedNodes, _ session.QueryNodes(repositoryId, branchId, query, pagination)

    // Find Nodes
    find := cloudcms.JsonObject{
        "search": "Shakespeare",
        "query": JsonObject{
            "_type": "store:book",
        }
    }
    findNodes, _ := session.FindNodes(repositoryId, branchId, find ,nil)
}
```

## Resources

* Cloud CMS: https://gitana.io
* Github: https://github.com/gitana/cloudcms-go-driver
* Go Driver Download: TODO
* Cloud CMS Documentation: https://gitana.io/documentation.html
* Developers Guide: https://gitana.io/developers.html

## Support

For information or questions about the Go Driver, please contact Cloud CMS
at [support@cloudcms.com](mailto:support@cloudcms.com).
