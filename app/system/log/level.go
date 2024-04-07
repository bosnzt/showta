package log

func Debugf(format string, v ...interface{}) {
	sugar.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	sugar.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	sugar.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	sugar.Errorf(format, v...)
}

func Debug(v ...interface{}) {
	sugar.Debug(v...)
}

func Info(v ...interface{}) {
	sugar.Info(v...)
}

func Warn(v ...interface{}) {
	sugar.Warn(v...)
}

func Error(v ...interface{}) {
	sugar.Error(v...)
}

func StdInfof(format string, v ...interface{}) {
	stdLogger.Infof(format, v...)
}

func StdErrorf(format string, v ...interface{}) {
	stdLogger.Errorf(format, v...)
}

func StdError(v ...interface{}) {
	stdLogger.Error(v...)
}
