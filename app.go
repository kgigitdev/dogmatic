package godgt

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"github.com/notnil/chess"
)

// DgtApp is the main class-like struct for the application.
type DgtApp struct {
	args     *DgtAppArgs
	startFEN string
	port     io.ReadWriteCloser
	board    *DgtBoard
	currGame *chess.Game
}

func NewDgtApp(args *DgtAppArgs) *DgtApp {
	return &DgtApp{
		args: args,
	}
}

// Run is the main entry point for DgtApp.
func (a *DgtApp) Run() {
	a.getStartFEN()
	a.openPort()
	a.createBoard()
	a.initialiseBoard()
	a.runForever()
}

func (a *DgtApp) getStartFEN() {
	startFEN := chess.NewGame().Position().String()
	a.startFEN = strings.Replace(startFEN, " w KQkq - 0 1", "", -1)

}

func (a *DgtApp) openPort() {
	a.log("Opening port ...\n")
	options := serial.OpenOptions{
		PortName:        a.args.Device,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	port, err := serial.Open(options)
	a.check(err)
	a.port = port
}

func (a *DgtApp) createBoard() {
	a.log("Creating board ...\n")
	a.board = NewDgtBoard(a.port)
}

func (a *DgtApp) initialiseBoard() {
	a.log("Starting byte reader ...\n")
	a.board.startByteReader()
	a.log("Starting command processor ...\n")
	a.board.startCommandProcessor()
	a.log("Resetting board ...\n")
	_, err := a.board.WriteSendResetCommand()
	a.check(err)
	_, err = a.board.WriteSendBoardCommand()
	a.check(err)
	_, err = a.board.WriteSendUpdateBoardCommand()
	a.check(err)
}

func (a *DgtApp) runForever() {
	for {
		select {
		case bm := <-a.board.GetBoardMessageChannel():
			a.processBoardMessage(bm)
		}
	}
}

func (a *DgtApp) processBoardMessage(bm *BoardMessage) {
	if bm.boardDumpFEN != "" {
		a.handleBoardDumpFEN(bm.boardDumpFEN)
	} else if bm.fieldUpdate != nil {
		a.handleFieldUpdate(bm.fieldUpdate)
	} else {
		a.log(bm.unhandledMessage + "\n")
	}
}

func (a *DgtApp) handleBoardDumpFEN(boardDumpFEN string) {
	// log.Println("In handleBoardDumpFEN()")

	// We get a board dump at the beginning. If we haven't
	// started the game yet and the FEN is the starting FEN,
	// start the game.

	if a.currGame == nil {
		if boardDumpFEN == a.startFEN {
			log.Print("New Game Started!")
			a.currGame = chess.NewGame()
		} else {
			a.log("Pieces not in position yet.\n")
		}
		return
	}

	currPos := a.currGame.Position()
	validMoves := a.currGame.ValidMoves()

	for _, candidateMove := range validMoves {
		candidatePos := currPos.Update(candidateMove)
		candidateFEN := candidatePos.String()
		miniFEN := strings.Fields(candidateFEN)[0]

		if miniFEN == boardDumpFEN {
			// We've matched a valid move!
			err := a.currGame.Move(candidateMove)
			if err != nil {
				a.log(err.Error())
				return
			}
			a.log("Accepted move: %s\n", candidateMove)
			a.printGameInfo()
			return
		}
	}

	a.log("Not a move.\n")
	fmt.Println(a.currGame.Position().String())
	// Print out the diffs to help backtrack.
	a.printBoardDiffs(boardDumpFEN)
	// a.printGameInfo()
}

func (a *DgtApp) printBoardDiffs(boardDumpFEN string) {
	logicalBoard := a.currGame.Position().Board()
	physicalFEN, err := chess.FEN(boardDumpFEN + " w KQkq - 0 1")
	if err != nil {
		a.log("Board is in illegal position.\n")
		return
	}
	physicalGame := chess.NewGame(physicalFEN)
	physicalBoard := physicalGame.Position().Board()

	for iRank := 7; iRank >= 0; iRank-- {
		for iFile := 0; iFile <= 7; iFile++ {
			// Should be a function to get this
			square := chess.Square((int(iRank) * 8) + int(iFile))

			logicalPiece := logicalBoard.Piece(square)
			physicalPiece := physicalBoard.Piece(square)

			if logicalPiece != physicalPiece {
				a.log("In %s, want %s have %s\n",
					square, logicalPiece, physicalPiece)
			}
		}
	}
}

func (a *DgtApp) printGameInfo() {
	fmt.Println(a.currGame.Position().String())
	fmt.Println(a.currGame)
	fmt.Println(a.currGame.Position().Board().Draw())
}

func (a *DgtApp) handleFieldUpdate(fieldUpdate *FieldUpdate) {
	// Every time we receive a field update, ask for a complete
	// board update. This is simple but inefficient, and will do
	// for the time being.
	//
	// Note that it's tempting to limit such requests only to
	// "piece drop" messages, but that's actually incorrect, for
	// at least two reasons. Firstly, we can receive the piece
	// drop and piece lift messages out of order. Secondly, a
	// capture executed as "use capturing piece to push captured
	// piece off the square, then remove captured piece" might end
	// with a final piece lift message.
	// log.Printf("Field update %s to %s\n", fieldUpdate.piece,
	//	fieldUpdate.square)
	_, err := a.board.WriteSendBoardCommand()
	a.check(err)
}

func (a *DgtApp) check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (a *DgtApp) log(format string, args ...interface{}) {
	fmt.Printf(format, args)
}
