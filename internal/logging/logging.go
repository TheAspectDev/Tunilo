package logging

type Logger interface {
	Logf(format string, args ...any)
	Errorf(err error, format string, args ...any)
}
