package dnastore

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const _dna = "dna"

type DnaStore struct {
	repository *mongo.Collection
}

func New(db *mongo.Database) *DnaStore {
	return &DnaStore{
		repository: db.Collection(_dna),
	}
}

// Create save dna model.
func (c *DnaStore) Create(ctx context.Context, dna DnaModel) error {
	_, err := c.repository.InsertOne(ctx, dna)
	return err
}

// Find search with the identifier if the DNA already exists.
func (c *DnaStore) Find(ctx context.Context, hash string) (DnaModel, error) {
	var result DnaModel
	filter := bson.M{"hash": hash}
	err := c.repository.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, nil
		}
	}

	return result, err
}

// GetStats check the number of human mutants entered.
func (c *DnaStore) GetStats(ctx context.Context) (int, int, error) {
	var countMutant int
	var countHuman int

	group := bson.D{
		{Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$isMutant"},
				{Key: "count",
					Value: bson.D{{Key: "$sum", Value: 1}},
				},
			},
		},
	}

	var results []bson.M
	cursor, err := c.repository.Aggregate(ctx, mongo.Pipeline{group})
	if err != nil {
		return 0, 0, err
	}

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0, err
	}

	for _, data := range results {
		isMutant := data["_id"].(bool)
		if isMutant {
			countMutant = int(data["count"].(int32))
		} else {
			countHuman = int(data["count"].(int32))
		}
	}

	return countMutant, countHuman, nil
}
