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

	r.GET("/", index)
	r.GET("/users", controller.CheckSession(uc.AllUsers, db))
	r.GET("/user", controller.CheckSession(uc.User, db))
	r.POST("/signup", uc.SignUp)
	r.GET("/verifyemail/:token/:userId", uc.ConfirmVerificationEmail)
	r.POST("/collectpassword", uc.CollectPassword)
	r.POST("/signin", uc.SignIn)
	r.GET("/signout", uc.SignOut)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html charset=utf8")
	w.WriteHeader(200)

	err := tpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
