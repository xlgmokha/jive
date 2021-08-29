package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "cd into a project directory",
	Long:  "cd into a project directory. e.g $ jive cd xlgmokha/net-hippie",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("<owner/project> is needed")
			os.Exit(1)
		}

		nwo := args[0]
		user, _ := user.Current()
		host := "github.com"
		projectPath := path.Join(user.HomeDir, "src", host, nwo)
		_, err := os.Stat(projectPath)

		if os.IsNotExist(err) {
			command := exec.Command(
				"git",
				"clone",
				fmt.Sprintf("git@%s:%s.git", host, nwo),
				projectPath,
			)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Run()
		}

		tasks := [][]string{
			[]string{"cd", projectPath},
			[]string{"ctags", projectPath},
			[]string{"setenv", fmt.Sprintf("JIVE_LAST_RUN=%s", time.Now().Unix())},
		}
		afterRun(tasks)
	},
}

func afterRun(tasks [][]string) {
	for i, _ := range tasks {
		command := tasks[i][0]
		args := tasks[i][1]
		toSend := fmt.Sprintf("%s:%s\n", command, args)

		_, err := syscall.Write(42, []byte(toSend))
		if err != nil {
			fmt.Errorf("%v", err)
			break
		}
	}
}

func init() {
	rootCmd.AddCommand(cdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
