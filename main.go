package main

import (
	"math/rand"
	"time"

	_interface "github.com/GrigoryKrasnochub/go-linear-programming-task/interface"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	program := _interface.InitInterface()
	program.ShowInterface()

	//Disable logs

	//log.SetOutput(ioutil.Discard)
	//linprogtask.CalcRandom(10, 100, 100, "filepath", true)
}
