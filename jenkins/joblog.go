package jenkins

import (
	"io"
	"log"
	"os"
)

func SaveConsoleText(jobConsoleText  string ) (logConsoleUrl string) {
	//创建日志文件
	logname := RandName() + ".txt"
	logConsoleUrl = "http://jks-devops.mobileztgame.com/" + logname
	f, err := os.OpenFile("/usr/local/jks-devops/logs/"+logname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//完成后，延迟关闭
	defer f.Close()
	// 设置日志输出到文件
	// 定义多个写入器
	writers := []io.Writer{
		f,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	logger := log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println(jobConsoleText)
	return
}
