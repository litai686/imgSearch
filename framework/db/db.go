package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//var err error
//var engin *gorose.Engin

func init() {
	// 全局初始化数据库,并复用
	// 这里的engin需要全局保存,可以用全局变量,也可以用单例
	// 配置&gorose.Config{}是单一数据库配置
	// 如果配置读写分离集群,则使用&gorose.ConfigCluster{}
	// mysql Dsn示例 "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=true"
	//engin, err = gorose.Open(&gorose.Config{Driver: "sqlite3", Dsn: "./dfs.db"})
}

//
//func DB() gorose.IOrm {
//	return engin.NewOrm()
//}

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
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(sql)
	if err == nil {
		return rows, nil
	} else {
		println("sql:", sql)
		return nil, err
	}
}

func Insert_sql(sql string) (int64, error) {
	db, err := connect()
	if err != nil {
		return 0, err
	}
	//tx, err := db.Begin()
	//if err != nil {
	//	return 0, err
	//}
	res, err := db.Exec(sql)
	if err != nil {
		return 0, err
	}
	//err = tx.Commit()
	//if err != nil {
	//	err = tx.Rollback()
	//	return 0, err
	//}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func Update_sql(sql string) (int64, error) {
	db, err := connect()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(sql)
	if err != nil {
		return 0, err
	}
	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, err
}

func Exec_sql(sql string) (int64, error) {
	db, err := connect()
	if err != nil {
		return 0, err
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(sql)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		return 0, err
	}
	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, err
}
