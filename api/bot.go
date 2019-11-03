package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func botMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var costumer Costumer
		db.Where("phone = ?", c.Param("phone")).FirstOrCreate(&costumer)
		c.Set("costumer", &costumer)

		message := strings.TrimSpace(c.Query("message"))
		c.Set("message", message)

		if costumer.ManagerID == 0 {
			var managers []Manager
			db.Find(&managers)

			costumer.Phone = c.Param("phone")

			for _, manager := range managers {
				if strings.Contains(message, manager.Name) {
					costumer.ManagerID = manager.ID
					break
				}
			}

			if costumer.ManagerID == 0 {
				costumer.ManagerID = 1
			}

			db.Save(&costumer)

			var manager Manager
			db.Where("ID = ?", costumer.ManagerID).First(&manager)

			c.JSON(200, gin.H{
				"ok":  true,
				"did": "registered",
				"data": []gin.H{
					{
						"Text": manager.Greeting,
					},
				},
			})
		} else {
			var manager Manager
			db.Where("ID = ?", costumer.ManagerID).First(&manager)
			c.Set("manager", &manager)
		}
	}
}

func botGetAnswer(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iCostumer, _ := c.Get("costumer")
		costumer := iCostumer.(*Costumer)

		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		iMessage, _ := c.Get("message")
		message := iMessage.(string)

		var qas []QA
		db.Model(&manager).Related(&qas)

		data := make([]interface{}, 0)

		if len(costumer.Write) > 0 {
			costumer.Next = 0
			db.Save(costumer)
		}

		if costumer.Next > 0 {
			for _, qa := range qas {
				if qa.ID == costumer.Next {
					data = append(data, qa)

					if len(qa.Write) > 0 {
						value := qa.Query
						if len(value) == 0 {
							value = message
						}

						if len(costumer.Fields) > 0 {
							costumer.Fields = strings.ReplaceAll(costumer.Fields, qa.Write, ":old:"+qa.Write)
							costumer.Fields += fmt.Sprintf("&%s=%s", qa.Write, value)
						} else {
							costumer.Fields += fmt.Sprintf("%s=%s", qa.Write, value)
						}
					}

					if len(qa.NewStatus) > 0 {
						costumer.Status = qa.NewStatus
						db.Save(costumer)
					}

					break
				}
			}

			costumer.Next = 0
			db.Save(costumer)
		} else {
			for _, qa := range qas {
				if strings.ToLower(qa.Query) == strings.ToLower(message) {
					data = append(data, qa)

					if len(qa.Next) > 0 {
						next, _ := strconv.Atoi(qa.Next)
						costumer.Next = uint(next)
						db.Save(costumer)
					}
				}
			}
		}

		if len(data) == 0 {
			for _, qa := range qas {
				if qa.Query == "<неопознанное>" {
					data = append(data, qa)
				}
			}

		}

		c.JSON(200, gin.H{
			"ok":       true,
			"data":     data,
			"manager":  manager,
			"costumer": costumer,
		})
	}
}

func botUpdateCostumerName(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iCostumer, _ := c.Get("costumer")
		costumer := iCostumer.(*Costumer)

		costumer.Name = c.Query("name")
		db.Save(costumer)

		c.JSON(200, gin.H{
			"ok":  true,
			"did": "renamed",
		})
	}
}

func botGetMailing(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var mailing Mailing
		db.Where("done = 0").First(&mailing)

		if mailing.ID == 0 {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Рассылок нет",
			})

			return
		}

		var manager Manager
		db.First(&manager, mailing.ManagerID)

		var card QA
		db.Where("id = ? and manager_id = ?", mailing.CardID, mailing.ManagerID).First(&card)

		var costumers []Costumer
		db.Model(&manager).Related(&costumers)

		phones := make([]string, 0)
		for _, costumer := range costumers {
			if mailing.Status == "all" || costumer.Status == mailing.Status {
				phones = append(phones, costumer.Phone)
			}
		}

		c.JSON(200, gin.H{
			"ok": true,
			"data": gin.H{
				"card":   card,
				"phones": phones,
			},
		})

		mailing.Done = true
		db.Save(&mailing)
	}
}
