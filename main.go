package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	dir := os.Args[1]
	git := GitCmd{dir}

	for {
		timer1 := time.NewTimer(2 * time.Second)
		<-timer1.C

		statuses, err := git.Run("status", "--short")
		if err != nil {
			log.Fatal(err)
		}

		if len(statuses) == 0 {
			fmt.Print(".")
			continue
		}

		_, err = git.Run("add", ".")
		if err != nil {
			log.Fatal(err)
		}

		_, err = git.Run("commit", "-m", "update")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("commit")
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
