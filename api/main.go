package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gernest/alien"
)

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func jsonResult(result apiResult) []byte {
	str, err := json.Marshal(&result)
	if err != nil {
		log.Printf("while encoding result %v: %v\n", result, err)
	}
	return str
}

func main() {
	dbInit("wazap.cvtuyrclurh0.ap-south-1.rds.amazonaws.com:3306", "admin", "952368741", "gym")

	api := alien.New()
	api.Use(middleware)

	api.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		result := wrapLogin(r.FormValue("username"), r.FormValue("password"))
		w.Write(jsonResult(result))
	})

	api.Get("/general", func(w http.ResponseWriter, r *http.Request) {
		result := wrapGetGeneral(r.FormValue("username"), r.FormValue("password"))
		w.Write(jsonResult(result))
	})
	api.Post("/general", func(w http.ResponseWriter, r *http.Request) {
		result := wrapUpdateGeneral(
			r.FormValue("username"),
			r.FormValue("password"),
			r.FormValue("name"),
			r.FormValue("linkTemplate"),
			r.FormValue("greeting"),
		)
		w.Write(jsonResult(result))
	})

	api.Get("/qa", func(w http.ResponseWriter, r *http.Request) {
		result := wrapGetQA(r.FormValue("username"), r.FormValue("password"))
		w.Write(jsonResult(result))
	})
	api.Post("/qa", func(w http.ResponseWriter, r *http.Request) {
		_, imageHeader, _ := r.FormFile("image")
		_, videoHeader, _ := r.FormFile("video")
		_, attachmentHeader, _ := r.FormFile("attachment")

		result := wrapUpdateQa(
			r.FormValue("username"),
			r.FormValue("password"),
			r.FormValue("index"),
			r.FormValue("query"),
			r.FormValue("description"),
			r.FormValue("show"),
			r.FormValue("text"),
			imageHeader, videoHeader, attachmentHeader,
		)
		w.Write(jsonResult(result))
	})
	api.Post("/qa/remove", func(w http.ResponseWriter, r *http.Request) {
		result := wrapRemoveQA(r.FormValue("username"), r.FormValue("password"), r.FormValue("index"))
		w.Write(jsonResult(result))
	})

	log.Fatal(http.ListenAndServe(":8090", api))
}
