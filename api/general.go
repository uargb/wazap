package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getGeneral(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(Manager)

		var admin Manager
		db.First(&admin)

		c.AsciiJSON(200, gin.H{
			"ok": true,
			"data": gin.H{
				"name":         manager.Name,
				"linkTemplate": manager.LinkTemplate,
				"greeting":     manager.Greeting,
				"link": fmt.Sprintf(
					"https://wa.me/%s/?text=%s",
					admin.Phone,
					strings.ReplaceAll(manager.LinkTemplate, "{name}", manager.Name),
				),
			},
		})
	}
}

func patchGeneral(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}

		manager := iManager.(Manager)

		linkTemplate, exist := c.GetPostForm("linkTemplate")
		if exist {
			manager.LinkTemplate = linkTemplate
		}

		greeting, exist := c.GetPostForm("greeting")
		if exist {
			manager.Greeting = greeting
		}

		db.Save(&manager)

		c.AsciiJSON(200, gin.H{
			"ok": true,
		})
	}
}
