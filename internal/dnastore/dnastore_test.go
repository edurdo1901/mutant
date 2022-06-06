package dnastore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"prueba.com/internal/dnastore"
	"prueba.com/internal/utilconection"
)

const defaultConection = "mongodb://127.0.0.1:27017"

func TestInsert(t *testing.T) {
	client := getClient(t)
	ctx := context.Background()
	err := client.Create(ctx, dnastore.DnaModel{
		Hash:     "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006ac",
		Data:     "ATDCGTT CAGTGCT TRATGTT AGAAGGT CCCCTAT TCACTGT TCACTGT",
		IsMutant: false,
	})

	assert.NoError(t, err)
}

func TestStats(t *testing.T) {
	dnaStorage := getClient(t)
	ctx := context.Background()
	dnastore.DropCollection(dnaStorage, ctx)
	populateData(t, dnaStorage)
	mutant, human, err := dnaStorage.GetStats(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, mutant)
	assert.Equal(t, 2, human)
}

func TestGetDna(t *testing.T) {
	client := getClient(t)
	ctx := context.Background()
	dnastore.DropCollection(client, ctx)
	populateData(t, client)
	model, err := client.Find(ctx, "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006ac")
	assert.NoError(t, err)
	assert.NotEmpty(t, model)

	model, err = client.Find(ctx, "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006aca")
	assert.NoError(t, err)
	assert.Empty(t, model)
}

func getClient(t *testing.T) *dnastore.DnaStore {
	dbConection, err := utilconection.OpenDB(defaultConection, "magento-DNA")
	require.NoError(t, err)
	client := dnastore.New(dbConection)
	require.NoError(t, err)
	return client
}

func populateData(t *testing.T, client *dnastore.DnaStore) {
	data := []dnastore.DnaModel{
		{
			Hash:     "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006ac",
			Data:     "ATDCGTT CAGTGCT TRATGTT AGAAGGT CCCCTAT TCACTGT TCACTDT",
			IsMutant: true,
		},
		{
			Hash:     "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006af",
			Data:     "ATDCGTT CAGTGCT TRATGTT AGAAGGT CCCCTAT TCACTGT TCACTTT",
			IsMutant: true,
		},
		{
			Hash:     "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006ag",
			Data:     "ATDCGTT CAGTGCT TRATGTT AGAAGGT CCCCTAT TCACTGT TCACTGA",
			IsMutant: false,
		},
		{
			Hash:     "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006ah",
			Data:     "ATDCGTT CAGTGCT TRATGTT AGAAGGT CCCCTAT TCACTGT TCACTGC",
			IsMutant: false,
		},
	}
	ctx := context.Background()

	for _, d := range data {
		err := client.Create(ctx, d)
		assert.NoError(t, err)
	}
}
