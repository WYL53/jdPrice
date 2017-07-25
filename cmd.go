package main

import (
	"bufio"
	"bytes"
	"os/exec"
	"time"
	"fmt"
)

func execCmd(command string, args ...string) *bytes.Buffer {
	ret := &bytes.Buffer{}
	retBuf := bufio.NewWriter(ret)
	cmd := exec.Command(command, args...)
	cmd.Stdout = retBuf
	cmd.Stderr = retBuf
	if err := cmd.Start();err != nil{
		fmt.Println("cmd.Start error:",err)
		return nil
	}
	go timeout2Kill(cmd,60,args[1])
	if err := cmd.Wait();err != nil{
		fmt.Println("cmd.Wait error:",err)
		return nil
	}
	return ret
}

//timeAfter秒后结束命令进程
func timeout2Kill(cmd *exec.Cmd, timeAfter uint,modelName string)  {
	defer func() {
		if err := recover();err != nil{
			fmt.Println("kill cmd err:",err)
		}
	}()

	var timer *time.Timer
	timer = time.AfterFunc(time.Duration(timeAfter)*time.Second, func() {
		timer.Stop()
		if cmd.Process != nil && cmd.ProcessState == nil{
			cmd.Process.Kill()
			fmt.Printf("查询型号【%s】的 cmd 被终结\n",modelName)
		}
	})
}