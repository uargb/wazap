package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	iManager, ok := c.Get("manager")
	if !ok {
		return
	}
	manager := iManager.(Manager)

	file, err := c.FormFile("file")
	if err != nil {
		c.AsciiJSON(200, gin.H{
			"ok": false,
		})
	}

	c.SaveUploadedFile(file, fmt.Sprintf("./public/%d-%s", manager.ID, file.Filename))
	c.AsciiJSON(200, gin.H{
		"ok": true,
	})
}

func listFiles(c *gin.Context) {
	iManager, ok := c.Get("manager")
	if !ok {
		return
	}
	manager := iManager.(Manager)

	files, _ := filepath.Glob(fmt.Sprintf("./public/%d-*", manager.ID))
	data, _ := json.Marshal(files)

	c.AsciiJSON(200, gin.H{
		"ok":   true,
		"data": string(data),
	})
}
