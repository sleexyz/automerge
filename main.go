package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var version = "1.0.0"

type CheckRun struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HTMLURL    string `json:"html_url"`
}

type CheckSuite struct {
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

type StatusCheck struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
	TargetURL   string `json:"target_url"`
}

type ChecksResponse struct {
	CheckRuns   []CheckRun   `json:"check_runs"`
	CheckSuites []CheckSuite `json:"check_suites"`
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("automerge version %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("Usage: automerge")
		fmt.Println("Polls GitHub API for status checks on current branch.")
		fmt.Println("Exits 0 if all checks pass, exits 1 if any fail.")
		os.Exit(0)
	}

	// Check if we're in a git repo
	if !isGitRepo() {
		fmt.Fprintf(os.Stderr, "Error: not in a git repository\n")
		os.Exit(1)
	}

	// Check if gh is available
	if !isGHAvailable() {
		fmt.Fprintf(os.Stderr, "Error: gh command not found. Please install GitHub CLI\n")
		os.Exit(1)
	}

	// Get current branch
	branch, err := getCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
		os.Exit(1)
	}

	// Get repo info
	owner, repo, err := getRepoInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting repository info: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Checking status for %s/%s on branch %s...\n", owner, repo, branch)

	// Poll for status checks
	for {
		allPassed, anyFailed, failureMessages := checkStatus(owner, repo, branch)
		
		if anyFailed {
			fmt.Fprintf(os.Stderr, "\nStatus checks failed:\n")
			for _, msg := range failureMessages {
				fmt.Fprintf(os.Stderr, "  ❌ %s\n", msg)
			}
			os.Exit(1)
		}
		
		if allPassed {
			fmt.Println("\n✅ All status checks passed!")
			os.Exit(0)
		}
		
		fmt.Print(".")
		time.Sleep(5 * time.Second)
	}
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func isGHAvailable() bool {
	cmd := exec.Command("gh", "--version")
	return cmd.Run() == nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getRepoInfo() (string, string, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "owner,name")
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	var repo struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(output, &repo); err != nil {
		return "", "", err
	}

	return repo.Owner.Login, repo.Name, nil
}

func checkStatus(owner, repo, branch string) (allPassed bool, anyFailed bool, failureMessages []string) {
	// Get the latest commit SHA for the branch
	cmd := exec.Command("gh", "api", fmt.Sprintf("/repos/%s/%s/branches/%s", owner, repo, branch))
	output, err := cmd.Output()
	if err != nil {
		return false, true, []string{fmt.Sprintf("Failed to get branch info: %v", err)}
	}

	var branchInfo struct {
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}

	if err := json.Unmarshal(output, &branchInfo); err != nil {
		return false, true, []string{fmt.Sprintf("Failed to parse branch info: %v", err)}
	}

	sha := branchInfo.Commit.SHA

	// Check status checks (older style)
	statusPassed, statusFailed, statusMessages := checkStatusChecks(owner, repo, sha)
	
	// Check check runs (newer style)
	checksPassed, checksFailed, checksMessages := checkCheckRuns(owner, repo, sha)

	// Combine results
	allPassed = statusPassed && checksPassed
	anyFailed = statusFailed || checksFailed
	failureMessages = append(statusMessages, checksMessages...)

	return allPassed, anyFailed, failureMessages
}

func checkStatusChecks(owner, repo, sha string) (bool, bool, []string) {
	cmd := exec.Command("gh", "api", fmt.Sprintf("/repos/%s/%s/commits/%s/status", owner, repo, sha))
	output, err := cmd.Output()
	if err != nil {
		return false, true, []string{fmt.Sprintf("Failed to get status checks: %v", err)}
	}

	var status struct {
		State    string        `json:"state"`
		Statuses []StatusCheck `json:"statuses"`
	}

	if err := json.Unmarshal(output, &status); err != nil {
		return false, true, []string{fmt.Sprintf("Failed to parse status checks: %v", err)}
	}

	if len(status.Statuses) == 0 {
		return true, false, nil // No status checks means we're good for this part
	}

	var failureMessages []string
	anyFailed := false

	for _, check := range status.Statuses {
		if check.State == "failure" || check.State == "error" {
			anyFailed = true
			failureMessages = append(failureMessages, fmt.Sprintf("%s: %s", check.Context, check.Description))
		}
	}

	allPassed := status.State == "success"
	return allPassed, anyFailed, failureMessages
}

func checkCheckRuns(owner, repo, sha string) (bool, bool, []string) {
	cmd := exec.Command("gh", "api", fmt.Sprintf("/repos/%s/%s/commits/%s/check-runs", owner, repo, sha))
	output, err := cmd.Output()
	if err != nil {
		return false, true, []string{fmt.Sprintf("Failed to get check runs: %v", err)}
	}

	var response ChecksResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return false, true, []string{fmt.Sprintf("Failed to parse check runs: %v", err)}
	}

	if len(response.CheckRuns) == 0 {
		return true, false, nil // No check runs means we're good for this part
	}

	var failureMessages []string
	anyFailed := false
	allCompleted := true

	for _, check := range response.CheckRuns {
		if check.Status != "completed" {
			allCompleted = false
			continue
		}
		
		if check.Conclusion == "failure" || check.Conclusion == "cancelled" || check.Conclusion == "timed_out" {
			anyFailed = true
			failureMessages = append(failureMessages, fmt.Sprintf("%s: %s", check.Name, check.Conclusion))
		}
	}

	allPassed := allCompleted && !anyFailed
	return allPassed, anyFailed, failureMessages
}