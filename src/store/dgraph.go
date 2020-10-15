package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

type dgraph struct {
	db *dgo.Dgraph
}

func New() *dgraph {
	host, _ := os.LookupEnv("DGRAPH_HOST")
	dialOptions := append([]grpc.DialOption{},
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	)
	connection, err := grpc.Dial(host, dialOptions...)
	if err != nil {
		log.Panic("Couldn't connect to dgraph")
	}

	db := dgo.NewDgraphClient(api.NewDgraphClient(connection))
	return &dgraph{db: db}
}

func (dg *dgraph) Setup() {
	err := dg.db.Alter(context.Background(), &api.Operation{
		Schema: `
			id: string @index(exact) .
			ip: string @index(exact) .
			from: uid @reverse .
			owner: uid @reverse .
			products_id: [uid] @reverse .
			name: string .
			age: int .
			price: string .
			device: string .
		`,
	})
	if err != nil {
		log.Panic("Couldn't setup schema to database")
	}
}

func (dg *dgraph) Save(content []byte) error {
	mutation := &api.Mutation{CommitNow: true, SetJson: content}
	_, err := dg.db.NewTxn().Mutate(context.Background(), mutation)
	return err
}

func (dg *dgraph) FindTransactions(id string) {
	variables := map[string]string{"$id": id}
	query := `
		query ($id: string) {
			var(func: eq(id, $id)) {
				ID as id
				~owner {
					from { IP as ip }
					products_id { PID as id }
				}
			}

			history(func: uid(ID)) {
				transactions: ~owner {
					id
					device
					products_id {
						id name price
					}
					from { ip }
				}
			}

			people(func: eq(ip, val(IP))) {
				sameIP: ~from {
					~owner @filter(not uid(PID)) {
						id name age
					}
				}
			}

			products(func: uid(PID)) {
				items: ~products_id @filter(not uid(PID)) {
					id name price
				}
			}
		}
	`

	res, err := dg.db.NewTxn().QueryWithVars(context.Background(), query, variables)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%s", res)
}
