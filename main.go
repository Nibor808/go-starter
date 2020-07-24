package main

import (
	"go-starter/controller"
	"go-starter/utils"
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file found")
	}

	tpl = template.Must(template.ParseFiles("./view/index.html"))
}

func main() {
	r := httprouter.New()
	db := utils.GetMongoSession()
	uc := controller.NewUserController(db)
	ac := controller.NewAuthController(db)
	fc := controller.NewFormController(db)

	r.GET("/", index)

	/* AUTH */
	r.POST("/signupemail", ac.SignUpEmail)
	r.GET("/confirmemail/:token/:userID", ac.ConfirmEmail)
	r.POST("/signuppassword", ac.SignUpPassword)
	r.POST("/signin", ac.SignIn)
	r.GET("/signout", ac.SignOut)

	/* USER */
	r.GET("/users", controller.CheckSession(uc.AllUsers, db))
	r.GET("/user", controller.CheckSession(uc.User, db))

	/* FORM */
	r.POST("/saveform", fc.SaveForm)
	r.POST("/updateform", fc.UpdateForm)

	log.Println("Listening on 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html charset=utf8")
	w.WriteHeader(200)

	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
