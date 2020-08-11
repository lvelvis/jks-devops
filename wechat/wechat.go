package wechat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"gopkg.in/ini.v1"
)

type JSON struct {
	Access_token string `json:"access_token"`
}

//MESSAGES wx
type MESSAGES struct {
	Touser  string `json:"touser"`
	Toparty int    `json:"toparty"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    struct {
		//Subject string `json:"subject"`
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

//Config 参数
type Config struct {
	ToUser     string
	ToParty    int
	AgentID    int
	CorpID     string
	CorpSecret string
}

var WebChatCfg = &Config{}

var cfg *ini.File

//InitConfig 初始化加载企业微信配置
func InitConfig() {
	var err error
	cfg, err = ini.LooseLoad("/usr/local/jks-devops/wechat/wechat.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'wechat.ini': %v", err)
	}
	mapTo("wechat", WebChatCfg)

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo "+section+" err: %v", err)
	}
}

//GetAccessToken 获取企业微信token
func GetAccessToken(corpid, corpsecret string) string {
	gettoken_url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret

	client := &http.Client{}
	req, _ := client.Get(gettoken_url)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	//fmt.Printf("\n%q",string(body))
	var jsonstr JSON
	json.Unmarshal([]byte(body), &jsonstr)
	return jsonstr.Access_token
}

//SendMessage 企业微信接口请求
func SendMessage(accessToken, msg string) {
	sendurl := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken
	client := &http.Client{}
	req, _ := http.NewRequest("POST", sendurl, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func messages(touser string, toparty int, agentid int, content string) string {
	msg := MESSAGES{
		Touser:  touser,
		Toparty: toparty,
		Msgtype: "text",
		Agentid: agentid,
		Safe:    0,
		Text: struct {
			//Subject string `json:"subject"`
			Content string `json:"content"`
		}{Content: content},
	}
	sedmsg, _ := json.Marshal(msg)
	return string(sedmsg)
}

//SendWebChat 企业微信报警
func SendWebChat(jobContent, jobName, jobStatus, branch, commitid, logUrl string, startTime int64) {
	InitConfig()

	accessToken := GetAccessToken(WebChatCfg.CorpID, WebChatCfg.CorpSecret)

	FullContent := "[定时发布]" + ": " + jobContent + "\n" +
		"[发布服务]" + ": " + jobName + "\n" +
		"[构建状态]" + ": " + jobStatus + "\n" +
		"[构建分支]" + ": " + branch + "\n" +
		"[构建HashID]" + ": " + commitid + "\n" +
		"[构建日志]" + ": " + logUrl + "\n" +
		"[构建时间]" + ": " + time.Unix(startTime/1000, 0).Format("2006-01-02 15:04:05")

	msg := messages(WebChatCfg.ToUser, WebChatCfg.ToParty, WebChatCfg.AgentID, FullContent)
	SendMessage(accessToken, msg)
}
