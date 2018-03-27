package tcp

type buffer []byte

// remove CRLF
func (buf buffer) TrimCRLF() buffer {
	n := len(buf)
	for ; n > 0; n -- {
		if c := buf[n-1]; c != '\r' && c != '\n' {
			break
		}
	}
	return buf[:n]
}

// return the first world
func (buf buffer) FirstWord() string {
	offset := 0
	inWord := false
	data := buf.TrimCRLF()

	for i, c := range data {
		switch c {
		case ' ', '\t':
			if inWord {
				return string(data[offset:i])
			}
			inWord = false
		default:
			if !inWord {
				offset = i
			}
			inWord = true
		}
	}
	return string(data[offset:])
}

// ParseInt parses an int
func (buf buffer) ParseInt() (int64, error) {
	data := buf.TrimCRLF()
	if len(data) < 2 {
		return 0, ProtoErrorf("Protocol error: expected ':', got ' '")
	} else if data[0] != ':' {
		return 0, ProtoErrorf("Protocol error: expected ':', got '%s'", string(data[0]))
	}

	n, m := int64(0), int64(1)
	for i, c := range data[1:] {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		} else if c == '-' && i == 0 {
			m = -1
		} else {
			return 0, ErrNotANumber
		}
	}
	return n * m, nil
}

// converts the line to a string
func (buf buffer) ParseMessage(prefix byte) (string, error) {
	data := buf.TrimCRLF()
	if len(data) < 1 {
		return "", ProtoErrorf("Protocol error: expected '%s', got ' '", string(prefix))
	} else if data[0] != prefix {
		return "", ProtoErrorf("Protocol error: expected '%s', got '%s'", string(prefix), string(data[0]))
	}
	return string(data[1:]), nil
}

// ParseSize parses a size with prefix
func (buf buffer) ParseSize(prefix byte, fallback error) (int64, error) {
	data := buf.TrimCRLF()

	if len(data) == 0 {
		return 0, ProtoErrorf("Protocol error: expected '%s', got ' '", string(prefix))
	} else if data[0] != prefix {
		return 0, ProtoErrorf("Protocol error: expected '%s', got '%s'", string(prefix), string(data[0]))
	} else if len(data) < 2 {
		return 0, fallback
	}

	var n int64
	for _, c := range data[1:] {
		if c >= '0' && c <= '9' {
			n = n*10 + int64(c-'0')
		} else {
			return 0, fallback
		}
	}
	if n < 0 {
		return 0, fallback
	}
	return n, nil
}
