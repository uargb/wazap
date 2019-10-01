package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Manager struct {
	gorm.Model
	Username     string
	Password     string
	Active       bool   `gorm:"default:'1'"`
	BotPhone     string `gorm:"default:''"`
	Name         string
	Phone        string
	LinkTemplate string `gorm:"default:''"`
	Greeting     string `gorm:"default:''"`
}

type QA struct {
	gorm.Model
	ManagerID  int
	Index      int
	Query      string
	Text       string `gorm:"default:''"`
	Image      string `gorm:"default:''"`
	Video      string `gorm:"default:''"`
	Attachment string `gorm:"default:''"`
	Next       string `gorm:"default:''"`
	Write      string `gorm:"default:''"`
}

type Costumer struct {
	gorm.Model
	ManagerID int
	Phone     string
	Name      string
	Fields    string
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		"admin",
		"952368741",
		"wazap.cvtuyrclurh0.ap-south-1.rds.amazonaws.com:3306",
		"gym",
	))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database %v", err))
	}

	db.AutoMigrate(&Manager{}, &QA{}, &Costumer{})

	/* db.Create(&Manager{
		Username: "test",
		Password: "test",
		Name:     "TEST",
		Phone:    "380970966546",
	}) */

	r := gin.Default()

	az := r.Group("/:user/:pwd")
	{
		az.GET("/check", check)
	}

	r.Run("0.0.0.0:8090")
}
