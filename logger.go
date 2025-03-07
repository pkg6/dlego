package dlego

import (
	"fmt"
	"strings"
)

type ILogger interface {
	LogF(tag string, format string, args ...any)
	Gets() []string
	Flush()
}
type DefaultLogger struct {
	records []string
}

var _ ILogger = (*DefaultLogger)(nil)

func (l *DefaultLogger) LogF(tag string, format string, args ...any) {
	temp := make([]string, len(args)+1)
	temp[0] = tag
	temp[1] = fmt.Sprintf(format, args...)
	l.records = append(l.records, strings.Join(temp, ": "))
}
func (l *DefaultLogger) Gets() []string {
	temp := make([]string, len(l.records))
	copy(temp, l.records)
	return temp
}
func (l *DefaultLogger) Flush() {
	l.records = make([]string, 0)
}
