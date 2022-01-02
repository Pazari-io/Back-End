package utils

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

// func WaterMarkPDF(file string,permissions bool) {
// return nill
// }

func ExecuteCommand(command string, timeOut time.Duration, args ...string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeOut*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, command, args...)

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		//log.Println("Command timed out")
		return "", errors.New("command timed out")
	}

	if err != nil {

		return "", err

	}

	//log.Println(string(out))

	// If there's no context error, we know the command completed (or errored).
	return string(out), nil

}
