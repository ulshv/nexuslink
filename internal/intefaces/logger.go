package intefaces

type Logger interface {
	// Log prints messages as is (i.e. <message> arg1 arg2 ...)
	Log(message string, args ...any)
	// Info prints messages with the following prefix: "<timestamp> INFO <message>"
	Info(message string, args ...any)
	// Error prints messages in red color with the following prefix: "<timestamp> ERROR <message>"
	Error(message string, args ...any)
}
