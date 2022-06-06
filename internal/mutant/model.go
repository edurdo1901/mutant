package mutant

//swagger:model
type StatsResponse struct {
	// CountMutant number of mutants processed.
	CountMutant int `json:"count_mutant_dna"`

	// CountHuman number of humans processed.
	CountHuman int `json:"count_human_dna"`

	// Ratio relationship between (mutants/humans).
	Ratio float64 `json:"ratio"`
}
