package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// QA describes Q&A card
type QA struct {
	gorm.Model
	ManagerID  uint
	Query      string
	Text       string `gorm:"default:''"`
	Image      string `gorm:"default:''"`
	Video      string `gorm:"default:''"`
	Attachment string `gorm:"default:''"`
	Next       string `gorm:"default:''"`
	Write      string `gorm:"default:''"`
}

// Costumer describes costumer
type Costumer struct {
	gorm.Model
	ManagerID uint
	Phone     string
	Name      string
	Fields    string
}

// Manager describes manager
type Manager struct {
	gorm.Model
	Username     string
	Password     string
	Active       bool `gorm:"default:'1'"`
	Name         string
	Phone        string
	LinkTemplate string `gorm:"default:''"`
	Greeting     string `gorm:"default:''"`
	QAs          []QA
	Costumers    []Costumer
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

func main() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		"admin",
		"952368741",
		"wazap.cvtuyrclurh0.ap-south-1.rds.amazonaws.com:3306",
		"gym",
	))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database %v", err))
	}

	db.AutoMigrate(&QA{}, &Costumer{}, &Manager{})

	r := gin.Default()
	r.Use(corsMiddleware())

	az := r.Group("/:user/:pwd")
	{
		az.Use(authMiddleware(db))

		az.POST("/check", check(db))

		az.GET("/general", getGeneral(db))
		az.POST("/general", patchGeneral(db))

		az.GET("/qa", getQA(db))
		az.POST("/qa/create", createQA(db))
		az.POST("/qa/modify", patchQA(db))
		az.POST("/qa/remove", removeQA(db))

		az.GET("/files", listFiles)
		az.POST("/upload", uploadFile)
	}

	r.Run("0.0.0.0:8090")
}
