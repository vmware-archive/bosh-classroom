package mocks

import "github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"

type ConnectAndRunCall struct {
	Receives struct {
		Host    string
		Command string
		Options *shell.ConnectionOptions
	}

	Returns struct {
		Stdout string
		Error  error
	}
}

type Runner struct {
	ConnectAndRunCalls     map[int]*ConnectAndRunCall
	ConnectAndRunCallCount int
}

func NewRunner(maxCallCount int) *Runner {
	calls := map[int]*ConnectAndRunCall{}
	for i := 0; i <= maxCallCount; i++ {
		calls[i] = &ConnectAndRunCall{}
	}
	return &Runner{ConnectAndRunCalls: calls}
}

func (r *Runner) ConnectAndRun(host, command string, options *shell.ConnectionOptions) (string, error) {
	call := r.ConnectAndRunCalls[r.ConnectAndRunCallCount]
	defer func() { r.ConnectAndRunCallCount++ }()

	call.Receives.Host = host
	call.Receives.Command = command
	call.Receives.Options = options
	return call.Returns.Stdout, call.Returns.Error
}
