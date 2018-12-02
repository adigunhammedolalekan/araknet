package main

import "github.com/adigunhammedolalekan/araknet"

func main() {

	arak := araknet.New("localhost:2299")
	arak.Connect()
}
