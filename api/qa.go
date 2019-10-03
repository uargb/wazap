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

		manager := iManager.(Manager)
		var qas []QA
		db.Model(&manager).Related(&qas)

		if id := c.Param("id"); id != "/" {
			pk, _ := strconv.Atoi(id[1:])
			for _, qa := range qas {
				if qa.ID == uint(pk) {
					c.AsciiJSON(200, gin.H{
						"ok":   true,
						"data": qa,
					})
					return
				}
			}
		} else {
			c.AsciiJSON(200, gin.H{
				"ok":   true,
				"data": qas,
			})
			return
		}
	}
}

func patchQA(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}

		manager := iManager.(Manager)
		id, _ := strconv.Atoi(c.Param("id"))

		var qas []QA
		db.Model(&manager).Related(&qas)

		for _, qa := range qas {
			if qa.ID == uint(id) {
				query, exist := c.GetPostForm("query")
				if exist {
					qa.Query = query
				}

				text, exist := c.GetPostForm("text")
				if exist {
					qa.Text = text
				}

				image, exist := c.GetPostForm("image")
				if exist {
					qa.Image = image
				}

				video, exist := c.GetPostForm("video")
				if exist {
					qa.Video = video
				}

				attachment, exist := c.GetPostForm("attachment")
				if exist {
					qa.Attachment = attachment
				}

				next, exist := c.GetPostForm("next")
				if exist {
					qa.Next = next
				}

				write, exist := c.GetPostForm("write")
				if exist {
					qa.Write = write
				}
			}
		}

		c.AsciiJSON(200, gin.H{
			"ok": true,
		})
	}
}
