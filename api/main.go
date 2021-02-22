package main

import (
	"fmt"
	"go-starter/controller"
	"go-starter/middleware"
	"go-starter/utils"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	db := utils.GetMongoSession()
	uc := controller.NewUserController(db)
	ac := controller.NewAuthController(db)
	dc := controller.NewDataController(db)

	r.GET("/wsconnect", handleWebSocket)

	/* AUTH */
	r.POST("/signupemail", ac.SignUpEmail)
	r.GET("/confirmemaildata/:token/:userID", ac.ConfirmEmail)
	r.POST("/signuppassword", middleware.CheckSession(ac.SignUpPassword, db))
	r.POST("/signin", ac.SignIn)
	r.GET("/signout", ac.SignOut)

	/* USER */
	r.GET("/users", middleware.CheckSession(uc.AllUsers, db))
	r.GET("/user", middleware.CheckSession(uc.User, db))

	/* DATA */
	r.GET("/alldata", middleware.CheckSession(dc.AllData, db))
	r.POST("/savedata", middleware.CheckSession(dc.SaveData, db))
	r.POST("/updatedata", middleware.CheckSession(dc.UpdateData, db))

	// DEPLOY_MODE is defined in docker-compose.yml
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

func handleWebSocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var upgrader = websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		if r.Header.Get("Origin") == "http://localhost:3000" {
			return true
		}

		return false
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR: ", err)
			break
		}

		err = conn.WriteMessage(mt, []byte("Server here. Message {"+string(message)+"} received!"))
		if err != nil {
			fmt.Println("ERROR: ", err)
			break
		}
	}
}
