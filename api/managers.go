package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getManagers(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}

		manager := iManager.(*Manager)

		if manager.Username != "admin" {
			c.JSON(200, gin.H{
				"ok":      "false",
				"message": "Разрешено только для администратора",
			})

			return
		}

		var managers []Manager
		db.Select("id, name").Find(&managers)

		c.JSON(200, gin.H{
			"ok":   true,
			"data": managers,
		})
	}
}

func getStats(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}

		manager := iManager.(*Manager)

		if manager.Username != "admin" {
			c.JSON(200, gin.H{
				"ok":      "false",
				"message": "Разрешено только для администратора",
			})

			return
		}

		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.JSON(200, gin.H{
				"ok":      "false",
				"message": "Невалидный ID",
			})

			return
		}

		var victim Manager
		db.First(&victim, id)

		var costumers []Costumer
		db.Model(&victim).Related(&costumers)

		stats := make(map[string]int)
		for _, costumer := range costumers {
			if _, ok := stats[costumer.Status]; !ok {
				stats[costumer.Status] = 0
			}

			stats[costumer.Status] += 1
		}

		i := 0
		keys := make([]string, len(stats))
		for k := range stats {
			if k == "" {
				k = "Без статуса"
			}

			keys[i] = k
			i++
		}

		c.JSON(200, gin.H{
			"ok": true,
			"data": gin.H{
				"stats": stats,
			},
		})
	}
}
