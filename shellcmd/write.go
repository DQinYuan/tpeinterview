package shellcmd

import (
	"fmt"
	"github.com/DQinYuan/tpeinterview/db"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)


var writeNum = -1

func WriteAction(cmd *cobra.Command, args []string) {
	recordNum, err := strconv.Atoi(args[0])
	if err != nil{
		fmt.Println(cmd.UsageString())
		return
	}
	mysqlDb := db.CreateDB()
	if mysqlDb == nil{
		return
	}

	var base int
	if writeNum == -1 {base = 0} else {base = writeNum}

	for i := 0; i < recordNum; i++{
		mysqlDb.Insert(genRecord(base + i))
	}

	if writeNum == -1{
		writeNum = recordNum
	} else {
		writeNum += recordNum
	}
}

/*
生成类似这样的数据作为一条测试记录：

genRecord(9)

{
	"9",
	"99",
	"999",
	"9999",
	"99999",
	"999999",
	"9999999",
	"99999999",
	"999999999",
	"9999999999",
	"99999999999",
}
 */
func genRecord(num int) []string {
	ele := strconv.Itoa(num)

	fieldNum := db.FieldCount + 1

	field := make([]string, 0, fieldNum)
	record := make([]string, 0, fieldNum)

	for i := 0; i < fieldNum; i++{
		field = append(field, ele)
		record = append(record, strings.Join(field, ""))
	}

	return record
}
