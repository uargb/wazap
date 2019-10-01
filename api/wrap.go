package main

import "github.com/gin-gonic/gin"

func check(c *gin.Context) {
	var manager Manager
	db.Where("username = ? and password = ?",
		c.Param("user"),
		c.Param("pwd"),
	).First(&manager)

	if manager.ID <= 0 {
		c.JSON(200, gin.H{
			"ok":      false,
			"message": "Неверное имя пользователя или пароль",
		})
		return
	}

	if manager.Active {
		c.JSON(200, gin.H{
			"ok":      false,
			"message": "Пользователь деактивирован. Обратитесь к администратору",
		})
	} else {
		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}

/*func wrapGetGeneral(user, pwd string) apiResult {
	ok, data := dbGetGeneral(user, pwd)

	if !ok {
		return apiResult{Ok: false, Message: "Ошибка при получении основных данных"}
	}

	result, err := json.Marshal(&map[string]string{
		"name":         string(data["name"].([]byte)),
		"linkTemplate": string(data["link_template"].([]byte)),
		"greeting":     string(data["greeting"].([]byte)),
		"link": fmt.Sprintf(
			"https://wa.me/%s?text=%s",
			string(data["bot_phone"].([]byte)),
			url.QueryEscape(
				strings.ReplaceAll(
					string(data["link_template"].([]byte)),
					"{name}",
					string(data["name"].([]byte)),
				),
			),
		),
	})

	if err != nil {
		log.Printf("while encoding general data: %v\n", err)
		return apiResult{Ok: false, Message: "Ошибка при обработке основных данных"}
	}

	return apiResult{Ok: true, Data: string(result)}
}

func wrapUpdateGeneral(user, pwd, name, linkTmpl, greeting string) apiResult {
	if dbUpdateGeneral(user, pwd, name, linkTmpl, greeting) {
		return apiResult{Ok: true}
	}

	return apiResult{Ok: false, Message: "При обновлении данных произошла внутренняя ошибка API"}
}

func wrapGetQA(user, pwd string) apiResult {
	ok, qa := dbGetQA(user, pwd)
	if !ok {
		return apiResult{Ok: false, Message: "Ошибка при получении меню"}
	}

	qaData, err := json.Marshal(qa)
	if err != nil {
		return apiResult{Ok: false, Message: "Ошибка при упаковке меню"}
	}

	return apiResult{Ok: true, Data: string(qaData)}
}

func wrapUpdateQa(user, pwd, index, q, dsc, show, text string, image, video, attachment *multipart.FileHeader) apiResult {
	realIndex, err := strconv.Atoi(index)
	if err != nil {
		log.Printf("while converting QA index to int: %v\n", err)
		return apiResult{Ok: false, Message: "Неверный индекс запроса"}
	}

	realShow, err := strconv.Atoi(show)
	if err != nil {
		log.Printf("while converting QA show flag to int: %v\n", err)
		return apiResult{Ok: false, Message: "Неверный показатель флага отображения"}
	}

	ok := dbUpdateQa(user, pwd, realIndex, q, dsc, realShow, text)
	if !ok {
		return apiResult{Ok: false, Message: "Ошибка при обновлении меню"}
	}

	if image != nil {
		temp := strings.Split(image.Filename, ".")
		if temp[len(temp)-1] != "png" &&
			temp[len(temp)-1] != "jpg" &&
			temp[len(temp)-1] != "jpeg" {
			return apiResult{Ok: false, Message: "Неподдерживаемый формат файла"}
		}

		file, err := image.Open()
		if err == nil {
			f, err := os.OpenFile("./public/"+image.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err == nil {
				io.Copy(f, file)
				dbAttachQa(user, pwd, realIndex, "image", image.Filename)
			} else {
				log.Printf("while creating file to save: %v\n", err)
			}
		} else {
			log.Printf("while opening image: %v\n", err)
		}
	}

	if video != nil {
		temp := strings.Split(video.Filename, ".")
		if temp[len(temp)-1] != "mp4" {
			return apiResult{Ok: false, Message: "Неподдерживаемый формат файла"}
		}

		file, err := video.Open()
		if err == nil {
			f, err := os.OpenFile("./public/"+video.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err == nil {
				io.Copy(f, file)
				dbAttachQa(user, pwd, realIndex, "video", video.Filename)
			} else {
				log.Printf("while creating file to save: %v\n", err)
			}
		} else {
			log.Printf("while opening video: %v\n", err)
		}
	}

	if attachment != nil {
		temp := strings.Split(attachment.Filename, ".")
		if temp[len(temp)-1] != "pdf" {
			return apiResult{Ok: false, Message: "Неподдерживаемый формат файла"}
		}

		file, err := attachment.Open()
		if err == nil {
			f, err := os.OpenFile("./public/"+attachment.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err == nil {
				io.Copy(f, file)
				dbAttachQa(user, pwd, realIndex, "attachment", attachment.Filename)
			} else {
				log.Printf("while creating file to save: %v\n", err)
			}
		} else {
			log.Printf("while opening attachment: %v\n", err)
		}
	}

	return apiResult{Ok: true}
}

func wrapRemoveQA(user, pwd, index string) apiResult {
	realIndex, err := strconv.Atoi(index)
	if err != nil {
		log.Printf("while converting QA index to int: %v\n", err)
		return apiResult{Ok: false, Message: "Неверный индекс запроса"}
	}

	ok := dbRemoveQa(user, pwd, realIndex)
	if !ok {
		return apiResult{Ok: false, Message: "Ошибка при обновлении меню"}
	}
	return apiResult{Ok: true}
}

func wrapGetCostumers(user, pwd string) apiResult {
	data := dbGetCostumers(user, pwd)
	if data == nil {
		return apiResult{Ok: false, Message: "Ошибка при получении клиентов"}
	}

	response, err := json.Marshal(data)
	if err != nil {
		return apiResult{Ok: false, Message: "Ошибка при формировании списка клиентов"}
	}

	return apiResult{Ok: true, Data: string(response)}
}
*/
