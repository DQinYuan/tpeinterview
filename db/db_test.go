package db

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMysqlDB(t *testing.T) {
	t.SkipNow()
	db := CreateDB()
	if db == nil{
		return
	}

	field := make([]string, 0, 11)
	record := make([]string, 0, 11)

	for i := 0; i < 11; i++{
		field = append(field, "1")
		record = append(record, strings.Join(field, ""))
	}

	db.Insert(record)

	all, err := db.QueryAll()
	if err != nil{
		fmt.Println("query err")
		os.Exit(1)
	}

	for _, row := range all{
		for _, field := range row{
			fmt.Printf("%s  ", field)
		}
		fmt.Println()
	}
}
