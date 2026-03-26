package main

import (
	"log"
	"net/http"
	"os"

	"fabel/app"
	"fabel/db"
	"fabel/llm"
	"fabel/web"
)

func main() {
	database, err := db.Open("fabel.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	llmClient := llm.New(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_BASE_URL"))

	a := app.New(database, llmClient, web.Dist)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Routes()))
}
