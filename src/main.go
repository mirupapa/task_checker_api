package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/task_checker_api/src/auth"
	"github.com/task_checker_api/src/controller"
	"github.com/urfave/negroni"
)

func main() {
	err := godotenv.Load(fmt.Sprintf("../%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		print("error_env")
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	task := r.PathPrefix("/task").Subrouter()
	{
		task.Handle("/", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.GetTasks),
		))
	}

	// cors
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET"})

	//サーバー起動
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
