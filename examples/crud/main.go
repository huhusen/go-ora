package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/sijms/go-ora/v2"

	"os"
	"time"
)

func dieOnError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}
func createTable(conn *sql.DB) {
	fmt.Println("Creating temporary table GOORA_TEMP_VISIT")
	sqlText := `CREATE TABLE GOORA_TEMP_VISIT(
	VISIT_ID	number(10)	NOT NULL,
	NAME		VARCHAR(200),
	VAL			number(10,2),
	VISIT_DATE	date,
	PRIMARY KEY(VISIT_ID)
	)`
	_, err := conn.Exec(sqlText)
	dieOnError("Cannot create temporary table", err)
}

func insertData(conn *sql.DB) {
	fmt.Println("Inserting values in the table")
	index := 1
	stmt, err := conn.Prepare(`INSERT INTO GOORA_TEMP_VISIT(VISIT_ID, NAME, VAL, VISIT_DATE) 
VALUES(:1, :2, :3, :4)`)
	dieOnError("Cannot prepare stmt for insert", err)
	defer func() {
		_ = stmt.Close()
	}()
	nameText := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	val := 1.1
	for index = 1; index <= 100; index++ {
		_, err = stmt.Exec(index, nameText, val, time.Now())
		errorText := fmt.Sprintf("Error during insert at index: %d", index)
		dieOnError(errorText, err)
		val += 1.1
	}
	fmt.Println("100 Rows inserted")
}

func queryData(conn *sql.DB) {
	fmt.Println("Query rows")
	rows, err := conn.Query("SELECT VISIT_ID, NAME, VAL, VISIT_DATE FROM GOORA_TEMP_VISIT")
	dieOnError("Cannot query after insert", err)
	var (
		id   int64
		name string
		val  float32
		date time.Time
	)
	for rows.Next() {
		err = rows.Scan(&id, &name, &val, &date)
		dieOnError("Cannot scan rows", err)
		fmt.Println("ID: ", id, "\tName: ", name, "\tval: ", val, "\tDate: ", date)
	}
	dieOnError("Query: ", rows.Err())
	fmt.Println("Finish query")
}

func updateData(conn *sql.DB) {
	fmt.Println("Updating values")
	updStmt, err := conn.Prepare(`UPDATE GOORA_TEMP_VISIT SET NAME = :1 WHERE VISIT_ID = :2`)
	dieOnError("Can't prepare stmt for update", err)
	nameText := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defer func() {
		_ = updStmt.Close()
	}()
	for index := 1; index <= 100; index++ {
		_, err = updStmt.Exec(nameText[:101-index], index)
		dieOnError("Can't update", err)
	}
	fmt.Println("Finish update")
}

func deleteData(conn *sql.DB) {
	fmt.Println("Deleting data")
	_, err := conn.Exec("delete from GOORA_TEMP_VISIT")
	dieOnError("Can't delete", err)
	fmt.Println("Finish delete")
}

func dropTable(conn *sql.DB) {
	fmt.Println("Dropping table")
	_, err := conn.Exec("drop table GOORA_TEMP_VISIT purge")
	dieOnError("Can't drop table", err)
	fmt.Println("Finish drop table")
}
func usage() {
	fmt.Println()
	fmt.Println("crud")
	fmt.Println("  a complete code of create table insert, update, query and delete then drop table.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println(`  curd -server server_url`)
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println(`  crud -server "oracle://user:pass@server/service_name"`)
	fmt.Println()
}

func main() {
	var (
		server string
	)

	flag.StringVar(&server, "server", "", "Server's URL, oracle://user:pass@server/service_name")
	flag.Parse()

	connStr := os.ExpandEnv(server)
	if connStr == "" {
		fmt.Println("Missing -server option")
		usage()
		os.Exit(1)
	}
	fmt.Println("Connection string: ", connStr)
	conn, err := sql.Open("oracle", server)
	dieOnError("Can't open the driver:", err)

	defer func() {
		_ = conn.Close()
	}()

	err = conn.Ping()
	dieOnError("Can't ping connection", err)

	createTable(conn)

	defer dropTable(conn)

	insertData(conn)
	queryData(conn)
	updateData(conn)
	queryData(conn)
	deleteData(conn)

}