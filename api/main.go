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
	ManagerID     uint
	Query         string
	Text          string
	Image         string
	Video         string
	Attachment    string
	Next          string
	Write         string
	NewStatus     string
	NotifyManager bool `gorm:"default:0"`
}

// Costumer describes costumer
type Costumer struct {
	gorm.Model
	ManagerID uint
	Phone     string
	Name      string
	Price     float64
	Period    int `gorm:"default:0"`
	Fields    string
	Next      uint
	Status    string
	Write     string
}

// Manager describes manager
type Manager struct {
	gorm.Model
	Username     string
	Password     string
	Active       bool   `gorm:"default:'1'"`
	Name         string `gorm:"default:''"`
	Phone        string `gorm:"default:''"`
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

	az := r.Group("/admin/:user/:pwd")
	az.Use(authMiddleware(db))
	{
		az.POST("/check", check(db))

		az.GET("/files", listFiles)
		az.POST("/upload", uploadFile)

		az.GET("/general", getGeneral(db))
		az.POST("/general", patchGeneral(db))

		az.GET("/qa", getQA(db))
		az.POST("/qa/create", createQA(db))
		az.POST("/qa/modify", patchQA(db))
		az.POST("/qa/remove", removeQA(db))

		az.GET("/costumers", getCostumers(db))
		az.POST("/costumers/modify", patchCostumer(db))
	}

	bot := r.Group("/bot/:phone")
	bot.Use(botMiddleware(db))
	{
		bot.GET("/answer", botGetAnswer(db))
		bot.GET("/rename", botUpdateCostumerName(db))
	}

	r.Run("0.0.0.0:8090")
}
