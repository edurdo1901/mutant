# Mutant API

[![Build Status](https://github.com/edurdo1901/mutant/workflows/Build/badge.svg?branch=main)](https://github.com/edurdo1901/mutant/actions?query=branch%3Amain) [![Coverage Status](https://coveralls.io/repos/github/edurdo1901/mutant/badge.svg?branch=main)](https://coveralls.io/github/edurdo1901/mutant?branch=main)

## Sitio de la aplicación

[Prueba mercado libre](https://challenge-golang.gentlesea-9f37728d.westus.azurecontainerapps.io) `https://challenge-golang.gentlesea-9f37728d.westus.azurecontainerapps.io`

## Documentación API

Los dos `endpoints` se encuentran documentados en el siguiente archivo [Mutant API swagger](docs/swagger.yaml)
## Requerimientos

- docker
- make
- go

## Ejecutar aplicación

Para ejecutar la aplicación, se debe crear una instancia local de mongo, la cual se crea el siguiente comando `docker run --name some-mongo -p 27017:27017 -d mongo`. Seguido de esto, ejecutamos `go run cmd/main.go` para correr el proyecto.

Se puede modificar el .env para ajustar el puerto por el que se está ejecutando la aplicación y la cadena de conexión.

## Pruebas

Para ejecutar las pruebas unitarias, el comando `make test-docker` crea una instancia de mongo para posteriormente ejecutar las pruebas en esa base de datos.

## Cobertura de pruebas

Para verificar la cobertura del código, se debe tener una instancia local de mongo. Para ello, se corre comando `docker run --name some-mongo -p 27017:27017 -d mongo` y posteriormente se ejecuta `make test-cover`.
___


> **⚠ WARNING: Puede ser que el puerto 27017 este ocupado y por eso no logre ejecutar las pruebas unitarias. ** 

