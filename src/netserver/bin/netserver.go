package main

import (
	"data"
	"flag"
	"fmt"
	_ "module/dummy"
	"netserver"
	"netserver/log"
	"os"
	"os/signal"
	"syscall"
)

var gLog = log.NewLogger()

func help() {
	fmt.Println("Usage:")
	flag.PrintDefaults()
}

func parseFlag() {

	wanIp := flag.String("wanIp", "127.0.0.1", "wanIp")
	wanPort := flag.Int("wanPort", 8877, "wanPort")
	logPath := flag.String("logPath", "./log", "log path")
	logName := flag.String("logName", "netserver", "log name")
	logLv := flag.Int("logLv", 1, "log level")

	flag.Parse()

	var config = &netserver.GlobalConfig

	config.LogPath = *logPath
	config.LogName = *logName
	config.LogLevel = *logLv
	config.WanIp = *wanIp
	config.WanPort = *wanPort
	config.DBSource = "game:ljy1314@tcp(106.12.16.96:3306)/game001?charset=utf8"

	config.MaxClientReq = 20
	config.FrontPingMs = 1000
	config.UserIdleTimeMs = 3600000
	config.CommonPackageLen = 5120
	config.MaxPackageLen = 5120
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

	serverAddr := fmt.Sprintf("%s:%d", netserver.GlobalConfig.WanIp, netserver.GlobalConfig.WanPort)
	server := netserver.NewServer(serverAddr)

	err := data.InitDb(netserver.GlobalConfig.DBSource, netserver.GlobalConfig.LogPath)
	if err != nil {
		gLog.Fatal("init dbsrc:%s error : %s", netserver.GlobalConfig.DBSource, err)
		return
	}

	netserver.InitWorker(100)

	gLog.Info("Init Worker finish %d", 100)
	server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	s := <-quit
	fmt.Println("signal:", s)
	if s == syscall.SIGTERM || s == os.Interrupt {
		fmt.Println("stoping server")
	}

	onExit()
}

func onExit() {
	gLog.Flush()
}
