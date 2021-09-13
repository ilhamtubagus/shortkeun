package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ilhamtubagus/urlShortener/api"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	//uncomment line below in production stage
	lib.LoadEnv(".env")
}
func main() {
	//Initialize database client
	client := lib.InitDatabaseClient()
	//Check if client has been found and connected to
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("MongoDB server was not found " + err.Error())
	} else {
		fmt.Println("Connected to MongoDB")
	}
	//Create new echo instance
	e := echo.New()
	api.StartApp(e, client)
	p := os.Getenv("PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

}
