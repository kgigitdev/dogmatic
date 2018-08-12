package godgt

import (
	"fmt"

	"github.com/notnil/chess"
)

type FieldUpdate struct {
	square chess.Square
	piece  chess.Piece
}

// BoardMessage represents a single, discrete, message from the board.
// These messages have already been translated from their low-level
// byte representation into a higher-level form that external code will
// find easier to handle.
//
// The board is capable of generating multiple message types; at the
// moment we only handle two message types: field updates and board
// updates. Each newly-handled type should probably be added as a new
// field.
type BoardMessage struct {
	boardDumpFEN     string
	fieldUpdate      *FieldUpdate
	unhandledMessage string
}

// The board dump squares are listed in the same order as in a FEN
// diagram; that is, scanning from left to right and down the board
// (a8-h8, a7-h7, ..., a1-h1). As such, squareIndex 0 is a8 and
// squareIndex 63 is h1.  This isn't directly compatible with
// chess.Square, and getSquare() is private.  Luckily, we can just
// reimplement it.
func getSquareFromIndex(squareIndex int) chess.Square {
	fileIndex := squareIndex % 8
	rankIndex := 7 - ((squareIndex - fileIndex) / 8)
	return getSquareFromRankAndFile(rankIndex, fileIndex)
}

func getSquareFromRankAndFile(rank, file int) chess.Square {
	return chess.Square((int(rank) * 8) + int(file))
}

func NewBoardUpdateMessage(bytes []byte) *BoardMessage {
	fenBuilder := NewFenBuilder()
	for _, pieceCode := range bytes {
		// square := getSquareFromIndex(squareIndex)
		// chessPiece := getChessPiece(pieceCode)
		fenChar := getFENChar(pieceCode)
		fenBuilder.Add(fenChar)
	}

	// Hack: a complete FEN also has other metadata like
	// side to move, move number, etc, so we add this on
	// to make our FEN parse.
	fenString := fenBuilder.String()

	return &BoardMessage{
		boardDumpFEN: fenString,
	}
}

func NewUnhandledMessage(b byte) *BoardMessage {
	return &BoardMessage{
		unhandledMessage: fmt.Sprintf("Unhandled message %02x", b),
	}
}

func NewFieldUpdateMessage(bytes []byte) *BoardMessage {
	fieldNumber := bytes[0]
	pieceCode := bytes[1]

	// The field number is encoded as follows:
	// 0b00rrrfff
	// The 3 "rrr" bits denote the rank (8 values in total).
	// The 3 "fff" bits denote the file (8 values in total).
	// The "rrr" bits can be masked using 0x07.
	// The "fff" bits need to be masked using 0b00111000, or
	// 0x38, then right-shifted 3 bits. Note that r=0 actually
	// corresponds to the 8th rank (from White's POV), so we need
	// to take this into account when computing rankIndex by "flipping"
	// the order.

	// The file index, 0=a, 7=h
	fileIndex := int(fieldNumber & 0x07)

	// The rank index, 0=1st rank, 7=8th rank (notice the flip)
	rankIndex := 7 - int((fieldNumber&0x38)>>3)

	square := getSquareFromRankAndFile(rankIndex, fileIndex)
	piece := getChessPiece(pieceCode)

	return &BoardMessage{
		fieldUpdate: &FieldUpdate{
			square: square,
			piece:  piece,
		},
	}
}

func getChessPiece(pieceCode byte) chess.Piece {
	switch pieceCode {
	case WPAWN:
		return chess.WhitePawn
	case WKNIGHT:
		return chess.WhiteKnight
	case WBISHOP:
		return chess.WhiteBishop
	case WROOK:
		return chess.WhiteRook
	case WQUEEN:
		return chess.WhiteQueen
	case WKING:
		return chess.WhiteKing
	case BPAWN:
		return chess.BlackPawn
	case BKNIGHT:
		return chess.BlackKnight
	case BBISHOP:
		return chess.BlackBishop
	case BROOK:
		return chess.BlackRook
	case BQUEEN:
		return chess.BlackQueen
	case BKING:
		return chess.BlackKing
	case EMPTY:
		return chess.NoPiece
	default:
		return chess.NoPiece
	}
}

func getFENChar(pieceCode byte) string {
	switch pieceCode {
	case WPAWN:
		return "P"
	case WKNIGHT:
		return "N"
	case WBISHOP:
		return "B"
	case WROOK:
		return "R"
	case WQUEEN:
		return "Q"
	case WKING:
		return "K"
	case BPAWN:
		return "p"
	case BKNIGHT:
		return "n"
	case BBISHOP:
		return "b"
	case BROOK:
		return "r"
	case BQUEEN:
		return "q"
	case BKING:
		return "k"
	case EMPTY:
		return ""
	default:
		return ""
	}

}
