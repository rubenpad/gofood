package store

import (
	"context"
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
			name: string .
			age: int .
			price: string .
			device: string .
			products_id: [uid] .	
		`,
	})
	if err != nil {
		log.Panic("Couldn't setup schema to database")
	}
}

func (dg *dgraph) MakeMutation(content []byte) error {
	mutation := &api.Mutation{CommitNow: true, SetJson: content}
	_, err := dg.db.NewTxn().Mutate(context.Background(), mutation)
	return err
}
