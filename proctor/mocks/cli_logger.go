package mocks

import "fmt"

type CLILogger struct{}

func (l *CLILogger) Println(indentation int, format string, args ...interface{}) {

}

func (l *CLILogger) Green(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (l *CLILogger) Red(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
