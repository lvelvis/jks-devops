# jks-devops

jenkins-x自动化管理

## 调用方式
```
   Usage: jks-devops -action="roll-update" -app="testjob" -branch="v0.0.1" -commitid="003ba51153266329e2b207f8824743876b53259c"

Options:
  -action string
    	灰度发布:roll-update或者版本回退:roll-back (default "roll-update")
  -app string
    	版本发布的jenkins的job名称
  -branch string
    	发布应用版本的分支或者tag
  -commitid string
    	发布应用版本的commitid
  -help
    	this help
    	
   
```

环境变量注入:  
```
export JENKINS_API_USER=devops
export JENKINS_API_TOKEN=xxxxx  //jenkins user api token
export JENKINS_HOST=http://172.28.200.3
export JENKINS_WeChat_INI=/usr/local/jks-devops/wechat/wechat.ini  //企业微信配置

```
