package tui

type LogType int8

const (
	LogInfo LogType = iota
	LogError
)

type LogMsg struct {
	Text string
	Err  error
	Type LogType
}

type UILogger struct{}

func (UILogger) Logf(format string, args ...any) {
	Logf(format, args...)
}

func (UILogger) Errorf(err error, format string, args ...any) {
	Errorf(err, format, args...)
}
