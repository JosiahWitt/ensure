package runcmd

import (
	"errors"
	"os/exec"
)

type ExecParams struct {
	PWD  string
	CMD  string
	Args []string
}

type RunnerIface interface {
	Exec(*ExecParams) (string, error)
}

type Runner struct{}

var _ RunnerIface = &Runner{}

// Exec the command defined in the provided params.
func (*Runner) Exec(params *ExecParams) (string, error) {
	//nolint:gosec
	c := exec.Command(params.CMD, params.Args...)
	c.Dir = params.PWD
	out, err := c.CombinedOutput()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			//nolint:goerr113
			return "", errors.New(string(out))
		}

		return "", err
	}

	return string(out), err
}
