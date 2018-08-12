package godgt

import (
	"io"
	"log"
)

// DgtBoard is a thin wrapper around the serial connection to a DGT
// board. It runs two goroutines, where each goroutine transforms
// messages from the board into higher-level messages.
//
// readBytesFromBoard() reads raw bytes and pushes them into a byte
// channel.
//
// processCommands() reads bytes from the byte channel and chunks them
// up into messages, writing them into a *BoardMessage channel.
//
// External code then reads messages from this channel.
type DgtBoard struct {
	port              io.ReadWriteCloser
	bytesFromBoard    chan byte
	messagesFromBoard chan *BoardMessage
}

func NewDgtBoard(port io.ReadWriteCloser) *DgtBoard {
	return &DgtBoard{
		port:              port,
		bytesFromBoard:    make(chan byte, 2048),
		messagesFromBoard: make(chan *BoardMessage, 2048),
	}
}

func (b *DgtBoard) WriteBytes(bytes []byte) (int, error) {
	return b.port.Write(bytes)
}

func (b *DgtBoard) WriteByte(B byte) (int, error) {
	bytes := []byte{B}
	return b.WriteBytes(bytes)
}

func (b *DgtBoard) WriteSendResetCommand() (int, error) {
	return b.WriteByte(DGT_SEND_RESET)
}

func (b *DgtBoard) WriteSendBoardCommand() (int, error) {
	return b.WriteByte(DGT_SEND_BRD)
}

func (b *DgtBoard) WriteSendUpdateBoardCommand() (int, error) {
	return b.WriteByte(DGT_SEND_UPDATE_BRD)
}

func (b *DgtBoard) startByteReader() {
	go b.readBytesFromBoard()
}

func (b *DgtBoard) readBytesFromBoard() {
	buf := make([]byte, 2048)
	for {
		n, err := b.port.Read(buf)

		if err != nil {
			log.Println("error reading bytes.")
			// Should we worry more? Should we reset the board?
			continue
		}

		if n == 0 {
			continue
		}

		for i := 0; i < n; i++ {
			oneByte := buf[i]
			b.bytesFromBoard <- oneByte
		}
	}
}

func (b *DgtBoard) startCommandProcessor() {
	go b.processCommands()
}

func (b *DgtBoard) processCommands() {
	var bytes []byte
	for {
		select {
		case B := <-b.bytesFromBoard:
			bytes = append(bytes, B)
			message, remainder := b.extractMessage(bytes)
			bytes = remainder
			if message != nil {
				b.messagesFromBoard <- message
			}
		default:
			// Do nothing.
		}
	}
}

func (b *DgtBoard) extractMessage(bytes []byte) (*BoardMessage, []byte) {
	return NewMessageExtractor(bytes).Extract()
}

func (b *DgtBoard) GetBoardMessageChannel() <-chan *BoardMessage {
	return b.messagesFromBoard
}
