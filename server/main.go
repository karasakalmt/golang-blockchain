package main

import (
	"blockchain/blockchain"
	proto "blockchain/proto"
	"context"
	"fmt"
	"net"
	"runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Block struct {
	PrevHash []byte
	Data     string
	Hash     []byte
	Pow      bool
}

type blockChainServer struct {
	proto.UnimplementedBlockChainServiceServer
}

func main() {

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterBlockChainServiceServer(srv, &blockChainServer{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *blockChainServer) PrintBlocks(ctx context.Context, req *proto.PrintRequest) (*proto.PrintResponse, error) {
	iter := blockchain.InitBlockChain().Iterator()
	printRes := &proto.PrintResponse{}
	fmt.Println("block")

	for {
		currentBlock := iter.Next()
		pow := blockchain.NewProof(currentBlock)

		block := &proto.PrintResponse_Result{
			PrevHash: currentBlock.PrevHash,
			Data:     string(currentBlock.Data),
			Hash:     currentBlock.Hash,
			Pow:      pow.Validate(),
		}
		if len(currentBlock.PrevHash) == 0 {
			break
		}
		fmt.Println(block)

		printRes.Results = append(printRes.Results, block)
	}
	return printRes, nil
}

func (s *blockChainServer) AddBlock(ctx context.Context, req *proto.AddBlockRequest) (*proto.AddBlockResponse, error) {
	addBlockData := req.GetData()
	fmt.Printf("%s", addBlockData)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	if addBlockData == "" {
		runtime.Goexit()
	}
	chain.AddBlock(addBlockData)

	return &proto.AddBlockResponse{Data: "Added Block!"}, nil
}
