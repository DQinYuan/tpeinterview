package shellcmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func RemoveAction(cmd *cobra.Command, args []string) {
	caseName := args[0]
	if _, ok := scriptsToRun[caseName]; !ok{
		fmt.Printf("test case %s do not exist \n", caseName)
		return
	}

	delete(scriptsToRun, caseName)
}
