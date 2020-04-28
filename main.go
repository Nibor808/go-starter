package main

import (
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go-starter/controller"
	"go-starter/utils"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file found")
	}

	tpl = template.Must(template.ParseFiles("view/index.html"))
}

func main() {
	r := httprouter.New()
	db := utils.GetMongoSession()
	uc := controller.NewUserController(db)
	ac := controller.NewAuthController(db)

	r.GET("/", index)

	/* AUTH */
	r.POST("/signupemail", ac.SignUpEmail)
	r.GET("/confirmemail/:token/:userId", ac.ConfirmEmail)
	r.POST("/signuppassword", ac.SignUpPassword)
	r.POST("/signin", ac.SignIn)
	r.GET("/signout", ac.SignOut)

	/* USER */
	r.GET("/users", controller.CheckSession(uc.AllUsers, db))
	r.GET("/user", controller.CheckSession(uc.User, db))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html charset=utf8")
	w.WriteHeader(200)

	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
