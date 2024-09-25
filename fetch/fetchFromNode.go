package fetch

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/managedblockchainquery"
)

func ListTransactions(srcAddr string) {
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
	// Call ListTransactions API. Transactions that have reached finality are always returned
	listTransactionRequest, listTransactionResponse := client.ListTransactionsRequest(&managedblockchainquery.ListTransactionsInput{
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
	})
	errors := listTransactionRequest.Send()

	if errors != nil {
		log.Fatal(errors)
	}
	f, err := os.Create(srcAddr + "_txn.chain")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(f, listTransactionResponse)
	f.Close()
}
