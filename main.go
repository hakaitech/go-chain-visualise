package main

import (
	"go-chain-visualise/fetch"
	"go-chain-visualise/nodes"
	"go-chain-visualise/utils"
	"log"
	"time"
)

func main() {
	go func() {
		for {
			utils.PrintMemUsage()
			time.Sleep(5 * time.Second)
		}
	}()
	network, err := nodes.InitNodes("bitcoin_transactions_8-9-24_1.csv")
	if err != nil {
		log.Fatal(err)
	}
	// g := graph.New(func(n nodes.Node) string { return n.Wallet.ID }, graph.Directed())
	// var w sync.WaitGroup
	// w.Add(len(network))
	// for _, node := range network {
	// 	_ = g.AddVertex(node)
	// }
	// for _, node := range network {
	// 	go func(node nodes.Node) {
	// 		defer w.Done()
	// 		for _, txn := range node.Txns {
	// 			if txn.Sender == node.Wallet.ID {
	// 				_ = g.AddEdge(node.Wallet.ID, txn.Reciever, graph.EdgeWeight(1))
	// 			}
	// 		}
	// 	}(node)
	// }
	// w.Wait()
	// file, _ := os.Create("map.gv")
	// _ = draw.DOT(g, file)
	net := nodes.Network(network[:20])
	net.GenVolumeGraph("vol.gv")
	fetch.ListTransactions("bc1qz9jpuvmex3hwpleczwqzw6vs3v7mmxxr02048a")
}
