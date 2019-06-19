package main

import (
	"flag"
	"fmt"
	"netserver"
	"netserver/log"
	"os"
)

var gLog = log.NewLogger()

func help() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func parseFlag() {

	logPath := flag.String("logPath", "./log", "log path")
	logName := flag.String("logName", "", "log name")
	logLv := flag.Int("logLv", 1, "log level")

	flag.Parse()

	var config = &netserver.GlobalConfig

	config.LogPath = *logPath
	config.LogName = *logName
	config.LogLevel = *logLv
}

func parseInIConf() {
	//fp,err := os.OpenFile("./netserver.ini",os.O_RDONLY,0644)
	//if err != nil {
	//
	//}
	//context := make([]byte,200)
	//_,err := fp.Read(context)

}

func initLog() {
	config := &netserver.GlobalConfig

	if len(config.LogName) == 0 {
		config.LogName = "netServer.log"
	}

	gLog.SetLogPath(config.LogPath, config.LogName)
	gLog.SetLogLevel(uint64(config.LogLevel))
}

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "help" {
		help()
		return
	}

	parseFlag()

	initLog()

	jobDispatcher = NewDispatcher()
	jobDispatcher.Start()
}
