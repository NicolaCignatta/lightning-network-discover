//go:generate go run github.com/99designs/gqlgen

package storm

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateNode(ctx context.Context, input NewNode) (*Node, error) {
	c := newClient()
	txn := c.NewTxn()
	defer txn.Discard(ctx)

	mu := &api.Mutation{
		CommitNow: true,
	}

	nb, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = nb
	assigned, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

	variables := map[string]string{"$id": assigned.Uids["blank-0"]}
	q := `query Node($id: string){
		node(func: uid($id)) {
			uid
			pubKey
			name
		}
	}`

	resp, err := c.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Fatal(err)
	}

	type Root struct {
		Node []Node `json:"node"`
	}

	var rt Root
	err = json.Unmarshal(resp.Json, &rt)
	if err != nil {
		log.Fatal(err)
	}

	return &rt.Node[0], nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Nodes(ctx context.Context) ([]*Node, error) {
	c := newClient()
	txn := c.NewReadOnlyTxn()

	const q = `
		{
			nodes (func: has(pubKey)) {
				name
			}
		}
	`

	resp, err := txn.Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var jres struct {
		Nodes []*Node `json:"nodes"`
	}

	if err := json.Unmarshal(resp.GetJson(), &jres); err != nil {
		return nil, err
	}

	return jres.Nodes, nil
}

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}
