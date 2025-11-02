package cmd

import (
	"fmt"

	"github.com/mesh-dell/github-activity/internal/activity"
)

func Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("error: please provide a github username")
	}

	userName := args[1]
	events, err := activity.GetUserActivity(userName)

	if err != nil {
		return err
	}

	return activity.PrintUserActivity(userName, events)
}
