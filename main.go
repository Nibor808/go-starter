package main

import (
	"go-starter/controller"
	"go-starter/middleware"
	"go-starter/utils"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No env file found")
	}

	tpl = template.Must(template.ParseFiles("./view/index.html"))
}

func main() {
	r := httprouter.New()
	db := utils.GetMongoSession()
	uc := controller.NewUserController(db)
	ac := controller.NewAuthController(db)
	dc := controller.NewDataController(db)

	r.GET("/", index)

	/* AUTH */
	r.POST("/signupemail", ac.SignUpEmail)
	r.GET("/confirmemail/:token/:userID", ac.ConfirmEmail)
	r.POST("/signuppassword", ac.SignUpPassword)
	r.POST("/signin", ac.SignIn)
	r.GET("/signout", ac.SignOut)

	/* USER */
	r.GET("/users", middleware.CheckSession(uc.AllUsers, db))
	r.GET("/user", middleware.CheckSession(uc.User, db))

	/* DATA */
	r.GET("/alldata", middleware.CheckSession(dc.AllData, db))
	r.POST("/savedata", middleware.CheckSession(dc.SaveData, db))
	r.POST("/updatedata", middleware.CheckSession(dc.UpdateData, db))

	mode, dExists := os.LookupEnv("DEPLOY_MODE")
	if !dExists {
		log.Println("Cannot get DEPLOY_MODE from .env")
	}

	var handler http.Handler

	if mode == "development" {
		handler = &middleware.Logger{Handler: r}
	} else {
		handler = r
	}

	log.Println("Listening on 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html charset=utf8")
	w.WriteHeader(http.StatusOK)

	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
