package logs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	rotatelog "github.com/lestrrat/go-file-rotatelogs"
	"github.com/op/go-logging"

	"github.com/opensourceways/app-robot-server/config"
)

const (
	logDir      = "logs"
	logSoftLink = "latest_log"
	Module      = "app-robot"
)

var Logger = logging.MustGetLogger(Module)

func Init() error {
	var backends [] logging.Backend
	cLog := config.Application.Log

	if cLog.SaveFile {
		backend, err := registerFile(cLog)
		if err != nil {
			return err
		}
		backends = append(backends, backend)
	}
	backends = append(backends, registerStdout(cLog))
	logging.SetBackend(backends...)
	return nil
}

func registerFile(log config.Log) (logging.Backend, error) {
	if ok := pathExists(logDir); !ok {
		fmt.Println("create log directory")
		_ = os.Mkdir(logDir, os.ModePerm)
	}
	fileWriter, err := rotatelog.New(
		logDir+string(os.PathSeparator)+"%Y-%m-%d.log",
		rotatelog.WithLinkName(logSoftLink),
		rotatelog.WithMaxAge(7*24*time.Hour),
		rotatelog.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}
	level, err := logging.LogLevel(log.Level)
	if err != nil {
		return nil, err
	}
	return createBackend(fileWriter, log, level), nil
}

func registerStdout(log config.Log) logging.Backend {
	level, err := logging.LogLevel(log.Level)
	if err != nil {
		fmt.Println(err)
	}
	return createBackend(os.Stdout, log, level)
}

func createBackend(w io.Writer, log config.Log, level logging.Level) logging.Backend {
	backend := logging.NewLogBackend(w, log.Prefix, 0)
	stoutWriter := false
	if w == os.Stdout {
		stoutWriter = true
	}
	format := getLogFormatter(stoutWriter)
	backendLeveled := logging.AddModuleLevel(logging.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, Module)
	return backendLeveled
}

func getLogFormatter(stdoutWriter bool) logging.Formatter {
	pattern := `%{time:2006/01/02 - 15:04:05.000} %{shortfile} %{color}▶ [%{level:.6s}] %{color:reset}%{message}`
	if !stdoutWriter {
		pattern = strings.Replace(pattern, "%{color}", "", -1)
		pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	}
	return logging.MustStringFormatter(pattern)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
