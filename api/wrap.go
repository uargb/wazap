package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
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

func wrapUpdateQa(user, pwd, index, q, dsc, text string, image, video, attachment *multipart.FileHeader) apiResult {
	realIndex, err := strconv.Atoi(index)
	if err != nil {
		log.Printf("while converting QA index to int: %v\n", err)
		return apiResult{Ok: false, Message: "Неверный индекс запроса"}
	}

	ok := dbUpdateQa(user, pwd, realIndex, q, dsc, text)
	if !ok {
		return apiResult{Ok: false, Message: "Ошибка при обновлении меню"}
	}

	if image != nil {

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
