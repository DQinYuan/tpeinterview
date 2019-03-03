package main

import (
	"fmt"
	"github.com/DQinYuan/tpeinterview/db"
	"github.com/DQinYuan/tpeinterview/dockerctl"
	"github.com/DQinYuan/tpeinterview/shellcmd"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
	"time"
)

/**
实现计划:

1. 用户可以通过命令行逐条解释执行相关的测试命令, 根据容器名kill容器 (允许的命令: 写入数据, 读取验证数据正确, 读取验证数据错误)
2. 也可以通过load命令load多个test_case(允许使用通配符加载多个),每次执行完一个case能够进行现场复原
3. 线索:关闭一个tikv数据依旧是完好的, 关闭两个就会出问题了。
        docker-compose每次启动的都是同一个容器, 并且只会启动当前缺失的容器


支持的命令:
    write 记录数
    kill 容器名
    auth
    load 文件名(允许使用*进行前后缀匹配)
    run
    nodes

 */


func main() {
	rootCmd := &cobra.Command{
		Use:   "tpeinterview docker-compose_path",
		Short: "a demo tool to test tidb",
		//至少要给出一个参数,即docker-compose的位置
		Args: cobra.MinimumNArgs(1),
		DisableFlagsInUseLine:true,
		Run: func(cmd *cobra.Command, args []string) {

			//启动tidb
			fmt.Println("tidb starting..., Please wait some seconds")
			dockerctl.DockerCompose(args[0], "up", "-d")

			//尝试建立连接,以验证tidb是否真正启动完毕
			for {
				testDB := db.CreateDB()
				if testDB != nil{
					break
				}
				time.Sleep(time.Second * 10)
			}


			l, err := readline.NewEx(&readline.Config{
				Prompt:            "\033[31m»\033[0m ",
				HistoryFile:       "/tmp/readline.tmp",
				InterruptPrompt:   "^C",
				EOFPrompt:         "^D",
				HistorySearchFold: true,
			})
			if err != nil {
				fmt.Println(err)
			}
			defer l.Close()


			for {
				line, err := l.Readline()
				if err != nil {
					if err == readline.ErrInterrupt {
						return
					} else if err == io.EOF {
						return
					}
					continue
				}
				if line == "exit" {
					fmt.Println("stopping tidb..., please wait some seconds")
					dockerctl.DockerCompose("", "stop", "")
					return
				}

				handleLine(line)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}

func handleLine(line string) {
	defer func() {
		if err := recover(); err != nil{
			fmt.Println(err)
		}
	}()

	args := strings.Split(strings.TrimSpace(line), " ")
	shellcmd.RunCmd(args)
}