package dnastore

type DnaModel struct {
	// Hash unique value of DNA.
	Hash string `bson:"hash"`

	// DNA separated by spaces.
	Data string `bson:"data"`

	// IsMutant If the DNA is mutant.
	IsMutant bool `bson:"isMutant"`
}
