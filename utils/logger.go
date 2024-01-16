package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	BuildTypeNone    = ""
	BuildTypeDev     = "dev"
	BuildTypeRelease = "release"

	logDirectory      = "logs"
	logFileNameFormat = "2006-01-02[15.04.05]"
)

func NewLogger(BuildType string) (*zap.SugaredLogger, error) {
	config := zap.Config{}

	if BuildType == BuildTypeRelease {
		if err := EnsureDirExists(logDirectory); err != nil {
			return nil, fmt.Errorf("ensure log directory exists: %w", err)
		}

		config = zap.NewProductionConfig()
		// TODO: Add rolling for the log file eg. https://stackoverflow.com/questions/45440491/how-to-configure-uber-go-zap-logger-for-rolling-filesystem-log
		currentTime := time.Now()
		formattedTime := currentTime.Format(logFileNameFormat)
		config.OutputPaths = []string{
			filepath.Join(logDirectory, fmt.Sprintf("%v.log", formattedTime)),
			// "stdout",
			"stderr",
		}

		config.Encoding = "console"

	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	}
	config.EncoderConfig.EncodeCaller = nil
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("create zap logger: %w", err)
	}

	sugar := logger.Sugar()

	return sugar, nil
}
