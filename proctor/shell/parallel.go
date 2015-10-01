package shell

type runner interface {
	ConnectAndRun(host, command string, options *ConnectionOptions) (string, error)
}

type ParallelRunner struct {
	Runner runner
}

type Result struct {
	Host   string
	Stdout string
	Error  error
}

func (r *ParallelRunner) ConnectAndRun(hosts []string, command string, options *ConnectionOptions) map[string]Result {
	runner := func(host string, c chan Result) {
		stdout, err := r.Runner.ConnectAndRun(host, command, options)
		c <- Result{
			Host:   host,
			Stdout: stdout,
			Error:  err,
		}
	}

	resultsChannel := make(chan Result, len(hosts))
	for _, host := range hosts {
		go runner(host, resultsChannel)
	}

	results := map[string]Result{}
	for range hosts {
		result := <-resultsChannel
		results[result.Host] = result
	}

	return results
}
