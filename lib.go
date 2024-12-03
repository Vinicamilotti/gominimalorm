package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Vinicamilotti/gominimalorm/reflection"
	errHelper "github.com/Vinicamilotti/gominimalorm/utils"
	_ "github.com/lib/pq"
)

type DbHandler struct {
	Dns  string
	Conn *sql.DB
}

func (d *DbHandler) Open() {
	d.Conn = errHelper.MustReturn(sql.Open("postgres", d.Dns))
}

func scan_row[T interface{}](rows *sql.Rows) (T, error) {
	var dest T
	fields := reflection.StructFieldPtr(&dest)
	err := rows.Scan(fields...)
	return dest, err
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
	return scan_row[T](rows)
}

type Test struct {
	Id    string
	Nome  string
	Idade int
}

func main() {
	db := &DbHandler{
		Dns: "host=localhost user=postgres password=0303 dbname=teste port=5432 sslmode=disable TimeZone=America/Sao_Paulo",
	}
	db.Open()
	res, err := Query[Test](db, "SELECT * FROM teste")
	if err != nil {
		panic(err)
	}
	for _, v := range res {
		fmt.Println(v.Nome)
	}

	res2, err := QuerySingle[Test](db, "SELECT * FROM teste WHERE nome = $1", "vini")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res2.Nome)

}
