package db

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/DQinYuan/tpeinterview/util"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
	"time"
)

const tableName string = "testtable"
const FieldCount = 10

type mysqlDB struct {
	db      *sql.DB

	bufPool *BufPool

	stmtCache map[string]*sql.Stmt
}

func (db *mysqlDB) createTable() error {

	if _, err := db.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)); err != nil {
		return err
	}

	fieldLength := 100

	buf := new(bytes.Buffer)
	s := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (TEST_KEY VARCHAR(64) PRIMARY KEY", tableName)
	buf.WriteString(s)

	for i := 0; i < FieldCount; i++ {
		buf.WriteString(fmt.Sprintf(", FIELD%d VARCHAR(%d)", i, fieldLength))
	}

	buf.WriteString(");")


	fmt.Println(buf.String())

	_, err := db.db.Exec(buf.String())
	return err
}

func (db *mysqlDB) getAndCacheStmt(query string) (*sql.Stmt, error) {
	state := db.stmtCache

	if stmt, ok := state[query]; ok {
		return stmt, nil
	}

	//产生Prepared statement
	stmt, err := db.db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, err
	}

	state[query] = stmt
	return stmt, nil
}

func (db *mysqlDB) execQuery(query string, args ...interface{}) error {
	fmt.Printf("%s %v\n", query, args)


	stmt, err := db.getAndCacheStmt(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(context.Background(), args...)
	db.clearCacheIfFailed(query, err)
	return err
}

func (db *mysqlDB) clearCacheIfFailed(query string, err error) {
	if err == nil {
		return
	}

	state := db.stmtCache
	delete(state, query)
}


func (db *mysqlDB) Insert(record []string) error {
	args := make([]interface{}, 0, len(record))
	for _, r:= range record{
		args = append(args, r)
	}

	buf := db.bufPool.get()
	defer db.bufPool.put(buf)

	buf.WriteString("INSERT INTO ")
	buf.WriteString(tableName)
	buf.WriteString(" (TEST_KEY")

	for i := 0; i < FieldCount; i++ {
		buf.WriteString(fmt.Sprintf(", FIELD%d", i))
	}
	buf.WriteString(") VALUES (?")

	for i := 0; i < FieldCount; i++ {
		buf.WriteString(" ,?")
	}

	buf.WriteByte(')')

	return db.execQuery(buf.String(), args...)
}

func (db *mysqlDB) QueryAll() ([][]string, error) {
	sql := "select * from " + tableName
	fmt.Println(sql)

	//10s没有返回就算超时失败
	timeoutContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := db.db.QueryContext(timeoutContext, sql)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil{
		return nil, err
	}

	result := make([][]string, 0)
	for rows.Next(){
		row := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++{
			row[i] = new(string)
		}
		if err = rows.Scan(row...); err != nil {
			return nil, err
		}

		strRow := make([]string, len(cols))
		for _, field := range row{
			strRow = append(strRow, *field.(*string))
		}

		result = append(result, strRow)
	}

	return result, rows.Err()
}

// BufPool is a bytes.Buffer pool
type BufPool struct {
	p *sync.Pool
}

// newBufPool creates a buffer pool.
func newBufPool() *BufPool {
	//p是Pool是地址
	p := &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	return &BufPool{
		p: p,
	}
}

// get gets a buffer.
func (b *BufPool) get() *bytes.Buffer {
	buf := b.p.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// put returns a buffer.
func (b *BufPool) put(buf *bytes.Buffer) {
	b.p.Put(buf)
}

var mydb *mysqlDB

var onceInterruptable util.OnceInterruptable

//创建测试用数据库,顺便建空表
func CreateDB() (*mysqlDB) {

	defer func() {
		if err := recover() ;err != nil{
			fmt.Println(err)
		}
	}()


	onceInterruptable.Do(func() {

		fmt.Println("test connect to tidb")

		d := new(mysqlDB)

		host := "127.0.0.1"
		port := 4000
		user := "root"
		password := ""
		dbName := "test"

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
		var err error
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Println("database conncect fail")
			os.Exit(1)
		}

		threadCount := 4
		db.SetMaxIdleConns(threadCount + 1)
		db.SetMaxOpenConns(threadCount * 2)
		d.db = db

		d.bufPool = newBufPool()

		if err := d.createTable(); err != nil {
			panic("table create fail, please waiting for tidb start up")
		}

		d.stmtCache = make(map[string]*sql.Stmt)

		if mydb != nil{
			mydb.db.Close()
		}

		mydb = d
	})

	return mydb
}

