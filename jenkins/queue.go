package jenkins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)


type parameter struct {
	Name  string
	Value string
}

type generalAction struct {
	Causes     []map[string]interface{}
	Parameters []parameter
}

type taskResponse struct {
	Actions                    []generalAction `json:"actions"`
	Blocked                    bool            `json:"blocked"`
	Buildable                  bool            `json:"buildable"`
	BuildableStartMilliseconds int64           `json:"buildableStartMilliseconds"`
	ID                         int64           `json:"id"`
	InQueueSince               int64           `json:"inQueueSince"`
	Params                     string          `json:"params"`
	Pending                    bool            `json:"pending"`
	Stuck                      bool            `json:"stuck"`
	Task                       struct {
		Color string `json:"color"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"task"`
	URL        string `json:"url"`
	Why        string `json:"why"`
	Executable struct {
		Number int64  `json:"number"`
		URL    string `json:"url"`
	} `json:"executable"`
}
//通过job的队列ID获取构建的number
func GetJobID(id int64) (*taskResponse, error) {
	queueURL := os.Getenv("JENKINS_HOST") + "/queue/item/" + strconv.FormatInt(id, 10) + "/api/json"

	client := &http.Client{}
	req, err := http.NewRequest("get",queueURL,nil)
	req.SetBasicAuth(os.Getenv("JENKINS_API_USER"), os.Getenv("JENKINS_API_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("query job id  is failed:%s", resp.Status)
	}
	var result taskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
	
}
