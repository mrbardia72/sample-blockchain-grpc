package main

import (
	"fmt"
	"github.com/mrbardia72/sample-blockchain-grpc/blockchain"
	"github.com/mrbardia72/sample-blockchain-grpc/config"
	"log"
	"net"

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

// AddBlock : adds new block to blockchain
func (s *Server) AddBlock(ctx context.Context, in *proto.AddBlockRequest) (*proto.AddBlockResponse, error) {
	fmt.Println("AddBlock")
	//var model blockchain.Block
	block := s.Blockchain.AddBlock(in.Data)

	insertResult, err := blockchainCollection.InsertOne(ctx, block)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult)

	return &proto.AddBlockResponse{
		//Hash: block.Hash,
		Hash: block.Hash,
	}, nil


}

// GetBlockchain : returns blockchain
func (s *Server) GetBlockchain(ctx context.Context, in *proto.GetBlockchainRequest) (*proto.GetBlockchainResponse, error) {
	fmt.Println("GetBlockchain")
	resp := new(proto.GetBlockchainResponse)
	for _, b := range s.Blockchain.Blocks {
		resp.Blocks = append(resp.Blocks, &proto.Block{
			PrevBlockHash: b.PrevBlockHash,
			Data:          b.Data,
			Hash:          b.Hash,
		})
	}
	return resp, nil
}
