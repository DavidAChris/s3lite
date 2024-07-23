package main

import (
	"context"
	"github.com/DavidAChris/s3lite/database"
	"github.com/DavidAChris/s3lite/services"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"os"
)

type Event struct {
	Query  string `json:"query"`
	Action string `json:"action"`
}

type Response struct {
	Status string `json:"status"`
	S3Uri  string `json:"s3Uri"`
	Msg    string `json:"msg"`
}

func LambdaHandler(ctx context.Context, event *Event) (*Response, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Println("Config Failure")
		return &Response{Msg: "Server Error"}, nil
	}
	s3Svc := services.NewS3Service(config)
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("BUCKET_NAME ENV REQUIRED")
	}
	details := &services.S3Details{
		BucketName: bucketName,
		Key:        "tmp/todos.db",
		FileName:   "/tmp/todos.db",
	}
	log.Println("About to Init DB")
	err = database.InitDb(s3Svc, details)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Processing Action")
	switch event.Action {
	case "WRITE":
		log.Println("WRITE ACTION")
		err = database.Write(event.Query, s3Svc, details)
	case "ONE":
		log.Println("ROWS ACTION")
		_, err = database.Query(event.Query)
	case "MANY":
		log.Println("MANY ACTION")
		err = database.QueryRowScan(event.Query)
	default:
		log.Println("Undefined Action")
	}
	if err != nil {
		log.Fatal(err)
	}
	/*
		Ensure sqlite exists in /tmp dir. (Incase of lambda cold start)
		/tmp/{file}.db
		Execute Query based on action
		Send a response back with a success or failure and reason
	*/
	return &Response{Msg: "Success"}, err
}

func main() {
	lambda.Start(LambdaHandler)
}
