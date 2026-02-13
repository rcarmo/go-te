package te

import "unicode/utf8"

type ByteStream struct {
	*Stream
	buffer []byte
}

func NewByteStream(screen EventHandler, strict bool) *ByteStream {
	return &ByteStream{Stream: NewStream(screen, strict)}
}

func (st *ByteStream) Feed(data []byte) error {
	if st.useUTF8 {
		st.buffer = append(st.buffer, data...)
		var out []rune
		for len(st.buffer) > 0 {
			r, size := utf8.DecodeRune(st.buffer)
			if r == utf8.RuneError && size == 1 {
				break
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

func (st *ByteStream) SelectOtherCharset(code string) {
	st.Stream.SelectOtherCharset(code)
	if st.useUTF8 {
		st.buffer = nil
	}
}
