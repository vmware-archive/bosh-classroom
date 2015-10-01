package mocks

import (
	"fmt"
	"sync"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

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
	Calls                  chan string
	Unlocker               chan struct{}
	mutex                  sync.Mutex
	ConnectAndRunCalls     map[int]*ConnectAndRunCall
	ConnectAndRunCallCount int
}

func NewRunner(maxCallCount int) *Runner {
	calls := map[int]*ConnectAndRunCall{}
	for i := 0; i <= maxCallCount; i++ {
		calls[i] = &ConnectAndRunCall{}
	}
	return &Runner{
		ConnectAndRunCalls: calls,
		Calls:              make(chan string, maxCallCount),
		Unlocker:           make(chan struct{}, maxCallCount),
	}
}
func (r *Runner) ConnectAndRun(host, command string, options *shell.ConnectionOptions) (string, error) {
	// fmt.Println(host + " connection starting...")
	r.Calls <- host
	// fmt.Println(host + " called")
	<-r.Unlocker
	// fmt.Println(host + " unlocked")

	r.mutex.Lock()
	call := r.ConnectAndRunCalls[r.ConnectAndRunCallCount]
	r.ConnectAndRunCallCount++
	r.mutex.Unlock()

	call.Receives.Host = host
	call.Receives.Command = command
	call.Receives.Options = options
	return fmt.Sprintf("some result from host %s", host), fmt.Errorf("some error from host %s", host)
}

type ParallelRunner struct {
	ConnectAndRunCall struct {
		Receives struct {
			Hosts   []string
			Command string
			Options *shell.ConnectionOptions
		}

		Returns struct {
			Results map[string]shell.Result
		}
	}
}

func (r *ParallelRunner) ConnectAndRun(hosts []string, command string, options *shell.ConnectionOptions) map[string]shell.Result {
	r.ConnectAndRunCall.Receives.Hosts = hosts
	r.ConnectAndRunCall.Receives.Command = command
	r.ConnectAndRunCall.Receives.Options = options
	return r.ConnectAndRunCall.Returns.Results
}
