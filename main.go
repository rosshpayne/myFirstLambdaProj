package main

import (
	"io/ioutil"
 	"context"
  	"flag"
  	"fmt"
  	"log"
  	"encoding/json"

	 "github.com/dgraph-io/dgraph/client"
  	"github.com/dgraph-io/dgraph/protos/api"
  	"google.golang.org/grpc"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxynRequest) (events.APIGatewayProxyResponse, error) {
	dgraph := string("ec2-54-206-32-30.ap-southeast-2.compute.amazonaws.com:9080")
	
	defer conn.Close()
        
	conn, err := grpc.Dial(dgraph, grpc.WithInsecure())
  	dg := client.NewDgraphClient(api.NewDgraphClient(conn))

  	resp, err := dg.NewTxn().Query(context.Background(), `{
  			bladerunner(func: eq(name@en, "Blade Runner")) {
    					uid
    					name@en
    					initial_release_date
    					netflix_id
  			}
			}`)
	
	/*
	index, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
        */
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(resp.Json),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil

}

func main() {
	lambda.Start(handler)
}
