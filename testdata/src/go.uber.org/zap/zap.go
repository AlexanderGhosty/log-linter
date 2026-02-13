package zap

type Logger struct{}

func NewExample() *Logger { return &Logger{} }

func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}
func (l *Logger) Debug(msg string, fields ...Field) {}

type Field struct{}

func String(key, val string) Field          { return Field{} }
func Int(key string, val int) Field         { return Field{} }
func Any(key string, val interface{}) Field { return Field{} }

func L() *Logger        { return &Logger{} }
func S() *SugaredLogger { return &SugaredLogger{} }

type SugaredLogger struct{}

func (s *SugaredLogger) Info(args ...interface{})                       {}
func (s *SugaredLogger) Infof(template string, args ...interface{})     {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{}) {}
