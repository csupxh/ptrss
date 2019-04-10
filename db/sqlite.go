package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zzyandzzy/ptrss/util"
    "os"
)

const (
	DataTableName   string = "t_data"
	ClientTableName string = "t_client"
	RSSTableName    string = "t_rss"
	createDataTable string = `
		CREATE TABLE t_data (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		tid INTEGER,
    		guid VARCHAR(64) UNIQUE,
    		hash VARCHAR(64),
    		title VARCHAR(128),
    		name VARCHAR(128),
    		url VARCHAR(256),
    		link VARCHAR(256),
    		rtype VARCHAR(64),
			size UNSIGNED BIG INT,
    		created VARCHAR(64)
		);
    `
	createClientTable string = `
		CREATE TABLE t_client (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  name VARCHAR(32),
		  local VARCHAR(64) UNIQUE,
		  user VARCHAR(64),
		  pwd VARCHAR(128)
		)
`
	createRSSTable string = `
		CREATE TABLE t_rss (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  url VARCHAR(256) UNIQUE,
		  cid INTEGER,
		  path VARCHAR(256),
		  pause INTEGER DEFAULT 1,
		  refresh INTEGER DEFAULT 300
		)
`
)

type Client struct {
	Id                     int
	Name, Local, User, Pwd string
}

type RSS struct {
	Id        int
	Url, Path string
	Pause     bool
	Refresh   int
}

var dbInstance *sql.DB

func Init() {
	if !ExistTable(DataTableName) {
		CreateTable(createDataTable)
	}
	if !ExistTable(ClientTableName) {
		CreateTable(createClientTable)
	}
	if !ExistTable(RSSTableName) {
		CreateTable(createRSSTable)
	}
}

func Instance() *sql.DB {
	if dbInstance != nil {
		return dbInstance
	} else {
        rssDir,_ := os.Getwd()
		dbInstance, err := sql.Open("sqlite3", rssDir + "./rss.db")
		util.CheckErr(err)
		return dbInstance
	}
}

// 判断表是否存在
func ExistTable(tableName string) bool {
	sqlStmt, err := Instance().Prepare("SELECT COUNT(*)  FROM main.sqlite_master WHERE type='table' AND name = ?")
	defer sqlStmt.Close()
	util.CheckErr(err)
	result, err := sqlStmt.Query(tableName)
	defer result.Close()
	util.CheckErr(err)
	for result.Next() {
		var id int
		result.Scan(&id)
		if id != 0 {
			return true
		}
	}
	return false
}

func CreateTable(sqlStmt string) {
	_, err := Instance().Exec(sqlStmt)
	util.CheckErr(err)
}

func InsertDatas(args ... interface{}) bool {
	stmt, err := Instance().Prepare(`INSERT INTO main.t_data
				(tid, guid, hash, title, name, url, link, rtype, size, created) 
				VALUES (?,?,?,?,?,?,?,?,?,?)`)
	util.CheckErr(err)
	return Insert(stmt, args...)
}

func InsertData(guid string, title string, url string, link string, rtype string, size int, created string) bool {
	stmt, err := Instance().Prepare(`INSERT INTO main.t_data(guid, title, url, link, rtype, size, created) 
				VALUES (?,?,?,?,?,?,?)`)
	util.CheckErr(err)
	return Insert(stmt, guid, title, url, link, rtype, size, created)
}

func InsertClient(args ... interface{}) bool {
	stmt, err := Instance().Prepare(`INSERT INTO main.t_client(name, local, user, pwd) VALUES (?,?,?,?)`)
	util.CheckErr(err)
	return Insert(stmt, args...)
}

func InsertRSS(args ... interface{}) bool {
	stmt, err := Instance().Prepare(`INSERT INTO main.t_rss(url, cid, path, pause, refresh) VALUES (?,?,?,?,?)`)
	util.CheckErr(err)
	return Insert(stmt, args...)
}

//插入数据
func Insert(stmt *sql.Stmt, args ... interface{}) bool {
	_, err := stmt.Exec(args...)
	defer stmt.Close()
	if err != nil {
		util.CheckErr(err)
		return false
	}
	return true
}

func QueryClient(clientName string) Client {
	stmt, err := Instance().Prepare("SELECT * FROM main.t_client WHERE name = ?")
	defer stmt.Close()
	util.CheckErr(err)
	var client Client
	stmt.QueryRow(clientName).Scan(&client.Id, &client.Name, &client.Local, &client.User, &client.Pwd)
	return client
}

func QueryRSS(clientName string) []RSS {
	stmt, err := Instance().Prepare(`SELECT t_rss.id,url,path,pause,refresh
				FROM main.t_rss,main.t_client 
				WHERE t_rss.cid = t_client.id AND t_client.name = ?`)
	defer stmt.Close()
	util.CheckErr(err)
	var rssArry [] RSS
	results, err := stmt.Query(clientName)
	defer results.Close()
	for results.Next() {
		var id int
		var url string
		var path string
		var pause bool
		var refresh int
		results.Scan(&id, &url, &path, &pause, &refresh)
		var rss = RSS{id, url, path, pause, refresh}
		rssArry = append(rssArry, rss)
	}
	return rssArry
}

func ExistData(guid string) bool {
	stmt, err := Instance().Prepare("SELECT id FROM main.t_data WHERE guid = ?")
	defer stmt.Close()
	util.CheckErr(err)
	var id int
	stmt.QueryRow(guid).Scan(&id)
	return id != 0
}

func ExistRSS(url string) bool {
	stmt, err := Instance().Prepare("SELECT id FROM main.t_rss WHERE url = ?")
	defer stmt.Close()
	util.CheckErr(err)
	var id int
	stmt.QueryRow(url).Scan(&id)
	return id != 0
}
