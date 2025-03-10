package dlego

import "fmt"

type ILogger interface {
	GetProviderName() string
	LogF(format string, args ...any)
	Gets() []string
	Flush()
}

func NewLoggerWithProviderName(providerName string) *DefaultLogger {
	return &DefaultLogger{
		records:      make([]string, 0),
		providerName: providerName,
	}
}

type DefaultLogger struct {
	providerName string
	records      []string
}

var _ ILogger = (*DefaultLogger)(nil)

func (l *DefaultLogger) GetProviderName() string {
	return l.providerName
}

func (l *DefaultLogger) LogF(format string, args ...any) {
	l.records = append(l.records, fmt.Sprintf(format, args...))
}
func (l *DefaultLogger) Gets() []string {
	temp := make([]string, len(l.records))
	copy(temp, l.records)
	return temp
}
func (l *DefaultLogger) Flush() {
	l.records = make([]string, 0)
	l.providerName = ""
}
