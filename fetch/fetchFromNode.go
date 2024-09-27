package fetch

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/managedblockchainquery"
)

func ListTransactions(srcAddr string) ([]*managedblockchainquery.TransactionOutputItem, error) {
	// Set up a session
	ambQuerySession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	}))
	client := managedblockchainquery.New(ambQuerySession)

	// Inputs for ListTransactions API
	ownerAddress := srcAddr
	network := managedblockchainquery.QueryNetworkBitcoinMainnet
	sortOrder := managedblockchainquery.SortOrderAscending
	fromTime := time.Date(1971, 1, 1, 1, 1, 1, 1, time.UTC)
	toTime := time.Now()
	nonFinal := "NONFINAL"
	txns := []*managedblockchainquery.TransactionOutputItem{}
	// Call ListTransactions API. Transactions that have reached finality are always returned

	errors := client.ListTransactionsPages(&managedblockchainquery.ListTransactionsInput{
		Address: &ownerAddress,
		Network: &network,
		Sort: &managedblockchainquery.ListTransactionsSort{
			SortOrder: &sortOrder,
		},
		FromBlockchainInstant: &managedblockchainquery.BlockchainInstant{
			Time: &fromTime,
		},
		ToBlockchainInstant: &managedblockchainquery.BlockchainInstant{
			Time: &toTime,
		},
		ConfirmationStatusFilter: &managedblockchainquery.ConfirmationStatusFilter{
			Include: []*string{&nonFinal},
		},
	}, func(lto *managedblockchainquery.ListTransactionsOutput, b bool) bool {
		txns = append(txns, lto.Transactions...)
		if lto.NextToken == nil {
			return false
		} else {
			lto.SetNextToken(*lto.NextToken)
			return true
		}
	})

	if errors != nil {
		return nil, errors
	}
	return txns, nil
}

func GeneratePassbook(txns []*managedblockchainquery.TransactionOutputItem, srcAddr string) ([][]string, error) {
	txnList, err := GetHashes(txns)
	if err != nil {
		return nil, err
	}
	var (
		Debit, Credit []string
	)
	for _, t := range txnList {
		if t["to"] == srcAddr {
			Credit = append(Credit, t["from"])
		} else {
			Debit = append(Credit, t["to"])
		}
	}
	return [][]string{Debit, Credit}, nil
}

func GetHashes(txns []*managedblockchainquery.TransactionOutputItem) (map[string]map[string]string, error) {
	ambQuerySession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
	}))
	client := managedblockchainquery.New(ambQuerySession)
	network := managedblockchainquery.QueryNetworkBitcoinMainnet
	ids := make(map[string]map[string]string)
	for _, txn := range txns {
		gto, errors := client.GetTransaction(&managedblockchainquery.GetTransactionInput{
			Network:         &network,
			TransactionHash: txn.TransactionHash,
		})
		if errors != nil {
			return ids, errors
		}
		if gto.Transaction.To == nil || gto.Transaction.From == nil {
			//these txns are representing mined blocks.
			// we can ignore for now
			continue
		}
		ids[*txn.TransactionHash] = map[string]string{
			"to":   *gto.Transaction.To,
			"from": *gto.Transaction.From,
		}
	}
	return ids, nil
}
