package godgt

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type DgtAppArgs struct {
	Device string `short:"d" long:"device" description:"Board device (e.g. /dev/ttyUSB0, /dev/ttyACM0)" default:"/dev/ttyACM0"`

	Pargs []string
}

func GetParsedArguments() *DgtAppArgs {
	var args DgtAppArgs
	parser := flags.NewParser(&args, flags.Default)
	pargs, err := parser.ParseArgs(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	args.Pargs = pargs

	return &args
}
