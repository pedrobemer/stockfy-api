package main

import (
	"context"
	"fmt"
	"os"
	"stockfyApi/api/router"
	"stockfyApi/database/postgresql"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/externalApi/finnhub"
	"stockfyApi/externalApi/firebaseApi"
	"stockfyApi/usecases"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func main() {
	DB_USER := viperReadEnvVariable("DB_USER")
	DB_PASSWORD := viperReadEnvVariable("DB_PASSWORD")
	DB_NAME := viperReadEnvVariable("DB_NAME")
	FIREBASE_API_WEB_KEY := viperReadEnvVariable("FIREBASE_API_WEB_KEY")
	ALPHA_VANTAGE_TOKEN := viperReadEnvVariable("ALPHA_VANTAGE_TOKEN")
	FINNHUB_TOKEN := viperReadEnvVariable("FINNHUB_TOKEN")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	DBpool, err := pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DBpool.Close(context.Background())

	auth := firebaseApi.SetupFirebase("stockfy-api-firebase-adminsdk-cwuka-f2c828fb90.json")
	firebaseInterface := firebaseApi.NewFirebase(auth)

	dbInterfaces := postgresql.NewPostgresInstance(DBpool)

	applicationLogics := usecases.NewApplications(dbInterfaces, firebaseInterface)

	finnhubInterface := finnhub.NewFinnhubApi(FINNHUB_TOKEN)
	alphaInterface := alphaVantage.NewAlphaVantageApi(ALPHA_VANTAGE_TOKEN)

	externalInt := externalapi.ThirdPartyInterfaces{
		FinnhubApi:      *finnhubInterface,
		AlphaVantageApi: *alphaInterface,
	}

	router.SetupRoutes("FIBER", FIREBASE_API_WEB_KEY, applicationLogics,
		externalInt)

}
