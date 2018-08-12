package main

import (
	"github.com/kgigitdev/dogmatic"
)

func main() {
	args := godgt.GetParsedArguments()
	app := godgt.NewDgtApp(args)
	app.Run()
}
