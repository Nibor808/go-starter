package main

import (
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go-starter/controller"
	"log"
	"net/http"
	"scheduler_backend/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file found")
	}
}

func main() {
	r := httprouter.New()
	db := utils.GetMongoSession()
	uc := controller.NewUserController(db)

	r.GET("/", index)

	// users
	r.GET("/users", uc.AllUsers)
	r.GET("/user/:id", uc.User)
	r.POST("/user/:id", uc.AddUser)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html charset=utf8")
	w.WriteHeader(200)

	_, err := w.Write([]byte(`<h1>Hello There!</h1>`))
	if err != nil {
		log.Fatalln("Index Response Error:", err)
	}
}

