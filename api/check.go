package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func authMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var manager Manager
		db.Where("username = ?",
			c.Param("user"),
		).First(&manager)

		if manager.Password != c.Param("pwd") {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Неверное имя пользователя или пароль",
			})
			return
		}

		if !manager.Active {
			c.JSON(200, gin.H{
				"ok":      false,
				"message": "Пользователь деактивирован. Обратитесь к администратору",
			})
			return
		}

		c.Set("manager", manager)
	}
}

func check(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		_, ok := c.Get("manager")
		if !ok {
			return
		}

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}
