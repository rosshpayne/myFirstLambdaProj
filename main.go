package main

import (
 	"context"
  	"log"
  	

	 "github.com/dgraph-io/dgraph/client"
  	"github.com/dgraph-io/dgraph/protos/api"
  	"google.golang.org/grpc"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	dgraph := string("ip-172-31-17-148.ap-southeast-2.compute.internal:9080")
	    
	conn, err := grpc.Dial(dgraph, grpc.WithInsecure())
	if err != nil {
          log.Fatal(err)
        }
	defer conn.Close()
	
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
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(handler)
}
