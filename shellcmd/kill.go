package shellcmd

import (
	"fmt"
	"github.com/DQinYuan/tpeinterview/dockerctl"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var cli *client.Client

func KillAction(cmd *cobra.Command, args []string) {
	containerName := args[0]

	dockerCtl := dockerctl.CreateClient()
	if dockerCtl == nil {
		return
	}

	container, err := dockerCtl.GetContainerByName(containerName)
	if err != nil{
		fmt.Println("get Container named " + containerName + " fail")
		fmt.Println(err)
		return
	}

	err = dockerCtl.StopContainer(container.ID)
	if err != nil{
		fmt.Println("stop container " + containerName + " fail")
		return
	}

	fmt.Println("kill container " + containerName + " ok")
}
