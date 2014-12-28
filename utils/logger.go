package utils

// Modeled after https://github.com/Sirupsen/logrus
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	// Should call os.Exit(1) after logging
	Fatalf(format string, args ...interface{})
	// Should call panic() after logging
	Panicf(format string, args ...interface{})
}
