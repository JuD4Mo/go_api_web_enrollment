package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
	"github.com/JuD4Mo/go_api_web_enrollment/pkg/bootstrap"
	"github.com/JuD4Mo/go_api_web_enrollment/pkg/handler"

	courseSDK "github.com/JuD4Mo/go_api_web_sdk/course"
	userSDK "github.com/JuD4Mo/go_api_web_sdk/user"
	"github.com/joho/godotenv"
)

func main() {

	//Cargamos las variables de entorno que están en el archivo .env por medio del package godotenv
	_ = godotenv.Load()

	//Instanciamos un logger propio
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}
	token := os.Getenv("API_COURSE_TOKEN")

	userTransport := userSDK.NewHttpClient(os.Getenv("API_USER_URL"), "")
	courseTransport := courseSDK.NewHttpClient(os.Getenv("API_COURSE_URL"), token)

	ctx := context.Background()

	enrollRepo := enrollment.NewRepo(db, l)
	enrollService := enrollment.NewService(l, enrollRepo, userTransport, courseTransport)
	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollService, enrollment.Config{LimitPage: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	//Se crea una instancia de un servidor
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST, PATCH, OPTIONS, DELETE, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
