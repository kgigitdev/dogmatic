package godgt

type MessageExtractor struct {
	bytes         []byte
	b0            byte
	b1            byte
	b2            byte
	m0            byte
	m1            byte
	m2            byte
	messageLength int
	removedBytes  []byte
	boardMessage  *BoardMessage
}

func NewMessageExtractor(bytes []byte) *MessageExtractor {
	return &MessageExtractor{
		bytes: bytes,
	}
}

func (m *MessageExtractor) Extract() (*BoardMessage, []byte) {
	// A  valid message is always at least 3 bytes, so
	// anything less than this can be ignored.
	if m.isTooShort() {
		return m.removeFirstNBytes(0)
	}

	m.extractFirstThreeBytes()

	if !m.isValidMessage() {
		return m.removeFirstNBytes(1)
	}

	m.computeMaskedMessageBytes()
	m.computeMessageLength()

	if !m.haveEnoughBytes() {
		return m.removeFirstNBytes(0)
	}

	if m.m0 == DGT_NONE {
		// The NONE command consists of nothing beyond the
		// 3-byte header, so remove this command entirely.
		return m.removeFirstNBytes(3)
	}

	// If we reach this point, we know that we have both (a) a
	// valid (non-NONE) command, and (b) all of the bytes that
	// make up this command. We can therefore remove ALL the
	// bytes making up this command from m.bytes. We don't need
	// them to parse the actual command type because we already
	// have them in m.m{012}.

	m.removeFirstNBytes(3)
	m.removeFirstNBytes(m.messageLength - 3)
	m.createNewBoardMessage()
	return m.boardMessage, m.bytes

}

func (m *MessageExtractor) isTooShort() bool {
	return len(m.bytes) < 3
}

func (m *MessageExtractor) extractFirstThreeBytes() {
	m.b0 = m.bytes[0]
	m.b1 = m.bytes[1]
	m.b2 = m.bytes[2]
}

func (m *MessageExtractor) isValidMessage() bool {
	// All valid messages start with the first byte
	// having its message bit set. Therefore, whenever
	// our first byte does not have the message bit
	// set, it's a corrupt message header. We should
	// therefore pop that byte off and return the
	// remainder in the hope of resynchronising.
	return m.b0|MESSAGE_BIT != 0
}

// removeFirstNBytes is a multi-purpose method. It returns a 2-tuple
// of (nil, []byte) for the common case where we want to return early
// from the messaage parsing without a command, with the second return
// value being the result of m.bytes *after* the required number of
// bytes have been removed. Finally, the bytes that are removed are
// also stored in the field m.removedBytes.
func (m *MessageExtractor) removeFirstNBytes(n int) (*BoardMessage, []byte) {
	m.removedBytes = m.bytes[0:n]
	m.bytes = m.bytes[n:]
	return nil, m.bytes
}

func (m *MessageExtractor) computeMaskedMessageBytes() {
	// All bytes in the message need to be parsed with their MSB
	// zeroed out, meaning that the individual bytes lie in the
	// range 0-127.
	m.m0 = m.b0 & MESSAGE_MASK
	m.m1 = m.b1 & MESSAGE_MASK
	m.m2 = m.b2 & MESSAGE_MASK
}

func (m *MessageExtractor) computeMessageLength() {
	// Combine the second and third bytes (that is, b1 b2) to
	// determine the message length. Since each byte now only
	// contains 7 bits, we left shift m1 by 7, not 8. Note
	// that the maximum message length is therefore 14 bits
	// (not 16).
	m.messageLength = (int(m.m1) << 7) + int(m.m2)
}

func (m *MessageExtractor) haveEnoughBytes() bool {
	return len(m.bytes) >= m.messageLength
}

func (m *MessageExtractor) createNewBoardMessage() {
	switch m.m0 {
	case DGT_BOARD_DUMP:
		m.boardMessage = NewBoardUpdateMessage(m.removedBytes)
	case DGT_FIELD_UPDATE:
		m.boardMessage = NewFieldUpdateMessage(m.removedBytes)
	default:
		m.boardMessage = NewUnhandledMessage(m.m0)
	}
}
