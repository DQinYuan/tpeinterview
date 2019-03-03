package shellcmd

import (
	"fmt"
	"github.com/DQinYuan/tpeinterview/db"
	"github.com/spf13/cobra"
	"strconv"
)


func AuthAction(cmd *cobra.Command, args []string) {

	if writeNum == -1{
		panic("please write some record first")
	}

	mysqlDB := db.CreateDB()
	if mysqlDB == nil{
		return
	}

	result, err := mysqlDB.QueryAll()
	if err != nil{
		panic("auth fail, query fail")
	}

	if len(result) != writeNum{
		panic(
			fmt.Sprintf("auth fail, expect record num %d, real write num %d", writeNum, len(result)))
	}

	for _, r := range result{
		ele, err := strconv.Atoi(r[11])
		if err != nil{
			panic(
				fmt.Sprintf("auth fail, expect record num , data is not expected, error data: %v", r))
		}
		record := genRecord(ele)
		if !arrEqual(record[11:], r){
			panic(
				fmt.Sprintf("auth fail, expect record: %v, error data: %v", record, r))
		}
	}

	fmt.Println("OK")
}

func arrEqual(strs1 []string , strs2 []string) bool {
	for i, v := range strs1{
		if v != strs2[i]{
			return false
		}
	}

	return true
}