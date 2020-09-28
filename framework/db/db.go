package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//连接数据库
func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./dfs.db")
	if err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

func Query_sql(sql string) (*sql.Rows, error) {
	db, err := connect()
	rows, err := db.Query(sql)
	defer db.Close()
	if err == nil {
		return rows, nil
	} else {
		println("sql:", sql)
		return nil, err
	}
}

func Insert_sql(sql string) (int64, error) {
	db, err := connect()
	defer db.Close()
	chk(err)
	tx, err := db.Begin()
	chk(err)
	res, err := db.Exec(sql)
	chk(err)
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		chk(err)
	}
	id, err := res.LastInsertId()
	chk(err)
	return id, err
}

func Update_sql(sql string) (int64, error) {
	db, err := connect()
	defer db.Close()
	chk(err)
	tx, err := db.Begin()
	chk(err)
	res, err := db.Exec(sql)
	chk(err)
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		chk(err)
	}
	id, err := res.RowsAffected()
	chk(err)
	return id, err
}

func Exec_sql(sql string) (int64, error) {
	db, err := connect()
	defer db.Close()
	chk(err)
	tx, err := db.Begin()
	chk(err)
	res, err := db.Exec(sql)
	chk(err)
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		chk(err)
	}
	id, err := res.RowsAffected()
	chk(err)
	return id, err
}

//检查错误
func chk(err error) {
	if err != nil {
		panic(err)
	}
}
