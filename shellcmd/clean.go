package shellcmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func CleanAction(cmd *cobra.Command, args []string) {
	scriptsToRun = make(map[string]string)
	fmt.Println("clean test cases ok")
}

