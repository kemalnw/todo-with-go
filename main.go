package main

import (
	"context"
	"log"
	"time"

	_todoHttp "github.com/kemalnw/todo-with-go/todo/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	viper.SetConfigFile(`config.json`)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("service RUN on DEBUG mode")
	}
}

func main() {
	dsn := viper.GetString(`database.dsn`)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatalln("dsn is incorrect,", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalln("unable to ping the database, make sure its running and authentication is correct, err:", err)
	}

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	database := client.Database(viper.GetString(`database.name`))

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	_todoHttp.NewTodoHandler(e, database)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
