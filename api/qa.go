package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getQA(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}

		manager := iManager.(*Manager)
		var qas []QA
		db.Model(&manager).Related(&qas)

		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			id = -1
		}

		if id != -1 {
			for _, qa := range qas {
				if qa.ID == uint(id) {
					c.JSON(200, gin.H{
						"ok":   true,
						"data": qa,
					})
					return
				}
			}
		} else {
			c.JSON(200, gin.H{
				"ok":   true,
				"data": qas,
			})
			return
		}
	}
}

func createQA(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		qa := QA{ManagerID: manager.ID}
		db.Create(&qa)

		c.JSON(200, gin.H{
			"ok":   true,
			"data": qa,
		})
	}
}

func patchQA(db *gorm.DB) func(*gin.Context) {
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

		var qas []QA
		db.Model(&manager).Related(&qas)

		for _, qa := range qas {
			if qa.ID == uint(id) {
				query, exist := c.GetPostForm("Query")
				if exist {
					qa.Query = query
				}

				text, exist := c.GetPostForm("Text")
				if exist {
					qa.Text = text
				}

				image, exist := c.GetPostForm("Image")
				if exist {
					qa.Image = image
				}

				video, exist := c.GetPostForm("Video")
				if exist {
					qa.Video = video
				}

				attachment, exist := c.GetPostForm("Attachment")
				if exist {
					qa.Attachment = attachment
				}

				next, exist := c.GetPostForm("Next")
				if exist {
					qa.Next = next
				}

				write, exist := c.GetPostForm("Write")
				if exist {
					qa.Write = write
				}

				newStatus, exist := c.GetPostForm("NewStatus")
				if exist {
					qa.NewStatus = newStatus
				}

				notifyManager, exist := c.GetPostForm("NotifyManager")
				if exist {
					if notifyManager == "true" {
						qa.NotifyManager = true
					} else {
						qa.NotifyManager = false
					}

				}

				db.Save(&qa)
			}
		}

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}

func removeQA(db *gorm.DB) func(*gin.Context) {
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

		db.Where("id = ? and manager_id = ?", id, manager.ID).Delete(&QA{})

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}
