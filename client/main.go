package main

import (
	proto "blockchain/proto"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type AddBlockBody struct {
	Data string `json:"data"`
}

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewBlockChainServiceClient(conn)

	g := gin.Default()

	g.GET("/print", func(ctx *gin.Context) {
		req := &proto.PrintRequest{}
		if response, err := client.PrintBlocks(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"Data": response.Results})
		} else {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		}
	})

	g.POST("/add-block", func(ctx *gin.Context) {
		var addBlockBody AddBlockBody

		if err := ctx.BindJSON(&addBlockBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "There is an error with the body you sent!"})
		}

		req := &proto.AddBlockRequest{Data: addBlockBody.Data}
		if response, err := client.AddBlock(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(response.Data),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error!"})
		}
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
