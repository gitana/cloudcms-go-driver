# cloudcms-go-driver

HTTP Driver for the for the [Cloud CMS](https://www.cloudcms.com) API

## Installation

TODO

Below are some examples of how you might use this driver:

```go
package main

import (
    "cloudcms"
)

func main() {
    // Connect to CloudCMS using gitana.json in working directory
    session, err := ConnectDefault()
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
    nodeObj := JsonObject{
        "title": "Twelfth Night",
        "description": "An old play",
    }
    nodeId, _ := session.createNode(repositoryId, branchId, nodeObj, nil)

    // Query Nodes
    query := JsonObject{
        "_type": "store:book",
    }
    pagination := JsonObject{
        "limit": 1,
    }
    queriedNodes, _ session.QueryNodes(repositoryId, branchId, query, pagination)

    // Find Nodes
    find := JsonObject{
        "search": "Shakespeare",
        "query": JsonObject{
            "_type": "store:book"
        }
    }
    findNodes, _ := session.FindNodes(repositoryId, branchId, find ,nil)
}
```

## Resources

* Cloud CMS: https://www.cloudcms.com
* Github: https://github.com/gitana/cloudcms-go-driver
* Go Driver Download: TODO
* Cloud CMS Documentation: https://www.cloudcms.com/documentation.html
* Developers Guide: https://www.cloudcms.com/developers.html

## Support

For information or questions about the Go Driver, please contact Cloud CMS
at [support@cloudcms.com](mailto:support@cloudcms.com).