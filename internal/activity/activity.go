package activity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repo struct {
	Name string `json:"name"`
}

type Event struct {
	Type    string `json:"type"`
	Repo    Repo   `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Head    string `json:"head"`
		Before  string `json:"before"`
	}
}

type CompareCommit struct {
	TotalCommits int `json:"total_commits"`
}

func PrintUserActivity(userName string, events []Event) error {
	if len(events) == 0 {
		return fmt.Errorf("no activity found")
	}

	fmt.Println()
	fmt.Printf("%s's recent activity\n", userName)
	fmt.Println()

	for _, event := range events {
		repoName := event.Repo.Name
		switch event.Type {
		case "CreateEvent":
			fmt.Println("- Created a new", event.Payload.RefType, repoName)
		case "PushEvent":
			commitCount, err := GetCommitCount(repoName, event.Payload.Before, event.Payload.Head)
			if err != nil {
				return err
			}
			fmt.Println("- Pushed ", commitCount, "commits to", repoName)
		case "WatchEvent":
			fmt.Println("- Starred repository", repoName)
		case "ForkEvent":
			fmt.Println("- Forked repository", repoName)
		case "PullRequestEvent":
			fmt.Println("-", event.Payload.Action, "pull request in", repoName)
		case "IssuesEvent":
			fmt.Println("-", event.Payload.Action, "an issue in", repoName)
		default:
			fmt.Println("-", event.Type, "in", repoName)
		}
	}
	return nil
}

func GetUserActivity(userName string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", userName)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("user not found. please check username")
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data: %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body%v", err)
	}

	var events []Event
	err = json.Unmarshal(body, &events)

	if err != nil {
		return nil, fmt.Errorf("error Decoding JSON")
	}

	return events, nil
}

func GetCommitCount(repoName, before, head string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/compare/%s...%s", repoName, before, head)
	res, err := http.Get(url)

	if err != nil {
		return 0, fmt.Errorf("error fetching commit count %s", err)
	}

	if res.StatusCode == 404 {
		return 0, fmt.Errorf("repository not found")
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching data: %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return 0, fmt.Errorf("error reading response body%v", err)
	}

	var compareCommits CompareCommit
	err = json.Unmarshal(body, &compareCommits)

	if err != nil {
		return 0, fmt.Errorf("error Decoding JSON: %s", err)
	}

	return compareCommits.TotalCommits, nil
}
