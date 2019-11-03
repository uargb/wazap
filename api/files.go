package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
)

func processCSVImport(db *gorm.DB, manager *Manager, filename string) error {
	file, err := os.Open(fmt.Sprintf("./public/%d-%s", manager.ID, filename))
	if err != nil {
		return err
	}

	data, err := gocsv.CSVToMaps(file)
	if err != nil {
		return err
	}

	for _, c := range data {
		db.Create(&Costumer{
			ManagerID: manager.ID,
			Phone:     c["Phone"],
			Name:      c["Name"],
			Status:    c["Status"],
		})
	}

	return nil
}

func uploadFile(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		iManager, ok := c.Get("manager")
		if !ok {
			return
		}
		manager := iManager.(*Manager)

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(200, gin.H{
				"ok": false,
			})
		}

		c.SaveUploadedFile(file, fmt.Sprintf("./public/%d-%s", manager.ID, file.Filename))
		if strings.HasSuffix(file.Filename, ".csv") {
			err = processCSVImport(db, manager, file.Filename)
			if err != nil {
				c.JSON(200, gin.H{
					"ok": false,
				})
				return
			}

			os.Remove(fmt.Sprintf("./public/%d-%s", manager.ID, file.Filename))
		}

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}

func listFiles(c *gin.Context) {
	iManager, ok := c.Get("manager")
	if !ok {
		return
	}
	manager := iManager.(*Manager)

	ext := c.DefaultQuery("ext", "*")

	files, _ := filepath.Glob(fmt.Sprintf("./public/%d-*.%s", manager.ID, ext))
	for i := range files {
		files[i] = strings.TrimPrefix(files[i], fmt.Sprintf("public/%d-", manager.ID))
		files[i] = strings.TrimPrefix(files[i], fmt.Sprintf("public\\%d-", manager.ID))
	}

	c.JSON(200, gin.H{
		"ok":   true,
		"data": files,
	})
}
