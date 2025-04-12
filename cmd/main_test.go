package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

var originalConnectToMongoDBFunc = connectToMongoDBFunc
var originalRunServerFunc = runServerFunc

func dummyConnectToMongoDB(uri, dbName string) (*mongo.Client, *mongo.Database, error) {
	return new(mongo.Client), new(mongo.Database), nil
}

func TestRun_MissingMongoURI(t *testing.T) {
	os.Unsetenv("MONGO_URI")
	os.Setenv("MONGO_DB_NAME", "testdb")
	defer os.Unsetenv("MONGO_DB_NAME")

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MONGO_URI")
}

func TestRun_MissingMongoDBName(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Unsetenv("MONGO_DB_NAME")
	defer func() {
		os.Unsetenv("MONGO_URI")
		os.Unsetenv("MONGO_DB_NAME")
	}()

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MONGO_DB_NAME")
}

func TestRun_ConnectError(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "testdb")
	defer func() {
		os.Unsetenv("MONGO_URI")
		os.Unsetenv("MONGO_DB_NAME")
	}()

	connectToMongoDBFunc = func(uri, dbName string) (*mongo.Client, *mongo.Database, error) {
		return nil, nil, errors.New("simulated connect error")
	}
	defer func() {
		connectToMongoDBFunc = originalConnectToMongoDBFunc
	}()

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated connect error")
}

func TestRun_ServerError(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "testdb")
	os.Setenv("PORT", "1234")
	defer func() {
		os.Unsetenv("MONGO_URI")
		os.Unsetenv("MONGO_DB_NAME")
		os.Unsetenv("PORT")
	}()

	originalDisconnectFunc := disconnectFunc
	disconnectFunc = func(client *mongo.Client) error { return nil }
	defer func() { disconnectFunc = originalDisconnectFunc }()

	runServerFunc = func(addr string, router *gin.Engine) error {
		return errors.New("simulated server error")
	}
	defer func() {
		connectToMongoDBFunc = originalConnectToMongoDBFunc
		runServerFunc = originalRunServerFunc
	}()

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated server error")
}

func TestRun_Success(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "testdb")
	os.Setenv("PORT", "5678")
	defer func() {
		os.Unsetenv("MONGO_URI")
		os.Unsetenv("MONGO_DB_NAME")
		os.Unsetenv("PORT")
	}()

	originalDisconnectFunc := disconnectFunc
	disconnectFunc = func(client *mongo.Client) error { return nil }
	defer func() { disconnectFunc = originalDisconnectFunc }()

	serverCalled := false
	runServerFunc = func(addr string, router *gin.Engine) error {
		serverCalled = true
		if addr != ":5678" {
			return errors.New("incorrect address")
		}
		return nil
	}
	defer func() {
		connectToMongoDBFunc = originalConnectToMongoDBFunc
		runServerFunc = originalRunServerFunc
	}()

	err := run()
	assert.NoError(t, err)
	assert.True(t, serverCalled, "runServerFunc debió haberse llamado")
}

func TestConnectToMongoDB_InvalidURI(t *testing.T) {
	_, _, err := connectToMongoDB("invalid-uri", "testdb")
	assert.Error(t, err)
}
