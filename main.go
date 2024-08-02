package main

import (
	"context"
	"github.com/DavidAChris/s3lite/database"
	"github.com/DavidAChris/s3lite/services"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/glebarez/go-sqlite"
	"log"
)

var (
	Create  = "CREATE"
	Refresh = "REFRESH"
	Logout  = "LOGOUT"
)

type Event struct {
	SessionId string `json:"ses_id"`
	Email     string `json:"email"`
	Action    string `json:"action"`
}

type Response struct {
	Status    uint8  `json:"status,omitempty"`
	SessionId string `json:"ses_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Msg       string `json:"msg,omitempty"`
}

func LambdaHandler(ctx context.Context, event *Event) (*Response, error) {
	s3Svc, details := services.AwsConfig()
	log.Println("About to Init DB")
	err := database.InitDb(s3Svc, details)
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

/*
session: uuid (7e67615f-458c-40ea-bc99-5f45b077425e)
jwt: asdlkasdkasldaskd
/login cookie session_id secure
	Authorization: jwt_token 4hrs
		- jwt_token will have session_id, as long as jwt is not expired ensure session_id is the same as the claims
		- If jwt is expired, use session_id to get client info and generate a new session_id as well as a new jwt
Actions:
	Create: get email, create session_id, return session_id
	Refresh: get session_id, update session_id return email, session_id
	Invalidate: get session_id, delete session_id, return success
*/
