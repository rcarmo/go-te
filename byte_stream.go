package te

import "unicode/utf8"

type ByteStream struct {
	stream   *Stream
	strict   bool
	buffer   []byte
	encoding string
}

func NewByteStream(screen ScreenLike, strict bool) *ByteStream {
	return &ByteStream{
		stream:   NewStream(screen, strict),
		strict:   strict,
		encoding: "UTF-8",
	}
}

func (st *ByteStream) Feed(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	st.buffer = append(st.buffer, data...)
	for len(st.buffer) > 0 {
		r, size := utf8.DecodeRune(st.buffer)
		if r == utf8.RuneError && size == 1 {
			break
		}
		st.buffer = st.buffer[size:]
		if err := st.stream.FeedString(string(r)); err != nil {
			return err
		}
	}
	return nil
}

func (st *ByteStream) SetEncoding(name string) {
	st.encoding = name
}
