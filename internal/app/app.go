package app

import (
	"avito_tech/internal/auth"
	"avito_tech/internal/middleware"
	"avito_tech/internal/service"
	"avito_tech/internal/storage"
	"avito_tech/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func Run() {
	dbParams := postgres.NewDBParams()
	dbPool, err := postgres.NewPostgresDB(dbParams)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	usersDb := storage.NewUsersDB(dbPool)
	jwtProvider := auth.NewJwtProvider()

	authService := auth.NewAuthorizationService(usersDb, jwtProvider)
	authHandle := auth.NewHandle(authService)

	storeService := service.NewStoreService(usersDb)
	coinsService := service.NewCoinsService(usersDb)
	storeHandler := service.NewHandler(storeService, coinsService)

	postgres.RunMigrations(dbParams)

	r := chi.NewRouter()
	r.Use(middleware.Middleware(jwtProvider))

	r.Post("/api/auth", authHandle.Authorization)
	r.Post("/api/buy/{item}", storeHandler.BuyItem)
	r.Post("/api/sendCoin", storeHandler.SendCoins)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Printf("Listening on port 8080")
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
