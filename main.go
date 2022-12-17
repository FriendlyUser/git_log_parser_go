package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GitCommit represents a commit in a Git repository.
type GitCommit struct {
	Headers map[string]string
	Sha     string
	Message string
	Files   []GitFileStatus
}

// GitFileStatus represents the status of a file in a Git commit.
type GitFileStatus struct {
	Status string
	File   string
}

func main() {
	var repo string
	flag.StringVar(&repo, "repo", "", "specify the repository")

	var since string
	flag.StringVar(&since, "since", "", "specify the since time")

	var author string
	flag.StringVar(&author, "author", "", "specify the since time")

	flag.Parse()

	if repo == "" {
		flag.Usage()
		fmt.Println("Verbose output enabled. Current Arguments: -v \n", since)
		fmt.Println("Quick Start Example! App is in Verbose mode!")
	} else {
		fmt.Printf("Current Arguments: -v %s\n", since)
		fmt.Println("Quick Start Example!")
		ChDir(repo)
	}

	output := AllLogs(since, author)
	fmt.Println(output)

	commits := ParseResults(output)
	fmt.Println(commits)

	entries := []string{}
	fmt.Println("Messages: ")
	for _, c := range commits {
		fmt.Println(c.Message)
		// check for regex #{number} and JIRA-1 test abc-2
		// ([\S]+) matches words and -\d+ matches -1
		re := regexp.MustCompile(`([\S]+)-\d+`)
		matches := re.FindAllString(c.Message, -1)
		for _, match := range matches {
			entries = append(entries, match)
		}

		// check for regex #{number}
		re = regexp.MustCompile(`#\d+`)
		matches = re.FindAllString(c.Message, -1)
		for _, match := range matches {
			entries = append(entries, match)
		}
	}

	fmt.Println("----------------")
	fmt.Println("Issues found: ")
	for _, e := range entries {
		fmt.Println(e)
	}
}

// CommandLineOptions represents the command-line options for the program.
type CommandLineOptions struct {
	Since  string `kong:"help='Since Time', default='yesterday'"`
	Author string `kong:"help='Author to search git logs for', default='David Li'"`
	Repo   string `kong:"help='local path to repository to parse'"`
}

// AllLogs returns the output of the "git log" command.
func AllLogs(since, author string) string {
	cmd := exec.Command("git", "log", "--since", since, "--author", author)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(output)
}

// ChDir changes the current working directory.
func ChDir(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		fmt.Println(err)
	}
}

// ParseResults parses the output of the "git log" command and returns a slice of GitCommit objects.
func ParseResults(output string) []*GitCommit {
	commits := []*GitCommit{}

	scanner := bufio.NewScanner(strings.NewReader(output))
	commit := &GitCommit{Headers: map[string]string{}}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// end of commit
			commits = append(commits, commit)
			commit = &GitCommit{Headers: map[string]string{}}
		} else if strings.HasPrefix(line, "commit ") {
			commit.Sha = strings.TrimPrefix(line, "commit ")
		} else if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			commit.Headers[key] = value
		} else {
			commit.Message += line + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return commits
}
