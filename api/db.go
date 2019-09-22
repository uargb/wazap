package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var dbConn *sqlx.DB

func dbInit(addr, user, pwd, db string) {
	var err error
	dbConn, err = sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s", user, pwd, addr, db))
	if err != nil {
		log.Fatalf("while creating connetcion to db: %s", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("while connecting to db: %v", err)
	}
}

func dbLogin(user, pwd string) int {
	result := make(map[string]interface{})
	err := dbConn.QueryRowx(
		"select id, active from managers where username = ? and password = ?",
		user, pwd).MapScan(result)

	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Printf("while authenticating manager: %v\n", err)
		return -1
	}

	if _, ok := result["id"]; ok && result["id"].(int64) > 0 {
		if result["active"].(int64) == 1 {
			return 0
		}

		return -2
	}

	return -3
}

func dbGetGeneral(user, pwd string) (bool, map[string]interface{}) {
	result := make(map[string]interface{})
	err := dbConn.QueryRowx(
		"select name, link_template, greeting, (select bot_phone from managers where id = 1) as bot_phone from managers where username = ? and password = ?",
		user, pwd).MapScan(result)
	if err != nil {
		log.Printf("while getting general info: %v\n", err)
		return false, nil
	}
	return true, result
}

func dbUpdateGeneral(user, pwd, name, linkTmpl, greeting string) bool {
	_, err := dbConn.Exec(
		"update managers set name = ?, link_template = ?, greeting = ? where username = ? and password = ?",
		name, linkTmpl, greeting, user, pwd,
	)

	if err != nil {
		log.Printf("while updating general info: %v\n", err)
		return false
	}

	return true
}

func dbGetQA(user, pwd string) (bool, []byte) {
	result := make(map[string]interface{})
	err := dbConn.QueryRowx(
		"select qa from managers where username = ? and password = ?",
		user, pwd).MapScan(result)
	if err != nil {
		log.Printf("while getting QA: %v\n", err)
		return false, nil
	}

	return true, result["qa"].([]byte)
}

func dbUpdateQa(user, pwd string, index int, q, dsc, text string) bool {
	result := make(map[string]interface{})
	err := dbConn.QueryRowx(
		"select qa from managers where username = ? and password = ?",
		user, pwd).MapScan(result)
	if err != nil {
		log.Printf("while getting QA: %v\n", err)
		return false
	}

	var qa []interface{}
	json.Unmarshal(result["qa"].([]byte), &qa)

	if d := index - (len(qa) - 1); d > 0 {
		for i := 0; i < d; i++ {
			qa = append(qa, make(map[string]interface{}))
		}
	}

	node := qa[index].(map[string]interface{})
	node["query"] = q
	node["description"] = dsc
	node["text"] = text
	qa[index] = node

	newQa, err := json.Marshal(&qa)
	if err != nil {
		log.Printf("while encoding new QA: %v\n", err)
		return false
	}

	_, err = dbConn.Exec(
		"update managers set qa = ? where username = ? and password = ?",
		string(newQa), user, pwd,
	)

	if err != nil {
		log.Printf("while updating QA: %v\n", err)
		return false
	}

	return true
}

func dbRemoveQa(user, pwd string, index int) bool {
	result := make(map[string]interface{})
	err := dbConn.QueryRowx(
		"select qa from managers where username = ? and password = ?",
		user, pwd).MapScan(result)
	if err != nil {
		log.Printf("while getting QA: %v\n", err)
		return false
	}

	var qa []interface{}
	json.Unmarshal(result["qa"].([]byte), &qa)

	if index > len(qa)+1 {
		log.Printf("while removing QA item: trying to remove to big index\n")
		return false
	}

	modified := qa[:index]
	if index < len(qa) {
		modified = append(modified, qa[index+1:]...)
	}

	newQa, err := json.Marshal(&modified)
	if err != nil {
		log.Printf("while encoding new QA: %v\n", err)
		return false
	}

	_, err = dbConn.Exec(
		"update managers set qa = ? where username = ? and password = ?",
		string(newQa), user, pwd,
	)

	if err != nil {
		log.Printf("while updating QA: %v\n", err)
		return false
	}

	return true
}
