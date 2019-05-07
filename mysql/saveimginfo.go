package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	USERNAME = ""
	PASSWORD = "*****"
	NETWORK  = "tcp"
	SERVER   = "******"
	PORT     = 3306
	DATABASE = "*****"
)

type Img struct {
	ID         int32     `db:"id"`
	Name       string    `db:"name"`
	Dir        int       `db:"dir"`
	Alias      string    `db:"alias"`
	Desc       string    `db:"desc"`
	UploadTime time.Time `db:"upload_time"`
}

func getConn() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Open mysql failed,err:%v\n", err)
		return nil
	}
	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(10)                   //设置最大连接数
	DB.SetMaxIdleConns(4)                    //设置闲置连接数
	return DB
}

func closeConn(db *sql.DB) bool {
	err := db.Close()
	if err != nil {
		return false
	}
	return true
}

//插入图片记录信息
func InsertImg(img *Img) error {
	DB := getConn()
	result, err := DB.Exec("insert INTO imgs(`name`,`dir`,`alias`,`desc`,`upload_time`) values(?,?,?,?,?)", img.Name, img.Dir, img.Alias, img.Desc, img.UploadTime)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		closeConn(DB)
		return err
	}
	lastInsertID, err := result.LastInsertId() //插入数据的主键id
	fmt.Println("LastInsertID:", lastInsertID)
	if err != nil {
		fmt.Printf("Get lastInsertID failed,err:%v", err)
		closeConn(DB)
		return err
	}
	closeConn(DB)
	return nil
}

//根据别名获取图片信息
func QueryByAlias(alias string) (Img, error) {
	DB := getConn()
	img := Img{}
	querySql := "select * from imgs where alias= '" + alias + "'"
	row := DB.QueryRow(querySql)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	err := row.Scan(&img.ID, &img.Name, &img.Dir, &img.Alias, &img.Desc, &img.UploadTime)
	if err != nil {
		fmt.Printf("scan failed, err:%v", err)
		closeConn(DB)
		return img, err
	}
	closeConn(DB)
	return img, err
}

func QueryMulti() {
	DB := getConn()
	img := Img{}
	rows, err := DB.Query("select * from imgs")
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&img.ID, &img.Name, &img.Dir, &img.Alias, &img.Desc, &img.UploadTime) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return
		}
		fmt.Println(img)
	}
}
