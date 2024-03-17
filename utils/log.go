package utils

import (
	"log/slog"
	"os"
	"strings"
)

var myLog *slog.Logger

func init() {
	Dev := "release"
	//Dev := "dev"
	writer := os.Stdout
	if Dev == "release" {
		file, err := os.OpenFile("./error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		writer = file
		if err != nil {
			panic("log file open err:" + err.Error())
		}
	}
	myLog = slog.New(slog.NewTextHandler(writer, nil))
}

func LogError(errMsg string, err error) {
	builder := strings.Builder{}
	builder.WriteString(errMsg)
	builder.WriteString(" err:")
	builder.WriteString(err.Error())
	myLog.Error(builder.String())
}
