# Stockfy API

A personal project with the purpose to create a REST API for a future website implementation where the customers will be able to follow their investment portfolio.
Specifically, this project is the backend foundation for such follow-up. The user will be able to register its orders, invested assets and earnings for example. 
Currently, this API will only works for assets from Brazil and United States.

This project is my first project coding in Golang and modeling a Backend application, where I expect to improve my knownledge in constructing a backend environment.
If you have some suggestion for improvement in this project feel free to contribute creating an issue. 

To construct such backend, we use this libraries and tools from Go:
- [PostgreSQL](https://www.postgresql.org/docs/): An open source object-relational database system.
- [Fiber](https://github.com/gofiber/fiber): A framework to construct our HTTP routes. It is inspired in the Express framework from Node.js.
- [Firebase](https://firebase.google.com/): Google Framework that we use for user authentication. Our work uses the [Firebase SDK](https://firebase.google.com/docs/auth)
for Go and the [REST API](https://firebase.google.com/docs/reference/rest/auth#section-api-usage) to facilitate the email verification. 
- [Pgx](https://github.com/jackc/pgx): A Golang toolkit for PostgreSQL implementations, which is the assumed database for our backend project.
- [Pgxmock](https://github.com/pashagolub/pgxmock): A mock library for the Pgx implementation. It is used in our unit tests for our database functions.
- [Net HTTP](https://pkg.go.dev/net/http): A Golang library that provied HTTP client and server implementations. In this project, it is used only to mock HTTP
requests for unit testing the routes from our REST API.

## REST API

#### Login
---
Uses the REST API from Google. Hence, we will expect that the Frontend deals with this exchange of information. Until this moment, we tested the login and 
authentication only for cases where the customer uses the sign in method with the register email and password in our Firebase project. Nowadays, we do not support
any other type of authentication such as the OAuth. Hence, the login is done using this endpoint:

HTTP Method: <code>POST</code>

HTTP Endpoint: <code>https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=[API_KEY]</code>

HTTP Body:
```json
{
	"email":"[EMAIL]",
	"password":"[PASSWORD]",
	"returnSecureToken":true
}
```

HTTP Authentication: <code>No authentication required</code>

#### Sign Up
---
Our backend works as a bridge between the user request and the Firebase, the reason for such method is because we need to save the user information in our PostgreSQL
database

HTTP Method : <code>POST</code>

HTTP Endpoint: <code>https://your.domain.path/api/signup</code>

HTTP Body:
```json
{
	"email": "[EMAIL]",
	"password": "[PASSWORD]",
	"displayName": "Pedro Test"
}
```

HTTP Authentication: <code>No authentication required</code>

#### Forgot Password
---
Our backend works as a bridge between the user request and the Firebase. As the name suggest, this request happen when the user does not remember its password and 
wants to update.

HTTP Method : <code>POST</code>

HTTP Endpoint: <code>https://your.domain.path/api/forgot-password</code>

HTTP Body:
```json
{
	"requestType": "PASSWORD_RESET",
	"email": "[EMAIL]"
}
```

HTTP Authentication: <code>No authentication required</code>

#### Delete User
---
Our backend works as a bridge between the user request and the Firebase. This endpoint has the purpose to be executed when the customer request to cancel its user
for example. Hence, this endpoint will only delete the logged user based on his valid token.

HTTP Method : <code>POST</code>

HTTP Endpoint: <code>https://your.domain.path/api/delete-user</code>

HTTP Body: <code>Empty</code>

HTTP Authentication: <code>Bearer</code>
```
Token: [FIREBASE_USER_TOKEN]
Prefix: Empty
```

#### Update User
---
Our backend works as a bridge between the user request and the Firebase. This endpoint is intended to be executed when the customer request wants to update their
information, such as the display name, password, and email. The request does not necessarily need to include those three pieces of information. Regarding the 
authentication to update the information from the user, the request needs to be authenticated with the corresponding valid token from that user, which is obtained
via the Login request presented previously.

HTTP Method : <code>POST</code>

HTTP Endpoint: <code>https://your.domain.path/api/update-user</code>

HTTP Body:
```json
{
	"displayName": "Pedro Test",
	"email": "[EMAIL]",
  "password": "[PASSWORD]"
}
```

HTTP Authentication: <code>Bearer</code>
```
Token: [FIREBASE_USER_TOKEN]
Prefix: Empty
```

## Database Organization


## File Organization

