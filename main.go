package main

import (
	"flag"
	"fmt"
	"github.com/bndr/gojenkins"
	"jks-devops/jenkins"
	"jks-devops/wechat"
	"os"
)

type Auth struct {
	Username string
	ApiToken string
	BaseUrl  string
}

var (
	help     bool
	action   string
	app      string
	branch   string
	commitid string
	jobCurID int
	jobStatus string
	jobLogUrl string
	jobStartTime int64
)

func init() {
	flag.BoolVar(&help, "help", false, "this help")
	flag.StringVar(&action, "action", "roll-update", "灰度发布:roll-update或者版本回退:roll-back")
	flag.StringVar(&app, "app", "", "版本发布的jenkins的job名称")
	flag.StringVar(&branch, "branch", "", "发布应用版本的分支或者tag")
	flag.StringVar(&commitid, "commitid", "", "发布应用版本的commitid")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `jenkins-x 版本发布状态检测 
Usage: jks-devops -action="roll-update" -app="testjob" -branch="v0.0.1" -commitid="003ba51153266329e2b207f8824743876b53259c"

Options:
`)
	flag.PrintDefaults()

}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	auth := &Auth{
		Username: os.Getenv("JENKINS_API_USER"),
		ApiToken: os.Getenv("JENKINS_API_TOKEN"),
		BaseUrl:  os.Getenv("JENKINS_HOST"),
	}
	JksClient := gojenkins.CreateJenkins(nil, auth.BaseUrl, auth.Username, auth.ApiToken)

	if _, err := JksClient.Init(); err != nil {
		panic("Jenkins init went wrong")
	}

	para := map[string]string{
		"Action":   action,
		"branch":   branch,
		"CommitID": commitid,
	}

	if action != "" || branch != "" || commitid != "" {
		//var wg sync.WaitGroup

		queueID, _ := JksClient.BuildJob(app, para)

		for {
			CurJob, err := jenkins.GetJobID(queueID)
			if err != nil {
				panic(err)
			}

			jobCurID = int(CurJob.Executable.Number)

			//jobLogUrl = CurJob.Executable.URL
			switch jobCurID {
			case 0:
				continue
			default:
				JobCheckInfo, err := jenkins.CheckJobStatus(JksClient,app,int64(jobCurID))
				if err !=nil {
					panic(err)
				}
				jobStatus = JobCheckInfo.Status
				jobStartTime = JobCheckInfo.StartTime
				goto END
			}

		}

	} else {
		fmt.Println("输入参数不正确")
		os.Exit(1)
	}

END:
	buildResult, err := JksClient.GetBuild(app, int64(jobCurID))
	if err != nil {
		panic(err)
	}
	jobLogUrl = jenkins.SaveConsoleText(buildResult.GetConsoleOutput())
	wechat.SendWebChat("自动化发布", app, jobStatus, branch, commitid, jobLogUrl, jobStartTime)
}
