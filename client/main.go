package main

import (
	"flag"
	"log"
	"time"

	"github.com/mrbardia72/sample-blockchain-grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client proto.BlockchainClient

func main() {
	
	addFlag, listFlag := MasterFlag()

	conn := GrpcDial()

	client = proto.NewBlockchainClient(conn)

	ChackFlag(addFlag, listFlag)
}

func ChackFlag(addFlag *bool, listFlag *bool) {
	if *addFlag {
		addBlock()
	}
	if *listFlag {
		getBlockchain()
	}
}

func GrpcDial() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}
	return conn
}

func MasterFlag() (*bool, *bool) {
	addFlag := flag.Bool("add", false, "Add new block")
	listFlag := flag.Bool("list", false, "List all blocks")
	flag.Parse()
	return addFlag, listFlag
}

func addBlock() {
	block, err := client.AddBlock(context.Background(), &proto.AddBlockRequest{
		Data: time.Now().String(),
	})
	if err != nil {
		log.Fatalf("unable to add block: %v", err)
	}
	log.Printf("new block hash: %s\n", block.Hash)
}

func getBlockchain() {
	blockchain, err := client.GetBlockchain(context.Background(), &proto.GetBlockchainRequest{})
	if err != nil {
		log.Fatalf("unable to get blockchain: %v", err)
	}

	log.Println("blocks:")
	for _, b := range blockchain.Blocks {
		log.Printf("hash %s, prev hash: %s, data: %s\n", b.Hash, b.PrevBlockHash, b.Data)
	}
}
