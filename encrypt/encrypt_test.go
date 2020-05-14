package encrypt

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const privatePem = "test_private.pem"
const publicPem = "test_public.pem"
const testFile = "test_file"

func createTestFile(lines int) (filename string, err error) {
	f, err := os.Create(testFile)
	if err != nil {
		return
	}
	defer f.Close()

	var data = "Hello, world!"
	w := bufio.NewWriter(f)
	for i := 0; i < lines; i++ {
		_, err = fmt.Fprintln(w, data)
	}
	w.Flush()

	filename = testFile
	return
}

func setup() (pub, priv string) {
	createTestFile(100)
	GenPrivateKeyPem(2048, privatePem)
	GenPublicKeyPem(privatePem, publicPem)
	return
}

func teardown() {
	os.Remove(publicPem)
	os.Remove(privatePem)
	os.Remove(testFile)
}

func GenPrivateKeyPem(bits int, privatePemPath string) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(privatePemPath)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

func GenPublicKeyPem(privatePemPath, publicPemPath string) error {
	privateContent, err := ioutil.ReadFile(privatePemPath)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(privateContent)
	if block == nil {
		return fmt.Errorf("private key error!")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	file, err := os.Create(publicPemPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
func TestVersion1(t *testing.T) {
	setup()
	defer teardown()

	enc := NewEncryptor(nil, 1)

	data, _ := ioutil.ReadFile(testFile)

	out, err := enc.EncryptData(data)
	if err != nil {
		t.Error(err)
	}

	f, err := os.Create("test_version1.meg")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	f.Write(out)
}

func TestVersion2(t *testing.T) {
	setup()
	defer teardown()

	filename := testFile

	enc := NewEncryptor(nil, 2)
	fi, _ := os.Stat(filename)
	f, _ := os.Open(filename)

	encReader, err := enc.GetEncryptStreamReader(f, uint32(fi.Size()))
	if err != nil {
		t.Error(err)
	}

	out, err := os.Create("test_version2.meg")
	if err != nil {
		t.Error(err)
	}
	defer out.Close()
	io.Copy(out, encReader)
}

func TestVersion3(t *testing.T) {
	setup()
	defer teardown()

	pubk, _ := ioutil.ReadFile("test_data/5097750102325959.pem")

	// 加密测试
	enc := NewEncryptorV3(5097750102325959, pubk)
	fi, _ := os.Stat(testFile)
	f, _ := os.Open(testFile)

	encReader, err := enc.GetEncryptStreamReader(f, uint32(fi.Size()))
	if err != nil {
		t.Error(err)
	}

	out, err := os.Create("test_version3.meg")
	if err != nil {
		t.Error(err)
	}
	defer out.Close()
	io.Copy(out, encReader)
}
