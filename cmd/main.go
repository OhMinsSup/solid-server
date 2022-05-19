package main

import (
	"flag"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"log"
	"os"
	"os/signal"
	"solid-server/server"
	config2 "solid-server/services/config"
	"syscall"
	"time"
)

// 공유 코드(dll)와 함께 사용되는 활성 서버
var pServer *server.Server

const (
	timeBetweenPidMonitoringChecks = 2 * time.Second
)

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))

	return err == nil
}

// monitorPid 는 다른(클라이언트 앱) 프로세스와 동기화된 서버 수명을 유지하는 데 사용됩니다.
func monitorPid(pid int, logger *mlog.Logger) {
	logger.Info("Monitoring PID", mlog.Int("pid", pid))

	go func() {
		for {
			if !isProcessRunning(pid) {
				logger.Info("Monitored process not found, exiting.")
				os.Exit(1)
			}

			time.Sleep(timeBetweenPidMonitoringChecks)
		}
	}()
}

func logInfo(logger *mlog.Logger) {
	logger.Info("Solid Server",
		mlog.String("version", model.CurrentVersion),
		mlog.String("build_number", model.BuildNumber),
		mlog.String("build_date", model.BuildDate),
		mlog.String("build_hash", model.BuildHash),
	)
}

func main() {
	pMonitorPid := flag.Int("monitorpid", -1, "a process ID")
	pPort := flag.Int("port", 0, "the port number")
	pDBType := flag.String("dbtype", "", "Database type")
	pDBConfig := flag.String("dbconfig", "", "Database config")
	pConfigFilePath := flag.String(
		"config",
		"",
		"Location of the JSON config file",
	)
	flag.Parse()

	config, err := config2.ReadConfigFile(*pConfigFilePath)
	if err != nil {
		log.Fatal("Unable to read the config file: ", err)
		return
	}

	logger, _ := mlog.NewLogger()
	cfgJSON := config.LoggingCfgJSON
	if config.LoggingCfgFile == "" && cfgJSON == "" {
		// if no logging defined, use default config (console output)
		cfgJSON = defaultLoggingConfig()
	}
	err = logger.Configure(config.LoggingCfgFile, cfgJSON, nil)
	if err != nil {
		log.Fatal("Error in config file for logger: ", err)
		return
	}
	defer func() { _ = logger.Shutdown() }()

	if logger.HasTargets() {
		restore := logger.RedirectStdLog(mlog.LvlInfo, mlog.String("src", "stdlog"))
		defer restore()
	}

	logInfo(logger)

	if pMonitorPid != nil && *pMonitorPid > 0 {
		monitorPid(*pMonitorPid, logger)
	}

	if pDBType != nil && len(*pDBType) > 0 {
		config.DBType = *pDBType
		logger.Info("DBType from commandline", mlog.String("DBType", *pDBType))
	}

	if pDBConfig != nil && len(*pDBConfig) > 0 {
		config.DBConfigString = *pDBConfig
		// Don't echo, as the confix string may contain passwords
		logger.Info("DBConfigString overriden from commandline")
	}

	if pPort != nil && *pPort > 0 && *pPort != config.Port {
		// Override port
		logger.Info("Port from commandline", mlog.Int("port", *pPort))
		config.Port = *pPort
	}

	db, err := server.NewStore(config, logger)
	if err != nil {
		logger.Fatal("server.NewStore ERROR", mlog.Err(err))
	}

	// permission services

	params := server.Params{
		Cfg:                config,
		DBStore:            db,
		Logger:             logger,
	}

	server, err := server.New(params)
	if err != nil {
		logger.Fatal("server.New ERROR", mlog.Err(err))
	}

	if err := server.Start(); err != nil {
		logger.Fatal("server.Start ERROR", mlog.Err(err))
	}

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	_ = server.Shutdown()
}

func defaultLoggingConfig() string {
	return `
	{
		"def": {
			"type": "console",
			"options": {
				"out": "stdout"
			},
			"format": "plain",
			"format_options": {
				"delim": " ",
				"min_level_len": 5,
				"min_msg_len": 40,
				"enable_color": true,
				"enable_caller": true
			},
			"levels": [
				{"id": 5, "name": "debug"},
				{"id": 4, "name": "info", "color": 36},
				{"id": 3, "name": "warn"},
				{"id": 2, "name": "error", "color": 31},
				{"id": 1, "name": "fatal", "stacktrace": true},
				{"id": 0, "name": "panic", "stacktrace": true}
			]
		}
	}`
}
