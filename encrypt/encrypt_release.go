//+build !debug

package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
)

func (e *Encryptor) encryptFileWithGCM(data []byte) (out []byte, err error) {
	block, err := aes.NewCipher(e.symmetricKey)
	if err != nil {
		return
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	// Create encrypt text
	var cipherSuite = uint8(1)
	var fileLength = uint32(len(data))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, []byte{'M', 'E', 'G', 'V'}) // 4B
	binary.Write(buf, binary.LittleEndian, cipherSuite)                // 1B
	binary.Write(buf, binary.LittleEndian, e.encKey)                   // 256B
	binary.Write(buf, binary.LittleEndian, fileLength)                 // 4B
	binary.Write(buf, binary.LittleEndian, nonce)                      // 12B
	binary.Write(buf, binary.LittleEndian, ciphertext)                 // xxxB

	out = buf.Bytes()
	return
}

func (e *Encryptor) EncryptData(data []byte) (out []byte, err error) {
	if e.version == 1 {
		return e.encryptFileWithGCM(data)
	}
	return nil, fmt.Errorf("Not supported")
}

//-----------------CBC encrypter--------------------//

// GetEncryptStreamReader return a reader
func (e *Encryptor) GetEncryptStreamReader(r io.Reader, fileLength uint32) (EncryptStreamReader, error) {
	if e.version == 2 {
		reader, err := e.createCBCReaderVersion2(r, fileLength)
		e.streamReader = reader
		return reader, err
	} else if e.version == 3 {
		reader, err := e.createCBCReaderVersion3(r, fileLength)
		e.streamReader = reader
		return reader, err
	}
	return nil, fmt.Errorf("Not supported")
}

func (e *Encryptor) createCBCReaderVersion2(r io.Reader, fileLength uint32) (*cbcReader, error) {
	block, err := aes.NewCipher(e.symmetricKey)
	if err != nil {
		return nil, err
	}

	// Create header
	var cipherSuite = uint8(2)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	aescbc := cipher.NewCBCEncrypter(block, iv)

	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, []byte{'M', 'E', 'G', 'V'}) // 4B
	binary.Write(header, binary.LittleEndian, cipherSuite)                // 1B
	binary.Write(header, binary.LittleEndian, e.encKey)                   // 256B
	binary.Write(header, binary.LittleEndian, fileLength)                 // 4B
	binary.Write(header, binary.LittleEndian, iv)                         // 16B

	// fmt.Printf("cipherSuite:%x, key:%x, fileLength:%d, iv:%x\n", cipherSuite, e.symmetricKey, fileLength, iv)

	return &cbcReader{
		r:            r,
		blockMode:    aescbc,
		headerReader: bytes.NewReader(header.Bytes()),
		buf:          new(bytes.Buffer),
		tmp:          make([]byte, aescbc.BlockSize()),
	}, nil
}

func (e *Encryptor) createCBCReaderVersion3(r io.Reader, fileLength uint32) (*cbcReader, error) {
	block, err := aes.NewCipher(e.symmetricKey)
	if err != nil {
		return nil, err
	}

	// Create header
	var cipherSuite = uint8(3)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	aescbc := cipher.NewCBCEncrypter(block, iv)

	header := new(bytes.Buffer)
	binary.Write(header, binary.LittleEndian, []byte{'M', 'E', 'G', 'V'}) // 4B
	binary.Write(header, binary.LittleEndian, cipherSuite)                // 1B
	binary.Write(header, binary.LittleEndian, e.pubKeyID)                 // 8B
	binary.Write(header, binary.LittleEndian, e.encKey)                   // 256B
	binary.Write(header, binary.LittleEndian, fileLength)                 // 4B

	// 计算header的HMAC
	hmacMaker := hmac.New(md5.New, e.symmetricKey)
	hmacMaker.Write(header.Bytes())
	hmacValue := hmacMaker.Sum(nil)

	binary.Write(header, binary.LittleEndian, hmacValue) // 16B
	// fmt.Println("header base64", len(header.Bytes()), base64.StdEncoding.EncodeToString(header.Bytes()))
	binary.Write(header, binary.LittleEndian, iv) // 16B

	// fmt.Printf("cipherSuite:%x, key_id:%d, dec_key:%x, enc_key:%x, fileLength:%d, hmac:%x, iv:%x\n", cipherSuite, e.pubKeyID, e.symmetricKey, e.encKey, fileLength, hmacValue, iv)

	return &cbcReader{
		r:            r,
		blockMode:    aescbc,
		headerReader: bytes.NewReader(header.Bytes()),
		buf:          new(bytes.Buffer),
		tmp:          make([]byte, aescbc.BlockSize()),
	}, nil
}

func (cr *cbcReader) Read(b []byte) (int, error) {
	if cr.err != nil {
		return 0, cr.err
	}

	if cr.headerReader.Len() > 0 {
		n, err := cr.headerReader.Read(b)
		cr.err = err
		return n, err
	}

	blockCount := len(b) / cr.blockMode.BlockSize()
	cr.buf.Reset()
	for i := 0; i < blockCount; i++ {
		num, err := io.ReadFull(cr.r, cr.tmp)
		if num > 0 {
			cr.blockMode.CryptBlocks(cr.tmp, cr.tmp)
			cr.buf.Write(cr.tmp)
		}
		if err == io.EOF {
			cr.err = io.EOF
			break
		}
	}
	defer cr.buf.Reset()

	return copy(b, cr.buf.Bytes()), cr.err
}
