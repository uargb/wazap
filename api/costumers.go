package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
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

		c.JSON(200, gin.H{
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
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Неверный ID клиента",
			})
		}

		for i := range costumers {
			if costumers[i].ID == uint(id) {
				status, exist := c.GetPostForm("Status")
				if exist {
					costumers[i].Status = status
				}

				qPrice, exist := c.GetPostForm("Price")
				if exist {
					price, err := strconv.ParseFloat(qPrice, 64)
					if err != nil {
						c.JSON(200, gin.H{
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
						c.JSON(200, gin.H{
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

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}

func removeCostumer(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "invalid ID",
			})
		}

		db.Where("id = ? and manager_id = ?", id, manager.ID).Delete(&Costumer{})

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}

func exportCostumers(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		var costumers []Costumer
		db.Model(&manager).Related(&costumers)

		file, err := os.Create(fmt.Sprintf("public/%d-export.csv", manager.ID))
		if err != nil {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Внутренняя ошибка сервера: file creation failed",
			})

			return
		}

		gocsv.MarshalFile(&costumers, file)
		file.Close()

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="export.csv"`,
		}

		file, err = os.Open(fmt.Sprintf("public/%d-export.csv", manager.ID))
		fi, err := file.Stat()
		c.DataFromReader(200, fi.Size(), "application/octet-stream", file, extraHeaders)
		file.Close()
	}
}

func mailing(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		status, exist := c.GetPostForm("status")
		if !exist {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Заполните поле: статус",
			})
		}

		card, exist := c.GetPostForm("card")
		if !exist {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Заполните поле: ID карточки",
			})
		}

		cardId, err := strconv.Atoi(card)
		if err != nil {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Неверно заполнено поле: ID карточки",
			})
		}

		mailing := Mailing{ManagerID: manager.ID, Status: status, CardID: cardId}
		db.Create(&mailing)

		c.JSON(200, gin.H{
			"ok":      true,
			"manager": manager.ID,
			"status":  status,
			"card":    card,
		})
	}
}
