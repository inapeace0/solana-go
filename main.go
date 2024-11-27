package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"

	"solMigrationBot/global"
)

var zero uint64 = 0

func filterLogsForBuyOrSell(log *ws.LogResult) solana.Signature {
	for _, str := range log.Value.Logs {
        if strings.Contains(str, "Transfer") {
            return log.Value.Signature;
        }
    }

	return solana.Signature{};
}

func main() {
	fmt.Println("Starting bot...")
	ctx := context.Background()
	clientRPC := rpc.New(rpc.MainNetBeta.RPC)
	clientWSS, err := ws.Connect(context.Background(), rpc.MainNetBeta.WS)
	if err != nil {
		panic(err)
	}
	fmt.Println("WS connection successful...")

	program := solana.MustPublicKeyFromBase58(global.PUMPFUN_PROGRAM_ID) // pump.fun

	{
		// Subscribe to log events that mention the provided pubkey:
		sub, err := clientWSS.LogsSubscribeMentions(
			program,
			rpc.CommitmentRecent,
		)
		if err != nil {
			panic(err)
		}
		defer sub.Unsubscribe()

		for {
			got, err := sub.Recv(ctx)
			if err != nil {
				panic(err)
			}
			// spew.Dump(got)
			signature := filterLogsForBuyOrSell(got);
			if (signature != solana.Signature{}) {
				fmt.Println(signature.String())
				start := time.Now()
				_, err := clientRPC.GetTransaction(
					context.TODO(),
					signature,
					&rpc.GetTransactionOpts{
						MaxSupportedTransactionVersion: &zero,
					},
				)

				elapsed := time.Since(start) // Calculate elapsed time
				fmt.Println("Execution time:", elapsed, "ms")

				if err != nil {
				panic(err)
				}
				//   spew.Dump(out)
			}
		}
	}
}