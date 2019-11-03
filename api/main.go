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
	Text          string `sql:"type:text"`
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
	Period    int    `gorm:"default:0"`
	Fields    string `sql:"type:text"`
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
	Greeting     string `gorm:"default:''"; sql:"type:text"`
	QAs          []QA
	Costumers    []Costumer
}

type Mailing struct {
	gorm.Model
	ManagerID uint
	Status    string
	CardID    int
	Done      bool `gorm:"default:'0'"`
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

func main() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		"wazap",
		"AwyYZZdZhcPSNXAk",
		"13.232.119.162:3306",
		"wazap",
	))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database %v", err))
	}

	db.AutoMigrate(&QA{}, &Costumer{}, &Manager{}, &Mailing{})
	db.Exec("SET NAMES utf8mb4")

	r := gin.Default()
	r.Use(corsMiddleware())

	az := r.Group("/admin/:user/:pwd")
	az.Use(authMiddleware(db))
	{
		az.POST("/check", check(db))

		az.GET("/files", listFiles)
		az.POST("/upload", uploadFile(db))

		az.GET("/general", getGeneral(db))
		az.POST("/general", patchGeneral(db))

		az.GET("/qa", getQA(db))
		az.POST("/qa/create", createQA(db))
		az.POST("/qa/modify", patchQA(db))
		az.POST("/qa/remove", removeQA(db))

		az.GET("/costumers", getCostumers(db))
		az.POST("/costumers/modify", patchCostumer(db))
		az.POST("/costumers/remove", removeCostumer(db))
		az.GET("/costumers/export", exportCostumers(db))
		az.POST("/costumers/send", mailing(db))

		az.GET("/managers", getManagers(db))
		az.GET("/stats", getStats(db))
	}

	bot := r.Group("/bot/:phone")
	bot.Use(botMiddleware(db))
	{
		bot.GET("/answer", botGetAnswer(db))
		bot.GET("/rename", botUpdateCostumerName(db))
		bot.GET("/mailing", botGetMailing(db))
	}

	r.Run("0.0.0.0:8090")
}
