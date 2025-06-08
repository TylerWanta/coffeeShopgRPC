package main

import (
	pb "coffeeshop/coffeeshop_proto"
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer conn.Close()

	client := pb.NewCoffeeShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	menuStream, err := client.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("Failed to get menu: %v", err)
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}

			if err != nil {
				log.Fatalf("Could not read menu: %v", err)
			}

			items = resp.Items
			log.Printf("Menu Items Received: %v", resp.Items)
		}
	}()

	<-done

	receipt, err := client.PlaceOrder(ctx, &pb.Order{
		Items: items,
	})

	if err != nil {
		log.Fatalf("Unable to place ordre: %v", err)
	}

	log.Printf("Receipt: %v", receipt)

	status, err := client.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Fatalf("Unable to get order status: %v", err)
	}

	log.Printf("Order status: %v", status)
}
