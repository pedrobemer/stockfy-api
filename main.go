package main

import (
	"context"
	"fmt"
	"os"
	"stockfyApi/api/router"
	"stockfyApi/client"
	"stockfyApi/database/postgresql"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/externalApi/finnhub"
	"stockfyApi/externalApi/firebaseApi"
	"stockfyApi/externalApi/oauth2"
	"stockfyApi/usecases"
	"stockfyApi/usecases/utils"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func main() {

	// Database Configuration
	filenamePath := "./"
	filename := "database"
	DB_USER := utils.ViperReadEnvVariable(filenamePath, filename, "DB_USER")
	DB_PASSWORD := utils.ViperReadEnvVariable(filenamePath, filename, "DB_PASSWORD")
	DB_NAME := utils.ViperReadEnvVariable(filenamePath, filename, "DB_NAME")
	DB_PORT := utils.ViperReadEnvVariable(filenamePath, filename, "DB_PORT")
	DB_HOST := utils.ViperReadEnvVariable(filenamePath, filename, "DB_HOST")

	// Access tokens or keys for third-party APIs
	FIREBASE_API_WEB_KEY := utils.ViperReadEnvVariable(filenamePath, filename,
		"FIREBASE_API_WEB_KEY")
	ALPHA_VANTAGE_TOKEN := utils.ViperReadEnvVariable(filenamePath, filename,
		"ALPHA_VANTAGE_TOKEN")
	FINNHUB_TOKEN := utils.ViperReadEnvVariable(filenamePath, filename,
		"FINNHUB_TOKEN")

	// Google OAuth2 Configuration
	GOOGLE_CLIENT_ID := utils.ViperReadEnvVariable(filenamePath, filename,
		"GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET := utils.ViperReadEnvVariable(filenamePath, filename,
		"GOOGLE_CLIENT_SECRET")
	GOOGLE_REDIRECT_URI := "http://localhost:3000/api/signin/oauth2/google"
	GOOGLE_SCOPE := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	GOOGLE_AUTHORIZATION_ENDPOINT := "https://accounts.google.com/o/oauth2/auth"
	GOOGLE_ACCESS_TOKEN_ENDPOINT := "https://oauth2.googleapis.com/token"

	googleOAuth2Config := oauth2.GoogleOAuthConfig(GOOGLE_CLIENT_ID,
		GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URI, GOOGLE_SCOPE,
		GOOGLE_AUTHORIZATION_ENDPOINT, GOOGLE_ACCESS_TOKEN_ENDPOINT)

	// Facebook OAuth2 Configuration
	FACEBOOK_CLIENT_ID := utils.ViperReadEnvVariable("./", filename,
		"FACEBOOK_CLIENT_ID")
	FACEBOOK_CLIENT_SECRET := utils.ViperReadEnvVariable("./", filename,
		"FACEBOOK_CLIENT_SECRET")
	FACEBOOK_REDIRECT_URI := "http://localhost:3000/api/signin/oauth2/facebook"
	FACEBOOK_SCOPE := []string{
		"email",
		"public_profile",
	}
	FACEBOOK_AUTHORIZATION_ENDPOINT := "https://www.facebook.com/v12.0/dialog/oauth"
	FACEBOOK_ACCESS_TOKEN_ENDPOINT := "https://graph.facebook.com/v12.0/oauth/access_token"

	facebookOAuth2Config := oauth2.FacebookOAuthConfig(FACEBOOK_CLIENT_ID,
		FACEBOOK_CLIENT_SECRET, FACEBOOK_REDIRECT_URI, FACEBOOK_SCOPE,
		FACEBOOK_AUTHORIZATION_ENDPOINT, FACEBOOK_ACCESS_TOKEN_ENDPOINT)

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

	DBpool, err := pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DBpool.Close(context.Background())

	auth := firebaseApi.SetupFirebase("stockfy-firebase-admin.json")
	firebaseInterface := firebaseApi.NewFirebase(auth)

	dbInterfaces := postgresql.NewPostgresInstance(DBpool)

	applicationLogics := usecases.NewApplications(dbInterfaces, firebaseInterface)

	finnhubInterface := finnhub.NewFinnhubApi(FINNHUB_TOKEN,
		client.RequestAndAssignToBody)
	alphaInterface := alphaVantage.NewAlphaVantageApi(ALPHA_VANTAGE_TOKEN,
		client.RequestAndAssignToBody)

	externalInt := externalapi.ThirdPartyInterfaces{
		FinnhubApi:      finnhubInterface,
		AlphaVantageApi: alphaInterface,
	}

	routerConfig := router.Config{
		RouteFramework: "FIBER",
		FirebaseWebKey: FIREBASE_API_WEB_KEY,
		GoogleOAuth2:   googleOAuth2Config,
		FacebookOAuth2: facebookOAuth2Config,
	}

	router.SetupRoutes(routerConfig, applicationLogics, externalInt)

}
