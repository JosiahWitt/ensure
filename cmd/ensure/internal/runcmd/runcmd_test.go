package runcmd_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/runcmd"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func TestRunnerExec(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("with valid command execution", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(&runcmd.ExecParams{
			PWD:  "/tmp",
			CMD:  "sh",
			Args: []string{"-c", "pwd"},
		})

		ensure(err).IsNotError()
		ensure(result).Equals("/tmp\n")
	})

	ensure.Run("with invalid command", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(&runcmd.ExecParams{
			CMD: "this-command-does-not-exist",
		})

		var expectedErr *exec.Error
		ensure(errors.As(err, &expectedErr)).IsTrue()
		ensure(result).IsEmpty()
	})

	ensure.Run("with failing command", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(&runcmd.ExecParams{
			CMD:  "sh",
			Args: []string{"-c", "echo 'abc'; exit 1"},
		})

		ensure(err.Error()).Equals("abc\n")
		ensure(result).IsEmpty()
	})
}
