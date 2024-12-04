package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Vinicamilotti/gominimalorm/reflection"
	errHelper "github.com/Vinicamilotti/gominimalorm/utils"
	_ "github.com/lib/pq"
)

type DbHandler struct {
	Dns  string
	Conn *sql.DB
}

func scan_row[T interface{}](rows *sql.Rows) (T, error) {
	var dest T
	fields := reflection.StructFieldPtr(&dest)
	err := rows.Scan(fields...)
	return dest, err
}

func generate_insert(fieldNames []string, tableName string) string {
	fields := strings.Join(fieldNames, ",")
	params := []string{}
	for i := range fieldNames {
		p := fmt.Sprintf("$%d", i)
		params = append(params, p)
	}
	sql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, fields, strings.Join(params, ","))
	return sql
}

func (d *DbHandler) Open() {
	d.Conn = errHelper.MustReturn(sql.Open("postgres", d.Dns))
}

func Query[T interface{}](db *DbHandler, sql string, parameters ...any) ([]T, error) {
	rows, err := db.Conn.Query(sql, parameters...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var ret []T = []T{}

	for rows.Next() {
		create, err := scan_row[T](rows)
		if err != nil {
			log.Println(err)
		}
		ret = append(ret, create)
	}

	return ret, nil
}

func QuerySingle[T interface{}](db *DbHandler, sql string, parameters ...any) (T, error) {
	rows, err := db.Conn.Query(sql, parameters...)
	if err != nil {
		log.Println(err)
		var def T
		return def, err
	}
	defer rows.Close()
	rows.Next()
	return scan_row[T](rows)
}

func SqlExec(db *DbHandler, sql string, parameters ...any) (sql.Result, error) {
	return db.Conn.Exec(sql, parameters...)
}

func Insert[T interface{}](db *DbHandler, tableName string, model T) (sql.Result, error) {
	fieldNames := reflection.StructFieldNames(model)
	fieldValues := reflection.StructFieldPtr(model)

	sql := generate_insert(fieldNames, tableName)

	return SqlExec(db, sql, fieldValues...)
}
