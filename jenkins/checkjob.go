package jenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"time"
)

type CheckJobResult struct {
	Name string `json:"name"`
	Id int  `json:"id"`
	Status string   `json:"status"`
	StartTime int64 `json:starttime`
}

var (
	JobStatus string
	JobResult CheckJobResult
)

func CheckJobStatus(jksclient *gojenkins.Jenkins, jobName string, number int64) (CheckJobResult, error){

START:

	if buildResult, err := jksclient.GetBuild(jobName,number); err != nil {
		return CheckJobResult{}, err
	} else {

		JobResult = CheckJobResult{jobName,int(number), buildResult.Raw.Result,buildResult.Raw.Timestamp}
		v := buildResult.Raw.Building
		for {
			switch v {
			case false:
				fmt.Printf("job执行成功,buildID: %d\n", int(buildResult.Raw.Number))
				goto END
			case true:
				time.Sleep(5*time.Second)
				fmt.Printf("job正在构建中....%s\n", time.Now().Format("2006-01-02 15:04:05"))
				goto START
			}
		}
	}
END:
	return JobResult, nil

}
