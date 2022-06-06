package mutant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"prueba.com/internal/dnastore"
	"prueba.com/internal/mutant"
)

type storageMock struct {
	mock.Mock
}

func (m *storageMock) Create(ctx context.Context, model dnastore.DnaModel) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *storageMock) Find(ctx context.Context, hash string) (dnastore.DnaModel, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).(dnastore.DnaModel), args.Error(1)
}

func (m *storageMock) GetStats(ctx context.Context) (int, int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Get(1).(int), args.Error(2)
}

type mockArgs struct {
	methodName string
	inputArgs  []interface{}
	returnArgs []interface{}
	times      int
}

func TestIsMutant(t *testing.T) {
	ctx := context.Background()
	sMockDefault := []mockArgs{
		{
			methodName: "Find",
			inputArgs:  []interface{}{ctx, mock.Anything},
			returnArgs: []interface{}{dnastore.DnaModel{}, nil},
			times:      1,
		},
		{
			methodName: "Create",
			inputArgs:  []interface{}{ctx, mock.Anything},
			returnArgs: []interface{}{nil},
			times:      1,
		},
	}

	tt := map[string]struct {
		dna            []string
		expectedMutant bool
		sMock          []mockArgs
		expecterErr    error
	}{
		"MutantIsTrue": {
			dna:            []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"},
			expectedMutant: true,
			sMock:          sMockDefault,
		},
		"MutantFalse": {
			dna:            []string{"TTGGTA", "CAGTGC", "TTATGT", "TGATGT", "CCCCTA", "TCACTG"},
			expectedMutant: false,
			sMock:          sMockDefault,
		},
		"MutantGetResult": {
			dna:            []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAGGG", "CCCCTA", "TCACTG"},
			expectedMutant: false,
			sMock: []mockArgs{
				{
					methodName: "Find",
					inputArgs:  []interface{}{ctx, mock.Anything},
					returnArgs: []interface{}{dnastore.DnaModel{
						Hash: "08129510eee94be205644d23eae66988eb34c23bde9ea9daf4b9b6e72a0006ac",
					}, nil},
					times: 1,
				},
			},
		},
		"MutantGetResultErrorPayloadLength": {
			dna:            []string{"ATGCGAA", "CAGTGC", "TTATGT", "AGAGGG", "CCCCTA", "TCACTG"},
			expectedMutant: false,
			expecterErr:    mutant.ErrInvalidLength,
		},
		"MutantGetResultErrorPayloadDNA": {
			dna:            []string{"ATGCGY", "CAGTGC", "TTATGT", "AGAGGG", "CCCCTA", "TCACTG"},
			expectedMutant: false,
			expecterErr:    mutant.ErrInvalidDna,
		},
		"MutantErrorGet": {
			dna:            []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAGGG", "CCCCTA", "TCACTG"},
			expectedMutant: false,
			sMock: []mockArgs{
				{
					methodName: "Find",
					inputArgs:  []interface{}{ctx, mock.Anything},
					returnArgs: []interface{}{dnastore.DnaModel{}, errors.New("Error")},
					times:      1,
				},
			},
			expecterErr: errors.New("Error"),
		},
		"MutantErrorSave": {
			dna:            []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAGGG", "CCCCTA", "TCACTG"},
			expectedMutant: false,
			sMock: []mockArgs{
				{
					methodName: "Find",
					inputArgs:  []interface{}{ctx, mock.Anything},
					returnArgs: []interface{}{dnastore.DnaModel{}, nil},
					times:      1,
				},
				{
					methodName: "Create",
					inputArgs:  []interface{}{ctx, mock.Anything},
					returnArgs: []interface{}{errors.New("Error")},
					times:      1,
				},
			},
			expecterErr: errors.New("Error"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			var sMock storageMock
			for _, configMock := range tc.sMock {
				setupMock(configMock, &sMock.Mock)
			}

			testMagento := mutant.New(&sMock)
			result, err := testMagento.IsMutant(ctx, tc.dna)
			assert.EqualValues(t, result, tc.expectedMutant)
			if tc.expecterErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expecterErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mock.AssertExpectationsForObjects(t, &sMock)
		})
	}
}

func TestStats(t *testing.T) {
	ctx := context.Background()

	tt := map[string]struct {
		expectedRatio float64
		sMock         mockArgs
		expecterErr   error
	}{
		"GetStats": {
			expectedRatio: 0.4,
			sMock: mockArgs{
				methodName: "GetStats",
				inputArgs:  []interface{}{ctx},
				returnArgs: []interface{}{40, 100, nil},
			},
		},
		"GetStatsError": {
			expectedRatio: 0,
			sMock: mockArgs{
				methodName: "GetStats",
				inputArgs:  []interface{}{ctx},
				returnArgs: []interface{}{0, 0, errors.New("Error")},
			},
			expecterErr: errors.New("Error"),
		},
	}

	for _, tc := range tt {
		var sMock storageMock
		setupMock(tc.sMock, &sMock.Mock)
		testMagento := mutant.New(&sMock)

		stats, err := testMagento.Stats(ctx)

		if tc.expecterErr != nil {
			assert.Error(t, err)
			assert.Empty(t, stats)
		} else {
			assert.NoError(t, err)
			assert.EqualValues(t, tc.expectedRatio, stats.Ratio)
		}
	}
}

func setupMock(params mockArgs, mock *mock.Mock) {
	if params.methodName != "" {
		mock.On(params.methodName, params.inputArgs...).
			Return(params.returnArgs...).
			Times(1)
	}
}
