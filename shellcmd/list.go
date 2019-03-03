package shellcmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func ListAction(cmd *cobra.Command, args []string) {
	fmt.Println("current test cases:")
	for path, _ := range scriptsToRun{
		fmt.Println(path)
	}
}