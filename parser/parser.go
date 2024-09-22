package parser

import (
	"encoding/csv"
	"log"
	"os"
	"slices"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type TxN struct {
	Sender    string
	Reciever  string
	TxnHash   string
	TimeStamp time.Time
	TxnAmt    float64
}

type Wallet struct {
	ID string
}

func ParseTxNs(path string) ([]TxN, []Wallet, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(file)
	csvReader.Read()
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	chunks := slices.Chunk(records, 500)
	var (
		g       errgroup.Group
		txns    []TxN
		loc     sync.Mutex
		wlock   sync.Mutex
		wallets = make(map[string]Wallet)
	)
	g.SetLimit((len(records) / 500) + (len(records) % 500))
	//attempt one using csv package.
	for chunk := range chunks {
		g.Go(func() error {
			for _, row := range chunk {
				s := row[0]
				r := row[2]
				txnHash := row[4]
				timestamp := row[5]
				ts, err := time.Parse("2006-01-02", timestamp)
				if err != nil {
					log.Println(err)
					return err
				}
				loc.Lock()
				txns = append(txns, TxN{
					Sender:    s,
					Reciever:  r,
					TxnHash:   txnHash,
					TimeStamp: ts,
				})
				loc.Unlock()
				//not checking for change in status
				wlock.Lock()
				if _, ok := wallets[s]; !ok {

					wallets[s] = Wallet{
						ID: s,
					}
				}
				if _, ok := wallets[r]; !ok {
					wallets[r] = Wallet{
						ID: r,
					}
				}
				wlock.Unlock()

			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	return txns, toArr(wallets), nil
}

func toArr[T any](in map[string]T) []T {
	var arr []T
	for _, e := range in {
		arr = append(arr, e)
	}
	return arr
}
