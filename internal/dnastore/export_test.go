package dnastore

import "context"

func DropCollection(client *DnaStore, ctx context.Context) {
	client.repository.Drop(ctx)
}
