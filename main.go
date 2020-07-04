package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/GrigoryKrasnochub/go-linear-programming-task/linprogtask"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	//program := _interface.InitInterface()
	//program.ShowInterface()

	//Disable logs

	log.SetOutput(ioutil.Discard)
	linprogtask.CalcRandom(10, 100, 100, "C:\\go\\src\\github.com\\GrigoryKrasnochub\\go-linear-programming-task\\builds\\test.txt", true)
}
