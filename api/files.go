package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	iManager, ok := c.Get("manager")
	if !ok {
		return
	}
	manager := iManager.(*Manager)

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
	manager := iManager.(*Manager)

	ext := c.DefaultQuery("ext", "*")

	files, _ := filepath.Glob(fmt.Sprintf("./public/%d-*.%s", manager.ID, ext))
	for i := range files {
		files[i] = strings.TrimPrefix(files[i], fmt.Sprintf("public/%d-", manager.ID))
		files[i] = strings.TrimPrefix(files[i], fmt.Sprintf("public\\%d-", manager.ID))
	}

	c.AsciiJSON(200, gin.H{
		"ok":   true,
		"data": files,
	})
}
