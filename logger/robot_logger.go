package logger

import "github.com/sirupsen/logrus"

var RobotLog *logrus.Logger

func InitRobotLogger(level string) {
	RobotLog = logrus.New()
	switch level {
	case "debug":
		RobotLog.Level = logrus.DebugLevel
	case "info":
		RobotLog.Level = logrus.InfoLevel
	case "warn":
		RobotLog.Level = logrus.WarnLevel
	case "error":
		RobotLog.Level = logrus.ErrorLevel
	case "fatal":
		RobotLog.Level = logrus.FatalLevel
	case "panic":
		RobotLog.Level = logrus.PanicLevel
	default:
		RobotLog.Level = logrus.DebugLevel
	}
}
