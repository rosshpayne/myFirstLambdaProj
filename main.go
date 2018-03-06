package main

import (
 	"context"
  	"log"
        "encoding/json"

	"github.com/dgraph-io/dgraph/client"
  	"github.com/dgraph-io/dgraph/protos/api"
  	"google.golang.org/grpc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

        log.Printf("\nResource: %s", request.Resource)
        log.Printf("\nPath: %s", request.Path)
        log.Printf("\nHTTPMethod: %s",request.HTTPMethod)
        log.Printf("\nBody: %s", request.Body)
	for k,v := range request.Headers {
		log.Printf("Header:  %s  %v",k,v)
        }
	for k,v := range request.QueryStringParameters {
		log.Printf("QueryString:  %s  %v",k,v)
        }
	for k,v := range request.PathParameters {
		log.Printf("PathParameters:  %s  %v",k,v)
        }
	for k,v := range request.StageVariables {
		log.Printf("StageVariable:  %s  %v",k,v)
        }

	dgraph := string("ip-172-31-17-148.ap-southeast-2.compute.internal:9080")
	    
	conn, err := grpc.Dial(dgraph, grpc.WithInsecure())
	if err != nil {
          log.Fatal(err)
        }
	defer conn.Close()
	
	log.Printf("%s\n","About to connect to dgraph using ip-172-31-17-148.ap-southeast-2.compute.internal:9080");
	
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
	
	log.Printf("%s\n","Completed Query..now return with JSON in body");
	
        //
        // Now unmarshal 
        //
        type map_recv map[string]string

        //  type outt struct {  Name []mapt  `json:"bladerunner"` }       // go's json decode will return a pointer to array so we must define a slice of pointers.
        //  type outt struct {  Bladerunner []mapt `json:"bladerunner"` } // this works but tag is uncessary (see next example). Note field names must be upper case to make visible.

        type decode_type  struct {  Bladerunner []map_recv }              //  this works as Go will check for field name in a case insensitive manner. 
        var lines_ decode_type
	

        if  err:=json.Unmarshal([]byte(resp.Json),&lines_); err != nil {  // pass in pointer so receiver  can be populated inplace.
             panic(err)
        }

        for i,v := range lines_.Bladerunner {                             // slice of maps
           for k2,v2 := range v {
	      log.Printf("\nKey, value   %d  %s  %s",i,k2,v2)
           }
        } 

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
