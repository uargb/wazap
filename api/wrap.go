package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type apiResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func wrapLogin(user, pwd string) apiResult {
	code := dbLogin(user, pwd)

	var result apiResult

	switch code {
	case -1:
		result.Ok = false
		result.Message = "Ошибка при работе с БД"
		break
	case -2:
		result.Ok = false
		result.Message = "Пользователь деактивирован. Обратитесь к администратору"
		break
	case -3:
		result.Ok = false
		result.Message = "Неверное имя пользователя или пароль"
		break
	case 0:
		result.Ok = true
		break
	}

	return result
}

func wrapGetGeneral(user, pwd string) apiResult {
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

	return apiResult{Ok: true, Data: string(qa)}
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
