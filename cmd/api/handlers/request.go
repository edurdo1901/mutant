package handlers

// swagger:model
type DnaRequest struct {
	// Dna dna information.
	// example:  ["ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"]
	// required: true
	Dna []string `json:"dna" binding:"required,min=4"`
}
