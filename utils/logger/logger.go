package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set log level
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "DEBUG":
		Log.SetLevel(logrus.DebugLevel)
	case "INFO":
		Log.SetLevel(logrus.InfoLevel)
	case "WARN":
		Log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if os.Getenv("LOG_FORMAT") == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	Log.SetOutput(os.Stdout)
}

// Helper functions for structured logging
func LogAPIRequest(method, path string, userID string) {
	Log.WithFields(logrus.Fields{
		"type":    "api_request",
		"method":  method,
		"path":    path,
		"user_id": userID,
	}).Info("API request received")
}

func LogAPIResponse(method, path string, statusCode int, userID string) {
	Log.WithFields(logrus.Fields{
		"type":        "api_response",
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"user_id":     userID,
	}).Info("API response sent")
}

func LogError(err error, context string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["type"] = "error"
	fields["context"] = context
	Log.WithFields(fields).Error(err.Error())
}

func LogInfo(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["type"] = "info"
	Log.WithFields(fields).Info(message)
}

func LogDebug(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["type"] = "debug"
	Log.WithFields(fields).Debug(message)
}

func LogWarning(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["type"] = "warning"
	Log.WithFields(fields).Warn(message)
}
