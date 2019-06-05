//go:generate go run github.com/99designs/gqlgen

package storm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NicolaCignatta/storm/models"
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

// func (r *Resolver) Node() NodeResolver {
// 	return &nodeResolver{r}
// }
// func (r *Resolver) Channel() ChannelResolver {
// 	return &channelResolver{r}
// }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type nodeResolver struct{ *Resolver }
type channelResolver struct{ *Resolver }

func (r *mutationResolver) CreateNode(ctx context.Context, input NewNode) (*models.Node, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}
	txn := c.NewTxn()
	defer txn.Discard(ctx)

	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
			name: string @index(term) .
			pubKey: string .
			~node: uid @count .
		`,
	})
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	nb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	mu.SetJson = nb
	assigned, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	type Root struct {
		Node []models.Node `json:"node"`
	}

	var rt Root
	err = json.Unmarshal(resp.Json, &rt)
	if err != nil {
		return nil, err
	}

	return &rt.Node[0], nil
}

func (r *mutationResolver) CreateChannel(ctx context.Context, input NewChannel) (*models.Channel, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}
	txn := c.NewTxn()
	defer txn.Discard(ctx)

	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
			node: uid @reverse .
			capacity: int @index(int) @count .
		`,
	})
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{CommitNow: true}

	nb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	mu.SetJson = nb
	assigned, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}

	variables := map[string]string{"$id": assigned.Uids["blank-0"]}
	q := `query Channel($id: string){
		channel(func: uid($id)) {
			uid
			node
			capacity
		}
	}`

	resp, err := c.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}

	type Root struct {
		Channel []models.Channel `json:"channel"`
	}

	var rt Root
	err = json.Unmarshal(resp.Json, &rt)
	if err != nil {
		return nil, err
	}

	return &rt.Channel[0], nil
}

func (r *queryResolver) Nodes(ctx context.Context) ([]*models.Node, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}
	txn := c.NewReadOnlyTxn()

	const q = `{
		nodes (func: has(pubKey)) {
			uid
			name
			pubKey
			channel: ~node {
				uid
				capacity
			}
		}
	}`

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var jres struct {
		Nodes []*models.Node `json:"nodes"`
	}

	if err := json.Unmarshal(resp.GetJson(), &jres); err != nil {
		return nil, err
	}

	return jres.Nodes, nil
}

func (r *queryResolver) Channels(ctx context.Context) ([]*models.Channel, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}
	txn := c.NewReadOnlyTxn()

	const q = `{
		channels (func: has(capacity)) {
			uid
			node {
				uid
				name
				pubKey
			}
			capacity
		}
	}`

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var jres struct {
		Channels []*models.Channel `json:"channels"`
	}

	if err := json.Unmarshal(resp.GetJson(), &jres); err != nil {
		return nil, err
	}

	return jres.Channels, nil
}

func (r *nodeResolver) Channel(ctx context.Context, obj *models.Node) ([]*models.Channel, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}
	txn := c.NewReadOnlyTxn()

	var cUIDs []interface{}
	for _, ch := range obj.Channel {
		cUIDs = append(cUIDs, ch)
	}
	q := fmt.Sprintf(`{
		channels (func: uid(%s)) {
			uid
			capacity
		}
	}`, cUIDs...)

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var jres struct {
		Channels []*models.Channel `json:"channels"`
	}

	if err := json.Unmarshal(resp.GetJson(), &jres); err != nil {
		return nil, err
	}

	return jres.Channels, nil
}

func (r *channelResolver) Node(ctx context.Context, obj *models.Channel) ([]*models.Node, error) {
	var ns []*models.Node
	for _, n := range obj.Node {
		nn, err := getNode(ctx, n.UID)
		if err != nil {
			return nil, err
		}
		ns = append(ns, nn)
	}
	return ns, nil
}

func getNode(ctx context.Context, nUID string) (*models.Node, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}

	txn := c.NewReadOnlyTxn()

	vars := map[string]string{"$id": nUID}
	const q = `query Node($id: string){
		node (func: uid($id)) {
			uid
			name
			pubKey
		}
	}`

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var jres struct {
		Node *models.Node `json:"node"`
	}

	if err := json.Unmarshal(resp.GetJson(), &jres); err != nil {
		return nil, err
	}

	return jres.Node, nil
}

func newClient() (*dgo.Dgraph, error) {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	), nil
}
