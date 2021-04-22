package main

import (
	"fmt"
	"github.com/mrbardia72/sample-blockchain-grpc/blockchain"
	"github.com/mrbardia72/sample-blockchain-grpc/config"
	"go.mongodb.org/mongo-driver/bson/primitive"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/status"
	"log"
	"net"
	"time"

	"github.com/mrbardia72/sample-blockchain-grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var blockchainCollection = config.DbConfig().Database("blockchain").Collection("block")

type Server struct {
	Blockchain *blockchain.Blockchain
}

func main() {
	listener := NetListen()

	srv := grpc.NewServer()
	proto.RegisterBlockchainServer(srv, &Server{
		Blockchain: blockchain.NewBlockchain(),
	})
	srv.Serve(listener)
}

func NetListen() net.Listener {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("unable to listen port 8080: %v", err)
	}
	return listener
}

// AddBlock : add a new block in blockchain
func (s *Server) AddBlock(ctx context.Context, in *proto.AddBlockRequest) (*proto.AddBlockResponse, error) {

	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println(formatted,"AddBlock")
	
	block := s.Blockchain.AddBlock(in.Data)

	insertResult, err := blockchainCollection.InsertOne(ctx, block)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult)

	return &proto.AddBlockResponse{
		Hash: block.Hash,
	}, nil
}

// GetBlockchain : returns blockchain
func (s *Server) GetBlockchain(ctx context.Context, in *proto.GetBlockchainRequest) (*proto.GetBlockchainResponse, error) {

	t := time.Now()
	formatted := fmt.Sprintf("data %d-%02d-%02d Time %02d:%02d:%02d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	fmt.Println(formatted,"GetBlockchain")

	blocks :=s.Blockchain.Blocks
	cr,_ := blockchainCollection.Find(ctx, blocks)

	//fmt.Println("get all block in blockchain")
	//var results []primitive.M
	for cr.Next(ctx) {

		var elem primitive.M
		err := cr.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		//results = append(results, elem)

		resp := new(proto.GetBlockchainResponse)
		resp.Blocks = append(resp.Blocks, &proto.Block{
			PrevBlockHash:	cr.PrevBlockHash,
			Data:          cr.Data,
			Hash:          cr.Hash,
		})
	}
	//for _, b := range s.Blockchain.Blocks {
	//	resp.Blocks = append(resp.Blocks, &proto.Block{
	//		PrevBlockHash: b.PrevBlockHash,
	//		Data:          b.Data,
	//		Hash:          b.Hash,
	//	})
	//}
	return resp,nil
}
