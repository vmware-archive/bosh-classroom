package shell

type runner interface {
	ConnectAndRun(host, command string) (string, error)
}

type ParallelRunner struct {
	Runner runner
}

type Result struct {
	Host   string
	Stdout string
	Error  error
}

func (r *ParallelRunner) ConnectAndRun(hosts []string, command string) map[string]Result {
	results := map[string]Result{}
	for _, h := range hosts {
		stdout, err := r.Runner.ConnectAndRun(h, command)
		results[h] = Result{
			Host:   h,
			Stdout: stdout,
			Error:  err,
		}
	}
	return results
}
