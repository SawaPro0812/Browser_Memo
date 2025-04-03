package main

import (
    "github.com/gin-gonic/gin"
    "fmt"
    "time"
    "crypto/sha512"

    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

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
	r.LoadHTMLGlob("view/html/*.html")

	// メモ画面を表示する
	r.GET("/memo", memoTop)
    // メモを登録する
    r.POST("/memo/create", memoCreate)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// メモ画面初期表示
func memoTop(c *gin.Context) {
    // DB接続情報取得
    db := gormConnect()
    //hashId := "6a50c6c7868cd35ce9cf32c64689247198ab387910fe6c6f389a1bc2d8935fbd2b26540f959fefbe66bb41861de4df6612eb0a96f373d7dd2b04d8e34323c114"
    //list := db.Find(&memos{}, "hash_id=?", hashId)
    // 構造体を宣言？
    memos := memos{}

    result := db.First(&memos)
    if result.Error != nil {
        fmt.Println("Error:", result.Error)
    } else {
        fmt.Printf("Memo: %+v\n", memos) // 取得したデータを表示
    }

    // システム日時を取得する
    t := time.Now().String()
    c.HTML(200, "browser_memo.html", gin.H{"sysTime": t, "memo": memos})
}

// メモ画面初期表示
func memoCreate(c *gin.Context) {
    memo := c.PostForm("memo")
    sysTime := c.PostForm("sysTime")
    // hashIdを生成する
    hashId := createHashId(sysTime)
    fmt.Println(hashId)
    // DB接続情報取得
    db := gormConnect()
    db.Create(&memos{Memo: memo, Hash_id: hashId})
    c.HTML(200, "browser_memo.html", gin.H{"memo": memo})
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