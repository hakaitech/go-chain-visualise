package main

import (
	"fmt"
	"go-chain-visualise/fetch"
	"go-chain-visualise/utils"
	"log"
	"time"
)

func main() {
	go func() {
		for {
			utils.PrintMemUsage()
			time.Sleep(time.Second)
		}
	}()
	t, err := fetch.ListTransactions("bc1qz9jpuvmex3hwpleczwqzw6vs3v7mmxxr02048a")
	if err != nil {
		log.Fatal(err)
	}
	pb, err := fetch.GeneratePassbook(t, "bc1qz9jpuvmex3hwpleczwqzw6vs3v7mmxxr02048a")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pb)
}
