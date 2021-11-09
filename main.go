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
	"stockfyApi/externalApi/oauth2"
	"stockfyApi/usecases"

	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func main() {

	// Database Configuration
	DB_USER := viperReadEnvVariable("DB_USER")
	DB_PASSWORD := viperReadEnvVariable("DB_PASSWORD")
	DB_NAME := viperReadEnvVariable("DB_NAME")

	// Access tokens or keys for third-party APIs
	FIREBASE_API_WEB_KEY := viperReadEnvVariable("FIREBASE_API_WEB_KEY")
	ALPHA_VANTAGE_TOKEN := viperReadEnvVariable("ALPHA_VANTAGE_TOKEN")
	FINNHUB_TOKEN := viperReadEnvVariable("FINNHUB_TOKEN")

	// Google OAuth2 Configuration
	GOOGLE_CLIENT_ID := viperReadEnvVariable("GOOGLE_CLIENT_ID")
	GOOGLE_CLIENT_SECRET := viperReadEnvVariable("GOOGLE_CLIENT_SECRET")
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
	FACEBOOK_CLIENT_ID := viperReadEnvVariable("FACEBOOK_CLIENT_ID")
	FACEBOOK_CLIENT_SECRET := viperReadEnvVariable("FACEBOOK_CLIENT_SECRET")
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

	routerConfig := router.Config{
		RouteFramework: "FIBER",
		FirebaseWebKey: FIREBASE_API_WEB_KEY,
		GoogleOAuth2:   googleOAuth2Config,
		FacebookOAuth2: facebookOAuth2Config,
	}

	router.SetupRoutes(routerConfig, applicationLogics, externalInt)

}
