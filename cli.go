package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/mikedutuandu/shippy-consignment-service/proto/consignment"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"github.com/micro/go-micro"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {

	// Create a new service
	service := micro.NewService(micro.Name("shippy.consignment.cli"))
	service.Init()


	// Create new consignment client
	consignmentClient := pb.NewShippingService("shippy.consignment.service", service.Client())
	println("Created client:",consignmentClient)

	// Contact the server and print out its response.
	file := defaultFilename
	var token string
	log.Println(os.Args)

	if len(os.Args) < 3 {
		log.Fatal(errors.New("Not enough arguments, expecing file and token."))
	}

	file = os.Args[1]
	token = os.Args[2]

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := consignmentClient.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := consignmentClient.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
