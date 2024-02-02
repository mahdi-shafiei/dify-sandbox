package service

import (
	"time"

	"github.com/langgenius/dify-sandbox/internal/core/runner/python"
	"github.com/langgenius/dify-sandbox/internal/static"
	"github.com/langgenius/dify-sandbox/internal/types"
)

type RunCodeResponse struct {
	Stderr string `json:"error"`
	Stdout string `json:"stdout"`
}

func RunPython3Code(code string) *types.DifySandboxResponse {
	runner := python.PythonRunner{}
	stdout, stderr, done, err := runner.Run(code, time.Duration(static.GetCoshubGlobalConfigurations().WorkerTimeout*int(time.Second)), nil)
	if err != nil {
		return types.ErrorResponse(-500, err.Error())
	}

	stdout_str := ""
	stderr_str := ""

	for {
		select {
		case <-done:
			return types.SuccessResponse(&RunCodeResponse{
				Stdout: stdout_str,
				Stderr: stderr_str,
			})
		case out := <-stdout:
			stdout_str += string(out)
		case err := <-stderr:
			stderr_str += string(err)
		}
	}
}
