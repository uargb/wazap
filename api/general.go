package main

import (
	"fmt"
	"net/url"
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
		manager := iManager.(*Manager)

		var admin Manager
		db.First(&admin)

		c.JSON(200, gin.H{
			"ok": true,
			"data": gin.H{
				"Name":         manager.Name,
				"LinkTemplate": manager.LinkTemplate,
				"Greeting":     manager.Greeting,
				"Phone":        manager.Phone,
				"Link": fmt.Sprintf(
					"https://wa.me/%s/?text=%s",
					admin.Phone,
					url.QueryEscape(strings.ReplaceAll(manager.LinkTemplate, "{name}", manager.Name)),
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

		manager := iManager.(*Manager)

		name, exist := c.GetPostForm("Name")
		if exist {
			manager.Name = name
		}

		linkTemplate, exist := c.GetPostForm("LinkTemplate")
		if exist {
			manager.LinkTemplate = linkTemplate
		}

		greeting, exist := c.GetPostForm("Greeting")
		if exist {
			manager.Greeting = greeting
		}

		phone, exist := c.GetPostForm("Phone")
		if exist {
			manager.Phone = phone
		}

		db.Save(&manager)

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}
