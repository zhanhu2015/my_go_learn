//+build debug

package encrypt

import (
	"io"
)

func (e *Encryptor) EncryptData(data []byte) (out []byte, err error) {
	return data[:], nil
}

func (e *Encryptor) GetEncryptStreamReader(r io.Reader, fileLength uint32) (EncryptStreamReader, error) {
	return &cbcReader{
		r:            r,
		blockMode:    nil,
		headerReader: nil,
		buf:          nil,
		tmp:          nil,
	}, nil
}

func (cr *cbcReader) Read(b []byte) (int, error) {
	return cr.r.Read(b)
}
