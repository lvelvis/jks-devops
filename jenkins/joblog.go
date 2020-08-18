package jenkins

import (
	"io"
	"log"
	"os"
)
//SaveConsoleText 创建日志文件
func SaveConsoleText(jobConsoleText  string ) (logConsoleUrl string) {
	logname := RandName() + ".txt"
	logConsoleUrl = "http://jks-devops.mobileztgame.com/" + logname

	f, err := os.OpenFile("/usr/local/jks-devops/logs/"+logname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	writers := []io.Writer{
		f,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	logger := log.New(fileAndStdoutWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println(jobConsoleText)
	return logConsoleUrl
}
