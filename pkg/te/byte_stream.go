package te

import "unicode/utf8"

// ByteStream wraps Stream to accept raw bytes.
type ByteStream struct {
	*Stream
	buffer []byte
}

// NewByteStream creates a ByteStream for the provided handler.
func NewByteStream(screen EventHandler, strict bool) *ByteStream {
	return &ByteStream{Stream: NewStream(screen, strict)}
}

// Feed processes terminal input.
func (st *ByteStream) Feed(data []byte) error {
	if st.useUTF8 {
		st.buffer = append(st.buffer, data...)
		var out []rune
		for len(st.buffer) > 0 {
			r, size := utf8.DecodeRune(st.buffer)
			if r == utf8.RuneError && size == 1 {
				out = append(out, rune(st.buffer[0]))
				st.buffer = st.buffer[1:]
				continue
			}
			st.buffer = st.buffer[size:]
			out = append(out, r)
		}
		return st.Stream.Feed(string(out))
	}
	out := make([]rune, len(data))
	for i, b := range data {
		out[i] = rune(b)
	}
	return st.Stream.Feed(string(out))
}

// SelectOtherCharset switches to the specified character set.
func (st *ByteStream) SelectOtherCharset(code string) {
	st.Stream.SelectOtherCharset(code)
	if st.useUTF8 {
		st.buffer = nil
	}
}
