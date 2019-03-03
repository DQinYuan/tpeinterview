package shellcmd

import (
	"fmt"
	"github.com/DQinYuan/tpeinterview/dockerctl"
	"github.com/spf13/cobra"
)

func NodesAction(cmd *cobra.Command, args []string) {
	client := dockerctl.CreateClient()
	if client == nil{
		return
	}

	containers, err := client.GetAllStartingContainer()
	if err != nil{
		fmt.Println("list containers error")
		return
	}

	for _, container := range containers{
		fmt.Println(container.Names[0][1:])
	}
}
