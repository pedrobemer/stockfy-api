# Stockfy API

A personal project with the purpose to create a REST API for a future website implementation where the customers will be able to follow their investment portfolio.
Specifically, this project is the backend foundation for such follow-up. The user will be able to register its orders, invested assets and earnings for example. 
Currently, this API will only works for assets from Brazil and United States.

This project is my first project coding in Golang and modeling a Backend application, where I expect to improve my knowledge in constructing a backend environment. In the project's current stage, we tried to follow the clean architecture principles for software development. If you do not know this type of architecture, please click on this link: [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

If you have some suggestion for improvement in this project feel free to contribute creating an issue. 

To construct such backend, we use these libraries and tools from Go:
- [PostgreSQL](https://www.postgresql.org/docs/): An open source object-relational database system.
- [Fiber](https://github.com/gofiber/fiber): A framework to construct our HTTP routes. It is inspired in the Express framework from Node.js.
- [Firebase](https://firebase.google.com/): Google Framework that we use for user authentication. Our work uses the [Firebase SDK](https://firebase.google.com/docs/auth)
for Go and the [REST API](https://firebase.google.com/docs/reference/rest/auth#section-api-usage) to facilitate the email verification. 
- [Pgx](https://github.com/jackc/pgx): A Golang toolkit for PostgreSQL implementations, which is the assumed database for our backend project.
- [Pgxmock](https://github.com/pashagolub/pgxmock): A mock library for the Pgx implementation. It is used in our unit tests for our database functions.
- [Net HTTP](https://pkg.go.dev/net/http): A Golang library that provides a HTTP client and server. In this project, it is used only to mock HTTP
requests for unit testing the routes from our REST API.
- [Finnhub](https://finnhub.io/docs/api): A RESTful API for real-time information regarding investiments around the world. Nevertheless, as a brazilian investor, this database does not have all the possible assets such as stocks without ownership in the company (BBDC4, ITUB4) and real estate funds.
- [Alpha Vantage](https://finnhub.io/docs/api): A RESTful API for real-time information regarding investments around the world. Unlike the Finnhub, this application has all the possible assets from Brazil. Nevertheless, the free version enables few request per hour in comparison with the Finnhub.

## Project Organization

In this project, we follow the clean architecture principles for software development, where our applications logic do not depend on the existence of a specific framework or library to work correctly. The idea is that we do not want our code to be attached, for example, with PostgreSQL, Finnhub, or any other third-party library or API. Hence, in this section, we will explain how this project is organized.

    .
    ├── api                   # API folder (It is in the layer of Framework & Drivers)
    ├── client                # HTTP client to send request from our API (It is in the layer of Framework & Drivers)
    ├── database              # Database Source Files (It is in the layer of Framework & Drivers)
    ├── entity                # Encapsulated wide method rules (It is in the Entities layer)
    ├── externalApi           # External API that we use in our backend (It is in the layer of Framework & Drivers)
    ├── usecases	      # Application logic folder (It is in the Use Cases layer)
    ├── main.go
    ├── envVariable.go
    ├── go.mod
    ├── go.sum
    └── README.md
    
   
   Here, we assume that lower layers cannot access types and methods from the upper layers. For example, our entity folder (the lowest layer of all) won't be able to execute or receive any method or type from the upper layers ("usecases" and "api" folder, for example). The diagram below demonstrates the organization's hierarchy from bottom to top with the layer's name and which folders are part of these layers. So, all the layers below "Entity" can execute and use this folder's type and methods. With the same logic applying to the other layers.
    
    .
    ├── 1°: Entity Layer - entity                  
    ├── 2°: Use Cases - usecases               
    ├── 3°: Presenters/Controllers - interfaces in the usescases folder           
    ├── 4°: Framework & Drivers - externalApi, api, client, database          


Our application logic developed in the Use Cases does not know what type of database or API is running in the Framework & Drivers. This is possible by implementing an intermediate layer that executes the controllers and presenter's methods. Basically, we do not have a specific folder for such layer, but in Golang, they can be deployed using interfaces. In our code, each folder from "usecases" has an interface file that specifies the methods necessary for that specific package. This interface connects the methods used in the database. Still, as they do not know which database is, any database can be used as long as the implementation follows that method declaration. So our API can send and retrieve some data from our database without knowing which kind of database is running. The same applies to the database; they do not know which sort of API is running.

## Database Organization (PostgreSQL)

![alt text](https://github.com/itelonog/stockfy-api/blob/clean_arch/database.png)

## REST API 

Currently, we are developing the REST API. So, this is not a final version. 

To see the current stage of our API, we are using Swagger to document the API following the Open API 3.0 convention:
[Documentation](https://app.swaggerhub.com/apis-docs/pedrobemer/Stockfy/1.0.0#/)

