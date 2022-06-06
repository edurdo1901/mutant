package mutant

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"regexp"
	"strings"

	"prueba.com/internal/dnastore"
)

const _maxSecuence = 2

var (
	_dnasequence     = []string{"CCCC", "AAAA", "TTTT", "GGGG"}
	regexValidateDna = regexp.MustCompile(`[^ATGC ]`)

	ErrInvalidDna    = errors.New("mutant: invalid character dna")
	ErrInvalidLength = errors.New("mutant: invalid length dna")
)

type RepositoryDna interface {
	Create(context.Context, dnastore.DnaModel) error
	Find(context.Context, string) (dnastore.DnaModel, error)
	GetStats(context.Context) (int, int, error)
}

type Mutant struct {
	repositoryDNA RepositoryDna
}

func New(repository RepositoryDna) *Mutant {
	return &Mutant{repository}
}

// IsMutant Validate if the entered DNA is mutant or human and save the result.
func (m *Mutant) IsMutant(ctx context.Context, dna []string) (bool, error) {
	var isMutant bool
	dnaCompact, err := validateDna(dna)
	if err != nil {
		return isMutant, err
	}

	hash := hex.EncodeToString(newSHA256([]byte(dnaCompact)))
	dnaResult, err := m.repositoryDNA.Find(ctx, hash)
	if err != nil {
		return isMutant, err
	}

	if dnaResult.Hash != "" {
		return dnaResult.IsMutant, nil
	}

	isMutant = compute(dna, _maxSecuence)

	err = m.repositoryDNA.Create(ctx, dnastore.DnaModel{
		Hash:     hash,
		Data:     dnaCompact,
		IsMutant: isMutant,
	})

	if err != nil {
		return false, err
	}

	return isMutant, nil
}

// Stats consult the number of DNA of mutants and humans that have been processed.
func (m *Mutant) Stats(ctx context.Context) (StatsResponse, error) {
	var ratio float64
	countMutant, countHuman, err := m.repositoryDNA.GetStats(ctx)
	if countHuman != 0 {
		ratio = math.Floor((float64(countMutant)/float64(countHuman))*10) / 10
	}

	return StatsResponse{
		CountMutant: countMutant,
		CountHuman:  countHuman,
		Ratio:       ratio,
	}, err
}

// validateDna validates the DNA being processed is valid.
func validateDna(dna []string) (string, error) {
	length := len(dna)
	for _, d := range dna {
		if len(d) != length {
			return "", ErrInvalidLength
		}
	}

	dnaCompact := strings.Join(dna, "")
	if regexValidateDna.Match([]byte(dnaCompact)) {
		return "", ErrInvalidDna
	}

	return dnaCompact, nil
}

// compute compute if the DNA entered is from a mutant.
func compute(dna []string, maxdna int) bool {
	var sbuilderv strings.Builder
	var sbuilderh strings.Builder
	var sbuilderdp strings.Builder
	var sbuilderds strings.Builder
	var sbuilderdp1 strings.Builder
	var sbuilderds1 strings.Builder
	length := len(dna[0])
	lengthDiagonal := length - 3

	//Armando filas y columnas
	for i := 0; i < length; i++ {
		for j := 0; j < len(dna[i]); j++ {
			sbuilderh.WriteByte(dna[i][j])
			sbuilderv.WriteByte(dna[j][i])
			if i < lengthDiagonal && j < length-i {
				sbuilderdp.WriteByte(dna[i+j][j])
				sbuilderds.WriteByte(dna[i+j][length-j-1])
				if i != 0 {
					sbuilderdp1.WriteByte(dna[j][i+j])
					sbuilderds1.WriteByte(dna[j][length-j-i-1])
				}
			}
		}

		countdna := containsSequence(sbuilderh.String(),
			sbuilderv.String(),
			sbuilderdp.String(),
			sbuilderds.String(),
			sbuilderdp1.String(),
			sbuilderds1.String())

		if countdna > 0 {
			maxdna -= countdna
		}

		if maxdna <= 0 {
			return true
		}

		sbuilderv.Reset()
		sbuilderh.Reset()
		sbuilderdp.Reset()
		sbuilderds.Reset()
		sbuilderdp1.Reset()
		sbuilderds1.Reset()
	}

	return false
}

// containsSequence counts the number of entered sequences that match mutant DNA.
func containsSequence(dna ...string) int {
	count := 0
	for _, word := range dna {
		if word != "" {
			for _, secuequence := range _dnasequence {
				if strings.Contains(word, secuequence) {
					count += 1
				}
			}
		}
	}

	return count
}

// newSHA256 create a unique hash for the entered sequence.
func newSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}
