package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var Cmd *cobra.Command

// 创建支持秒级精度的Cron实例
var c = cron.New(cron.WithSeconds())

func init() {
	spec := ""
	cmdStr := ""
	s := ""
	Cmd = &cobra.Command{
		Use:   "cron",
		Short: "cron task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(spec) == 0 {
				return errors.New("spec is empty")
			}
			if len(cmdStr) == 0 {
				return errors.New("cmd is empty")
			}
			_, err := c.AddFunc(spec, func() {
				executeCommand(cmdStr)
			})
			if err != nil {
				return err
			}
			c.Start()
			return nil
		},
	}
	Cmd.Flags().StringVarP(&spec, "spec", "s", "", "example: * * * * * *")
	Cmd.Flags().StringVarP(&cmdStr, "cmd", "c", "", "example: echo 1")
	Cmd.Flags().StringVarP(&s, "_", "", "", "example:cron -s '* * * * * *' -c 'echo 1'")
}

func main() {
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
	}
	select {}
}

// 跨平台命令执行函数
func executeCommand(cmdStr string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	output, err := cmd.CombinedOutput()
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("\n[%s] Executing: %s\n", timestamp, cmdStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Output:\n%s\n", string(output))
}
