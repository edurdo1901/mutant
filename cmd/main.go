package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"prueba.com/cmd/api/handlers"
	"prueba.com/internal/dnastore"
	"prueba.com/internal/mutant"
	"prueba.com/internal/utilconection"
)

const _db = "magento-DNA"

func main() {
	loadEnv()
	port := os.Getenv("port")
	conectionString := os.Getenv("conectionString")
	db, err := utilconection.OpenDB(conectionString, _db)
	if err != nil {
		panic(err)
	}

	dnaStorage := dnastore.New(db)
	mutant := mutant.New(dnaStorage)
	r := gin.Default()
	initPrometheus(r)
	handler := handlers.New(mutant)
	handler.API(r)
	r.Run(":" + port)
}

// initPrometheus start Prometheus metrics.
func initPrometheus(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// loadEnv load the environment variables
func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}
