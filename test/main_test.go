package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/JuD4Mo/go_api_web_domain/domain"
	"github.com/JuD4Mo/go_api_web_enrollment/internal/enrollment"
	"github.com/JuD4Mo/go_api_web_enrollment/pkg/bootstrap"
	"github.com/JuD4Mo/go_api_web_enrollment/pkg/handler"
	courseSdkMock "github.com/JuD4Mo/go_api_web_sdk/course/mock"
	userSdkMock "github.com/JuD4Mo/go_api_web_sdk/user/mock"
	"github.com/joho/godotenv"
	"github.com/ncostamagna/go_http_client/client"
)

var cli client.Transport

func TestMain(m *testing.M) {
	//Cargamos las variables de entorno que están en el archivo .env por medio del package godotenv
	_ = godotenv.Load("../.env")

	//Instanciamos un logger propio
	l := log.New(io.Discard, "", 0)

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	tx := db.Begin()

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	userSdk := &userSdkMock.UserSdkMock{
		GetMock: func(id string) (*domain.User, error) {
			return nil, nil
		},
	}

	courseSdk := &courseSdkMock.CourseSdkMock{
		GetMock: func(id string) (*domain.Course, error) {
			return nil, nil
		},
	}

	ctx := context.Background()

	enrollRepo := enrollment.NewRepo(tx, l)
	enrollService := enrollment.NewService(l, enrollRepo, userSdk, courseSdk)
	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollService, enrollment.Config{LimitPage: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)
	cli = client.New(nil, "http://"+address, 0, false)
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

	r := m.Run()

	err = srv.Shutdown(context.Background())
	if err != nil {
		l.Println(err)
	}
	tx.Rollback()
	os.Exit(r)
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
