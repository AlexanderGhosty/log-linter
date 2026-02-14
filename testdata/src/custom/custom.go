package custom

type MyLogger struct{}

func (l *MyLogger) CustomInfo(msg string) {}

func DoLog() {
	l := &MyLogger{}
	l.CustomInfo("Starting server") // want "log message should start with a lowercase letter"
	l.CustomInfo("starting server") // OK
}
