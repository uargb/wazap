package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getCostumers(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		var costumers []Costumer
		db.Model(&manager).Related(&costumers)

		for i := range costumers {
			costumers[i].Fields = strings.ReplaceAll(costumers[i].Fields, "&", "\n")
			costumers[i].Fields = strings.ReplaceAll(costumers[i].Fields, "=", ": ")
		}

		c.AsciiJSON(200, gin.H{
			"ok":   true,
			"data": costumers,
		})
	}
}

func patchCostumer(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		var costumers []Costumer
		db.Model(&manager).Related(&costumers)

		qID := c.Query("id")
		id, err := strconv.Atoi(qID)
		if err != nil {
			c.AsciiJSON(200, gin.H{
				"ok":      false,
				"message": "Неверный ID клиента",
			})
		}

		for i := range costumers {
			if costumers[i].ID == uint(id) {
				qPrice, exist := c.GetPostForm("Price")
				if exist {
					price, err := strconv.ParseFloat(qPrice, 64)
					if err != nil {
						c.AsciiJSON(200, gin.H{
							"ok":      false,
							"message": "Неверное значение суммы оплаты",
						})
					}

					costumers[i].Price = price
				}

				qPeriod, exist := c.GetPostForm("Period")
				if exist {
					period, err := strconv.Atoi(qPeriod)
					if err != nil {
						c.AsciiJSON(200, gin.H{
							"ok":      false,
							"message": "Неверное значение периода",
						})
					}

					costumers[i].Period = period
				}

				db.Save(&costumers[i])
				break
			}
		}

		c.AsciiJSON(200, gin.H{
			"ok": true,
		})
	}
}
