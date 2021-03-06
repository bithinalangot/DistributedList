package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/bithinalangot/DistributedList/list"
)

const (
	port       = ":50053"
	reply_port = ":50052"
)

// Represent a node
type Node struct {
	data int32
	next *Node
	prev *Node
}

// Represent a linked list
type List struct {
	id       int32
	address  string
	replicas []string

	head *Node
	tail *Node

	reply chan *pb.InputResponse
}

//inserting data into linked list
func (L *List) InsertNode(ctx context.Context, in *pb.NodeRequest) (*pb.NodeResponse, error) {
	newNode := &Node{
		data: in.Data,
		next: nil,
		prev: nil,
	}
	cliConn, err := grpc.Dial("localhost"+reply_port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Didn't connect %v", err)
	}
	relayserver := pb.NewListClient(cliConn)
	resp, err := relayserver.HelloRelay(context.Background(), &pb.NodeRequest{Data: 1})

	if err != nil {
		log.Fatalf("could not relay: %v", err)
	}

	log.Print("Server1 got reponse", resp.Reply)

	if L.head == nil && L.tail == nil {
		L.head = newNode
		L.tail = newNode
	} else {
		L.tail.next = newNode
		temp := L.tail
		L.tail = newNode
		newNode.prev = temp
	}
	return &pb.NodeResponse{Success: true}, nil
}

func (L *List) HelloRelay(ctx context.Context, in *pb.NodeRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: in.String() + "Hello"}, nil
}

func (L *List) InsertInput(ctx context.Context, in *pb.InputMsg) (*pb.InputResponse, error) {
	log.Printf("Input Message %v", in.Data)
	return &pb.InputResponse{Resp: true}, nil
}

//Function to do multicast to all the replicas
func (L *List) ProcessInput(ctx context.Context, in *pb.InputMsg) (*pb.InputResponse, error) {
	for _, replica := range L.replicas {
		cliConn, err := grpc.Dial(replica, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Didn't connect %v", err)
		}
		multiserver := pb.NewListClient(cliConn)
		resp, err := multiserver.InsertInput(context.Background(), &pb.InputMsg{Data: 1})
		if err != nil {
			log.Printf("Error connecting %v", err)
		}
		log.Print(resp)
	}
	return &pb.InputResponse{Resp: true}, nil
}

//Printing the linked list
func (L *List) Printing(nodes *pb.LinkRequest, stream pb.List_PrintingServer) error {
	for temp := L.head; temp != nil; temp = temp.next {
		node := &pb.Nodes{
			Node: temp.data,
		}
		if err := stream.Send(node); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	l := List{}
	l.id = 1
	l.address = "localhost:50053"
	l.replicas = append(l.replicas, "localhost:50051", "localhost:50052", "localhost:50053", "localhost:50050")
	s := grpc.NewServer()
	pb.RegisterListServer(s, &l)
	s.Serve(lis)
}
