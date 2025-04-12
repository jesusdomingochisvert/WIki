package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/application/usecases"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/adapters"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/interfaces"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/infrastructure/repositories"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/presentations/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var connectToMongoDBFunc = connectToMongoDB

var disconnectFunc = func(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}

var runServerFunc = func(addr string, router *gin.Engine) error {
	return router.Run(addr)
}

func connectToMongoDB(uri, dbName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	database := client.Database(dbName)
	return client, database, nil
}

func setupRouter(db interfaces.Database) *gin.Engine {
	router := gin.Default()

	repository := repositories.NewUserRepository(db)
	useCase := usecases.NewGetAllUsersUseCase(repository)
	handler := handlers.NewUserHandler(useCase)

	router.GET("/users", func(c *gin.Context) {
		users, err := handler.GetAllUsers(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"users": users})
	})

	return router
}

func run() error {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return errors.New("the Env Var MONGO_URI is mandatory")
	}

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		return errors.New("the Env Var MONGO_DB_NAME is mandatory")
	}

	client, database, err := connectToMongoDBFunc(mongoURI, mongoDBName)
	if err != nil {
		return fmt.Errorf("error conectándose a MongoDB: %v", err)
	}

	defer func() {
		if err := disconnectFunc(client); err != nil {
			log.Printf("error to disconnect of MongoDB: %v", err)
		}
	}()

	dbAdapter := adapters.NewMongoDatabaseAdapter(database)
	router := setupRouter(dbAdapter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("Iniciando servidor en %s", addr)
	if err := runServerFunc(addr, router); err != nil {
		return fmt.Errorf("error al iniciar el servidor: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("Hello World")
}
