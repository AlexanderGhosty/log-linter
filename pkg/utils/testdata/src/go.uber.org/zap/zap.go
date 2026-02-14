package zap

type Logger struct{}

func NewProduction() (*Logger, error) {
	return &Logger{}, nil
}

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{}
}

func (l *Logger) Info(msg string, fields ...Field) {}

type SugaredLogger struct{}

func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {}

type Field struct{}

func String(key string, val string) Field { return Field{} }
