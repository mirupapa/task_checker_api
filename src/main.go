package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"main/src/auth"
	"main/src/controller"

	"cloud.google.com/go/logging"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

// ENVLoad Env load
func ENVLoad() {
	env := os.Getenv("ENV")
	if env == "development" {
		err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
		if err != nil {
			print("error_env")
		}
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			fmt.Println(pair[0] + ":" + pair[1])
		}
	}
}

// Logging Logging
func Logging() {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the log to write to.
	logName := "task-checker-log"

	logger := client.Logger(logName).StandardLogger(logging.Info)

	// Logs "hello world", log entry is visible at
	// Cloud Logs.
	logger.Println("logging start")
}

func main() {
	ENVLoad()
	// Logging()
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0] + ":" + pair[1])
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()
	env := os.Getenv("ENV")
	if env != "development" {
		r.Schemes("https")
	}
	// ログイン
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	// サインアップ
	r.HandleFunc("/signUp", auth.SignUpHandler).Methods("POST")
	// 承認
	r.Handle("/auth", negroni.New(
		negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
		negroni.Wrap(auth.ExportUserInfo),
	))

	task := r.PathPrefix("/task").Subrouter()
	{
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.GetTasksHandler),
		)).Methods("GET")
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.PostTask),
		)).Methods("POST")
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.PutTask),
		)).Methods("PUT")
		task.Handle("", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.DeleteTask),
		)).Methods("DELETE")
		task.Handle("/done", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.PutDone),
		)).Methods("PUT")
		task.Handle("/upSort", negroni.New(
			negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext),
			negroni.Wrap(controller.UpSort),
		)).Methods("PUT")
	}

	// cors
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "PUT", "DELETE"})

	//サーバー起動
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}
