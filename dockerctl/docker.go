package dockerctl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/DQinYuan/tpeinterview/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"os"
	"os/exec"
	"time"
)

type dockerCtl struct {
	cli *client.Client
}

var ctl *dockerCtl

var onceInterruptable util.OnceInterruptable

func CreateClient() *dockerCtl {

	defer func() {
		if err := recover() ;err != nil{
			fmt.Println(err)
		}
	}()

	onceInterruptable.Do(func() {
		cli, err := client.NewClient(client.DefaultDockerHost, "v1.12", nil, nil)
		if err != nil {
			panic("docker client init error")
		}

		ctl = &dockerCtl{cli}
	})

	return ctl
}

func (ctl *dockerCtl) GetContainerByName(name string) (*types.Container, error) {
	args := filters.NewArgs()
	args.Add("name", name)

	containers, err := ctl.cli.ContainerList(context.Background(), types.ContainerListOptions{Filters:args})
	if err != nil{
		return nil, err
	}
	if len(containers) < 1{
		return nil, errors.New("no contain has name " + name)
	}

	return &containers[0], nil
}

// 根据容器id停止容器
func (ctl *dockerCtl) StopContainer(containerID string) error {
	timeout := time.Second * 10
	err := ctl.cli.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

func (ctl *dockerCtl) GetAllStartingContainer() ([]types.Container, error)  {
	return ctl.cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

var dockerComposePath *string = nil

/*
利用docker-compose启动或停止tidb
 */
func DockerCompose(path string, action string) {
	if dockerComposePath == nil{
		dockerComposePath = &path
	}

	out, err := execShell("docker-compose", "-f", *dockerComposePath, action, "-d")
	if err != nil{
		fmt.Println(out)
		os.Exit(1)
	}

	fmt.Printf("tidb %s ok \n", action)
}

//执行shell命令
func execShell(s string, args ...string) (string, error){
	cmd := exec.Command(s, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	return out.String(), err
}


