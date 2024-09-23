package nodes

import (
	"fmt"
	"go-chain-visualise/parser"
	"os"
	"strings"
	"sync"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Node struct {
	Wallet *parser.Wallet
	Txns   []parser.TxN
}

var (
	ErrNoSuchNode = fmt.Errorf("no such node exists")
)

func NewNode(walletid string, txns []parser.TxN) Node {
	return Node{
		Wallet: &parser.Wallet{
			ID: walletid,
		},
		Txns: txns,
	}
}

func InitNodes(src string) ([]Node, error) {
	txn, wallets, err := parser.ParseTxNs(src)
	if err != nil {
		return nil, err
	}
	var (
		loc     sync.Mutex
		network []Node
		g       sync.WaitGroup
		txnMaps = make(map[string][]parser.TxN)
	)
	g.Add(len(wallets))
	for _, w := range wallets {
		loc.Lock()
		txnMaps[w.ID] = []parser.TxN{}
		loc.Unlock()
		go func(wid string, txn []parser.TxN) {
			defer g.Done()
			for _, txnDetails := range txn {
				if txnDetails.Reciever == wid || txnDetails.Sender == wid {
					loc.Lock()
					txnMaps[wid] = append(txnMaps[wid], txnDetails)
					loc.Unlock()
				}
			}
		}(w.ID, txn)
	}
	g.Wait()

	for wid, txns := range txnMaps {
		network = append(network, NewNode(wid, txns))
	}
	return network, nil
}

func (n *Node) GenTxnVolume() map[string]int64 {
	x := make(map[string]int64)
	for _, t := range n.Txns {
		tID := n.Wallet.ID + "|" + t.Reciever
		x[tID] += 1
	}
	return x
}

type Network []Node

func (n *Network) GenVolumeGraph(outputFile string) error {
	g := graph.New(func(n *Node) string { return n.Wallet.ID }, graph.Directed())
	var out []map[string]int64
	for _, node := range *n {
		_ = g.AddVertex(&node)
		out = append(out, node.GenTxnVolume())
	}
	for _, vols := range out {
		for tid, vol := range vols {
			_ = g.AddEdge(strings.Split(tid, "|")[0], strings.Split(tid, "|")[1], graph.EdgeAttribute("txn_volume", fmt.Sprint(vol)), graph.EdgeWeight(1), graph.EdgeAttribute("Label", fmt.Sprint(vol)))
		}
	}
	file, _ := os.Create(outputFile)
	_ = draw.DOT(g, file)
	return nil
}
