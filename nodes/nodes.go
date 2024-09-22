package nodes

import (
	"fmt"
	"go-chain-visualise/parser"
	"sync"
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
