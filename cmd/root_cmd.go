package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	fqdn "github.com/Showmax/go-fqdn"
	"github.com/go-cmd/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:    "discordify",
	Short:  "Execute a shell command and notify about exit status in Discord or Slack",
	Run:    run,
	PreRun: checkRequired,
}

func checkRequired(cmd *cobra.Command, args []string) {
	if viper.GetString("webhook") == "" {
		cmd.Println("A webhook URL is required!")
		os.Exit(1)
	}
}

func run(c *cobra.Command, args []string) {
	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	// Create Cmd with options
	executable, arguments := args[0], args[1:]
	exec := cmd.NewCmdOptions(cmdOptions, executable, arguments...)
	var stdout strings.Builder
	var stderr strings.Builder

	// Print STDOUT and STDERR lines streaming from Cmd
	go func() {
		for {
			select {
			case line := <-exec.Stdout:
				fmt.Println(line)
				stdout.WriteString(line)
			case line := <-exec.Stderr:
				fmt.Fprintln(os.Stderr, line)
				stderr.WriteString(line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	<-exec.Start()

	// Cmd has finished but wait for goroutine to print all lines
	for len(exec.Stdout) > 0 || len(exec.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	c.Println(stdout.Len(), "bytes in stdout")
	c.Println(stderr.Len(), " bytes in stderr")
	c.Println("Execution finished, notifying Channel...")
	postResults(exec)
}

// Payload sent to webhook
type Payload struct {
	Content string `json:"content"`
}

func postResults(cmd *cmd.Cmd) {
	statusEmoji := ""
	if cmd.Status().Exit == 0 {
		statusEmoji = ":white_check_mark:"
	} else {
		statusEmoji = ":x:"
	}

	user, _ := user.Current()

	payload := Payload{
		Content: fmt.Sprintf("%v Your `%v` command on `%v` started by %v just finshed after %v seconds",
			statusEmoji,
			cmd.Name,
			fqdn.Get(),
			user.Username,
			cmd.Status().Runtime),
	}

	url := viper.GetString("webhook")

	jsonPayload, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("There was an error communicating with the API: ", err.Error())
	}
}
