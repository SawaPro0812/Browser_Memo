package main

import (
    "github.com/gin-gonic/gin"
    "fmt"
    "time"
    "crypto/sha512"
    "net/url"
    "net/http"

    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

var r *gin.Engine

// memosテーブルモデル宣言
type memos struct {
	gorm.Model
	ID  uint    `gorm:"primaryKey"`
	Memo string
	Hash_id string
}

// DB接続情報
func gormConnect() *gorm.DB {
    dsn := "root:mysql@tcp(db:3306)/browser_memo?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
	    fmt.Println("DB_STATUS: Error")
		panic(err.Error())
	}

	return db
}

func main() {
    // DB接続
    db := gormConnect()
    // テーブルのマイグレーション
    db.AutoMigrate(&memos{})
    //defer db.Close()
    //db.LogMode(true)

	r := gin.Default()
	r.Static("/css", "view/css")
	r.Static("/js", "view/js")
	r.LoadHTMLGlob("view/html/*.html")

	// メモ画面を表示する
	r.GET("/memo", memoTop)
    // メモを登録する
    r.POST("/memo/create", memoCreate)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// メモ画面初期表示
func memoTop(c *gin.Context) {
    // リクエスト
    r := c.Request
    // URL取得
    myUrl, _ := url.Parse(r.URL.String())
    // クエリパラメータ取得
    params, _ := url.ParseQuery(myUrl.RawQuery)

    // hashIdを設定する
    var hashId string
    if (len(params["hashId"]) > 0) {
            hashId = params["hashId"][0]
    }

    memos, _ := selectMemosByHashId(hashId)

    // システム日時を取得する
    t := time.Now().String()
    if len(hashId) > 0 {
        c.HTML(200, "browser_memo.html", gin.H{"sysTime": t, "memos": memos, "hashId": hashId})
    } else {
        c.HTML(200, "browser_memo.html", gin.H{"sysTime": t})
    }
}

// メモ画面初期表示
func memoCreate(c *gin.Context) {
    memo := c.PostForm("memo")
    sysTime := c.PostForm("sysTime")
    hashId := c.PostForm("hashId")
    if len(hashId) == 0 {
        // hashIdを生成する
        hashId = createHashId(sysTime)
    }

    // DB接続情報取得
    db := gormConnect()
    db.Create(&memos{Memo: memo, Hash_id: hashId})
    c.Redirect(http.StatusFound, "/memo?hashId=" + hashId)
}

// hashIdを生成する
func createHashId(sysTime string) string {
    // システム日時を取得する
    t := time.Now().String()
    // メモ画面初期表示時に設定したシステム日時と結合する
    joinTime := sysTime + t
    p := []byte(joinTime)
    // ハッシュ化する
    sha512 := sha512.Sum512(p)
    return fmt.Sprintf("%x", sha512)
}

// ハッシュIDをキーとしてメモのリストを取得する
func selectMemosByHashId(hashId string) ([]memos, error) {
    // DB接続情報取得
    db := gormConnect()
    var memosList []memos
    if len(hashId) > 0 {
        result := db.Where("hash_id = ?", hashId).Find(&memosList)
        if result.Error != nil {
            fmt.Println("Error:メモが取得できませんでした")
            return nil, result.Error
        }
        // rows affected でレコードの有無を確かめる
        if result.RowsAffected == 0 {
            return nil, nil
        }
        return memosList, nil
    }
    return nil, nil
}