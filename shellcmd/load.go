package shellcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

// 脚本文件路径 -> 文件内容
var scriptsToRun = make(map[string]string)

func LoadAction(cmd *cobra.Command, args []string) {
	path := args[0]
	files, err := filepath.Glob(path)
	if err != nil{
		fmt.Println("file path error")
		return
	}

	for _, file := range files{
		bytes, err := ioutil.ReadFile(file)
		if err != nil{
			fmt.Printf("file %s open fail", file)
			continue
		}

		scriptsToRun[file] = string(bytes)
	}

	ListAction(cmd, args)
}