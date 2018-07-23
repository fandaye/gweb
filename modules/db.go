package modules

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type DB struct {
	Conn   sql.DB
	Config map[string]string
}

func (db *DB) Connect() (err error) {
	mysqlConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", db.Config["mysql_username"], db.Config["mysql_password"], db.Config["mysql_host"], db.Config["mysql_port"], db.Config["mysql_dbname"], db.Config["mysql_charset"])
	connect, err := sql.Open("mysql", mysqlConn)
	if err == nil {
		db.Conn = *connect //数据库连接
		return nil
	} else {
		return err
	}
}

//插入  -------
func (db *DB) Insert(sqlstr string, args ...interface{}) (int64, error) {
	stmtIns, err := db.Conn.Prepare(sqlstr)
	if err != nil {
		return 1, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 1, err
	}
	return result.LastInsertId()
}

//更新
func (db *DB) Update(sqlstr string, args ...interface{}) (int64, error) {
	stmtIns, err := db.Conn.Prepare(sqlstr)
	if err != nil {
		return 1, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 1, err
	}
	return result.RowsAffected()
}

//删除
func (db *DB) Delete(sqlstr string, args ...interface{}) (int64, error) {
	stmtIns, err := db.Conn.Prepare(sqlstr)
	if err != nil {
		return 1, err
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 1, err
	}
	return result.RowsAffected()
}

//查询sql函数，传sql
func (db *DB) SelectAll(sqlstr string, args ...interface{}) (*[]map[string]string, error) {
	ret := make([]map[string]string, 0)

	stmtOut, err := db.Conn.Prepare(sqlstr)
	if err != nil {
		return &ret, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		return &ret, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return &ret, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return &ret, err
		}
		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return &ret, nil
}

//取一行数据，注意这类取出来的结果都是string  --------
func (db *DB) SelectOne(sqlstr string, args ...interface{}) (map[string]string, error) {
	stmtOut, err := db.Conn.Prepare(sqlstr)
	if err != nil {
		return make(map[string]string), err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		return make(map[string]string), err
	}
	columns, err := rows.Columns()
	if err != nil {
		return make(map[string]string), err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	ret := make(map[string]string, len(scanArgs))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		var value string

		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			ret[columns[i]] = value
		}
		break //get the first row only
	}
	return ret, nil
}