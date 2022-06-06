[![Build Status](https://github.com/edurdo1901/mutant/workflows/Build/badge.svg?branch=main)](https://github.com/edurdo1901/mutant/actions?query=branch%3Amain) [![Coverage Status](https://coveralls.io/repos/github/edurdo1901/mutant/badge.svg?branch=main)](https://coveralls.io/github/edurdo1901/mutant?branch=main)

## Requerimientos
- docker
- make
- go

## Test

Para ejecutar la pruebas unitarias solo ejecutar el comando `make test-docker` este comando crear un mongo para posteriormente ejecutar las pruebas en esa base de datos.

Para correr la covertura del codigo se tiene que tener una instancia local de mongo que se crea con el siguiente comando `docker run --name some-mongo -p 27017:27017 -d mongo` y posteriormente ejecutar `make test-cover`.

## Correr la aplicación

Para correr la aplicación se tiene que tener una instancia local de mongo que se crea con el siguiente comando `docker run --name some-mongo -p 27017:27017 -d mongo` despues de ejecutar esta comando ya podemos correr el proyecto `go run cmd/main.go`

Se puede modificar el .env para ajustar el puerto por el que se esta ejecutando la aplicación y la cadeba de conección.

## Documentación API

Los dos end point se encuentran documentados en el siguiente archivo [Swagger](docs/swagger.yaml)

## Sitio de la aplicación

[Prueba mercado libre](https://challenge-golang.gentlesea-9f37728d.westus.azurecontainerapps.io)



