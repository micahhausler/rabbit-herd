package herd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// A struct for storing the `name` parameter of the API
type Node struct {
	Name string `json:"name"`
}

// Get the node's hostname portion of the name.
func (n Node) Hostname() string {
	parts := strings.Split(n.Name, "@")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// Get the connected nodes from the API from "/api/nodes"
func GetApiHosts(api, user, pass string) []string {
	hosts := []string{}
	uri := fmt.Sprintf("%s/api/nodes", api)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return hosts
	}
	req.SetBasicAuth(user, pass)
	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return hosts
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return hosts
	}
	nodes := []Node{}
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return hosts
	}
	for _, node := range nodes {
		hosts = append(hosts, node.Hostname())
	}
	return hosts
}
