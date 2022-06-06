package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Post Is mutant
// swagger:operation POST /mutant ismutant ismutant
//
// validates if the entered DNA is mutant
//
// returns 200 if it is mutant and 403 if it is human
// ---
// produces:
// - application/json
// parameters:
// - name: body
//   in: body
//   required: true
//   description: information DNA
//   schema:
//    "$ref": "#/definitions/DnaRequest"
// responses:
//   '200':
//    description: the entered DNA is mutant
//   '403':
//    description: the entered DNA is humant
//   '422':
//    description: error in the DNA entered
//    examples:
//     application/json:
//      code: "Unprocessable Entity"
//      message: "mutant: invalid length dna"
//    schema:
//     "$ref": "#/definitions/Error"
//   '500':
//    description: unknown error
//    examples:
//     application/json:
//      code: "Internal Server Error"
//      message: "unknown error"
//    schema:
//     "$ref": "#/definitions/Error"
func ismutant(m Mutant) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dna DnaRequest
		if err := c.ShouldBindJSON(&dna); err != nil {
			c.Error(err)
			return
		}

		result, err := m.IsMutant(c.Request.Context(), dna.Dna)
		if err != nil {
			c.Error(err)
			return
		}

		if result {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusForbidden)
		}
	}
}

// Get Is mutant
// swagger:operation GET /stats stats stats
//
// the statistics of the processed adns
//
// returns the statistics of the processed adns
// ---
// produces:
// - application/json
// responses:
//   '200':
//    description: returns the statistics of the processed adns
//    examples:
//     application/json:
//      count_mutant_dna: 40
//      count_human_dna: 100
//      ratio: 0.4
//    schema:
//     "$ref": "#/definitions/StatsResponse"
//   '500':
//    description: unknown error
//    examples:
//     application/json:
//      code: "Internal Server Error"
//      message: "unknown error"
//    schema:
//     "$ref": "#/definitions/Error"
func stats(mutant Mutant) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := mutant.Stats(c.Request.Context())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
