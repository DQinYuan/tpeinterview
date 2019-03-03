package shellcmd

import (
	"fmt"
	"github.com/DQinYuan/tpeinterview/dockerctl"
	"github.com/spf13/cobra"
	"strings"
)

func RunAction(cmd *cobra.Command, args []string) {

	if len(scriptsToRun) < 1{
		fmt.Println("there is no test case, please load first")
		return
	}

	failCases := make([]string, 0)

	for path, content := range scriptsToRun{
		fmt.Printf("======start run test case: %s \n", path)
		runScript(path, content, &failCases)
	}

	fmt.Printf("\n\n fail cases: %v \n", failCases)
}

func runScript(path string, content string, cases *[]string) {

	var lineNum int

	//脚本执行出错,测试用例失败,将其捕获
	defer func() {
		if err := recover(); err != nil{
			*cases = append(*cases, path)
			fmt.Printf("======%s test case fail, line num %d fail, error message: %s \n", path, lineNum, err)
		}
	}()

	//还原现场,再跑新的test case
	dockerctl.DockerCompose("", "up", "-d")

	//逐行执行脚本
	for i, lineStr := range strings.Split(content, "\n") {
		lineNum = i
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			return
		}
		args := strings.Split(lineStr, " ")
		RunCmd(args)
	}

	fmt.Printf("======%s test case success \n", path)
}




func RunCmd(args []string) {
	cmd := &cobra.Command{
		Use:   "",
		Short: "shell command",
	}


	cmd.SetArgs(args)
	cmd.ParseFlags(args)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "auth",
			Short: "assert db is readable and data is the same as what is write",
			Run: AuthAction,
			DisableFlagsInUseLine: true,
		},
		&cobra.Command{
			Use:   "kill container_name",
			Short: "kill tidb node container by name",
			Args:  cobra.MinimumNArgs(1),
			Run:   KillAction,
			DisableFlagsInUseLine: true,
		},
		&cobra.Command{
			Use:   "load file_name",
			Short: "load a script file as test case",
			Args:  cobra.MinimumNArgs(1),
			Run:   LoadAction,
			DisableFlagsInUseLine: true,
		},
		&cobra.Command{
			Use:   "run",
			Short: "run all loaded test case",
			Run:   RunAction,
			DisableFlagsInUseLine: true,
		},
		&cobra.Command{
			Use:   "write record_num",
			Short: "write record_num record into database",
			Args:  cobra.MinimumNArgs(1),
			Run:   WriteAction,
			DisableFlagsInUseLine: true,
		},
		&cobra.Command{
			Use:   "clean",
			Short: "clean test cases loaded",
			Run:   CleanAction,
			DisableFlagsInUseLine: true,
		},
		&cobra.Command{
			Use: "nodes",
			Short: "list all starting nodes name",
			Run: NodesAction,
			DisableFlagsInUseLine:true,
		},
		&cobra.Command{
			Use: "list",
			Short: "list all test cases",
			Run:ListAction,
			DisableFlagsInUseLine:true,
		},
		&cobra.Command{
			Use:"remove test_case_path",
			Short:"remove a loaded test case",
			Run:RemoveAction,
			Args:  cobra.MinimumNArgs(1),
			DisableFlagsInUseLine:true,
		},
	)

	if err := cmd.Execute(); err != nil {
		fmt.Println(cmd.UsageString())
	}
}
