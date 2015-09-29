package commands

import "github.com/onsi/say"

type CliLogger struct{}

func (l *CliLogger) Println(i int, f string, a ...interface{}) {
	say.Println(i, f, a...)
}
func (l *CliLogger) Green(f string, a ...interface{}) string {
	return say.Green(f, a...)
}
func (l *CliLogger) Red(f string, a ...interface{}) string {
	return say.Red(f, a...)
}
