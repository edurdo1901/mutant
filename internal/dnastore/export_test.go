package dnastore

import "context"

func DropCollection(client *DnaStore, ctx context.Context) {
	_ = client.repository.Drop(ctx)
}
