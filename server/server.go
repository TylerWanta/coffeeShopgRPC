package main

import (
	pb "coffeeshop/coffeeshop_proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(menuRequest *pb.MenuRequest, stream pb.CoffeeShop_GetMenuServer) error {
	items := []*pb.Item{
		{
			Id:   "1",
			Name: "Black Coffee",
		},
		{
			Id:   "2",
			Name: "Americano",
		},
		{
			Id:   "3",
			Name: "Vanilla Soy Chai Latte",
		},
	}

	for i := range items {
		err := stream.Send(&pb.Menu{
			Items: items[0 : i+1],
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) PlaceOrder(context.Context, *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
	}, nil
}

func (s *server) GetOrderStatus(context context.Context, reciept *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: reciept.Id,
		Status:  "IN PROGRESS",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCoffeeShopServer(grpcServer, &server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to server: %s", err)
	}
}
