package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func main() {
	dir := os.Args[1]
	git := GitCmd{dir}

	var err error
	for {
		time.Sleep(2 * time.Second)

		_, err = git.Run("add", ".")
		if err != nil {
			log.Fatal(err)
		}

		rawStatus, err := git.Run("status", "--short")
		if err != nil {
			log.Fatal(err)
		}

		if len(rawStatus) == 0 {
			fmt.Print(".")
			continue
		}

		status := ParseGitStatus(rawStatus)
		commitMsg := makeCommitMessage(status)

		_, err = git.Run("commit", "-m", commitMsg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("C")
	}
}

type GitCmd struct {
	dir string
}

func (g *GitCmd) Run(cmds ...string) (string, error) {
	allCmds := append([]string{"-C", g.dir}, cmds...)
	cmd := exec.Command("git", allCmds...)
	bout, err := cmd.Output()
	out := string(bout)
	return out, err
}

type StatusCode int

const (
	STATUS_UNKNOWN StatusCode = iota
	STATUS_ADDED
	STATUS_MODIFIED
)

type GitStatus struct {
	Items []GitStatusItem
}

type GitStatusItem struct {
	Code     StatusCode
	Filename string
}

func NewGitStatusItem(rawCode, filename string) GitStatusItem {
	code := STATUS_UNKNOWN
	switch rawCode {
	case "A":
		code = STATUS_ADDED
	case "M":
		code = STATUS_MODIFIED
	}

	filename = strings.TrimSpace(filename)

	// A filename is surrounded by quotes when it contains a space.
	if filename[0] == '"' {
		filename = filename[1 : len(filename)-1]
	}
	return GitStatusItem{code, filename}
}

func ParseGitStatus(status string) GitStatus {
	lines := strings.Split(status, "\n")

	// Remove the empty last item.
	lines = lines[0 : len(lines)-1]

	items := make([]GitStatusItem, len(lines))
	for i, l := range lines {
		parts := strings.SplitN(l, " ", 2)
		items[i] = NewGitStatusItem(parts[0], parts[1])
	}

	return GitStatus{items}
}

func makeCommitMessage(status GitStatus) string {
	items := append([]GitStatusItem{}, status.Items...)

	sort.Slice(items, func(ia, ib int) bool {
		a := items[ia]
		b := items[ib]
		if a.Code != b.Code {
			return a.Code < b.Code
		}
		return strings.Compare(a.Filename, b.Filename) == -1
	})

	msg := ""
	var lastCode StatusCode = -1
	for i, it := range items {
		if lastCode != it.Code {
			if i > 0 {
				msg += " , "
			}
			switch it.Code {
			case STATUS_ADDED:
				msg += "add "
			case STATUS_MODIFIED:
				msg += "update "
			case STATUS_UNKNOWN:
				msg += "?? "
			default:
				msg += "???? "
			}
			msg += it.Filename
			lastCode = it.Code
		} else {
			msg += fmt.Sprintf(", %s", it.Filename)
		}
	}

	return msg
}
