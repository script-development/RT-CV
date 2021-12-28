package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"hash"
	"io"

	"github.com/script-development/RT-CV/helpers/numbers"
)

// chunkSize is the max size of every encrypted chunk
const chunkSize = 8192

// EncryptWriter encrypts files using the io.Writer interface and writes the output directly to another io.Writer
// This is handy to use when you want to encrypt a large amount of data
type EncryptWriter struct {
	closed            bool
	gcm               cipher.AEAD
	nonce             []byte
	dst               io.Writer
	buff              []byte
	encryptedDataBuff []byte
	nonceHasher       hash.Hash
}

// NewEncryptWriter creates a new instance of the encrypt writer
func NewEncryptWriter(key []byte, dst io.Writer) (*EncryptWriter, error) {
	c, err := aes.NewCipher(NormalizeKey(key))
	if err != nil {
		return nil, err
	}

	ew := &EncryptWriter{
		closed:            false,
		dst:               dst,
		buff:              []byte{},
		encryptedDataBuff: []byte{},
		nonceHasher:       sha512.New(),
	}

	ew.gcm, err = cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	ew.nonce = make([]byte, ew.gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, ew.nonce)
	if err != nil {
		return nil, err
	}

	err = ew.mustWriteToDst(ew.nonce)
	if err != nil {
		return nil, err
	}

	return ew, nil
}

// Write Implements io.Writer
func (ew *EncryptWriter) Write(p []byte) (n int, err error) {
	if ew.closed {
		return 0, errors.New("writer is closed")
	}

	ew.buff = append(ew.buff, p...)
	return len(p), ew.tryWriterChunk()
}

// Close Implements io.Closer
// This writes the remaining data to the underlying writer
func (ew *EncryptWriter) Close() error {
	if ew.closed {
		return errors.New("cannot close a already closed writer is closed")
	}

	defer func() {
		ew.closed = true
		ew.gcm = nil
		ew.nonce = nil
		// ew.dst
		ew.buff = nil
		ew.encryptedDataBuff = nil
		ew.nonceHasher = nil
	}()

	err := ew.tryWriterChunk()
	if err != nil {
		return err
	}

	if len(ew.buff) > 0 {
		// Write the remaining data
		ew.encryptedDataBuff = ew.gcm.Seal(ew.encryptedDataBuff[:0], ew.nonce, ew.buff, nil)
		encryptedDataBuffSize := numbers.UintToBytes(uint64(len(ew.encryptedDataBuff)), 32)

		err = ew.mustWriteToDst(encryptedDataBuffSize)
		if err != nil {
			return err
		}
		err = ew.mustWriteToDst(ew.encryptedDataBuff)
	}

	return err
}

// tryWriterChunk tries to write a chunk to the underlying writer if the buffer is large enough
func (ew *EncryptWriter) tryWriterChunk() error {
	for {
		if len(ew.buff) < chunkSize {
			// We do not have enough data to create a chunk yet
			return nil
		}

		ew.encryptedDataBuff = ew.gcm.Seal(ew.encryptedDataBuff[:0], ew.nonce, ew.buff[:chunkSize], nil)
		encryptedDataBuffSize := numbers.UintToBytes(uint64(len(ew.encryptedDataBuff)), 32)

		// Write the size of the chunk
		err := ew.mustWriteToDst(encryptedDataBuffSize)
		if err != nil {
			return err
		}

		// Write the encrypted chunk
		err = ew.mustWriteToDst(ew.encryptedDataBuff)
		if err != nil {
			return err
		}

		ew.genNextNonce()
		ew.buff = ew.buff[chunkSize:]
	}
}

// genNextNonce generates the next nonce in the chain
func (ew *EncryptWriter) genNextNonce() {
	genNextNonce(ew.nonceHasher, ew.nonce)
}

func genNextNonce(hasher hash.Hash, nonce []byte) {
	hasher.Reset()
	hasher.Write(nonce)
	nonce = append(nonce[:0], hasher.Sum(nonce[:0])[:len(nonce)]...)
}

// mustWriteToDst makes sure that ALL data of b is written to the destination
func (ew *EncryptWriter) mustWriteToDst(b []byte) error {
	offset := 0
	for {
		n, err := ew.dst.Write(b[offset:])
		if err != nil {
			return err
		}
		offset += n
		if offset == len(b) {
			return nil
		}
	}
}

// EncryptReader decrypts files using the io.Reader
type EncryptReader struct {
	key             []byte
	nonce           []byte
	chunk           []byte
	chunkReadOffset int
	gcm             cipher.AEAD
	source          io.Reader
	nonceHasher     hash.Hash
	eof             bool
}

// NewEncryptReader creates a new instance of the encrypt reader
// This can read from another reader where the data is encrypted using EncryptWriter
func NewEncryptReader(key []byte, source io.Reader) (*EncryptReader, error) {
	r := &EncryptReader{
		key:         key,
		chunk:       make([]byte, chunkSize),
		source:      source,
		nonceHasher: sha512.New(),
	}

	c, err := aes.NewCipher(NormalizeKey(key))
	if err != nil {
		return nil, err
	}

	r.gcm, err = cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := r.gcm.NonceSize()

	r.nonce, err = r.mustReadFromSource(nonceSize)
	if err != nil {
		return nil, err
	}

	err = r.readNextChunk()
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Read Implements io.Reader
func (r *EncryptReader) Read(p []byte) (n int, err error) {
	if r.eof {
		return 0, io.EOF
	}

	if r.chunkReadOffset+len(p) <= len(r.chunk) {
		// Requested less bytes than the chunk size
		copy(p, r.chunk[r.chunkReadOffset:r.chunkReadOffset+len(p)])
		r.chunkReadOffset += len(p)
		return len(p), nil
	}

	// Requested more bytes than the current chunk can offer

	if r.chunkReadOffset < len(r.chunk) {
		// Lets return the left over chunk data
		copy(p, r.chunk[r.chunkReadOffset:])
		n = len(r.chunk) - r.chunkReadOffset
		r.chunkReadOffset = len(r.chunk)
		return n, nil
	}

	// Read all the data from the last chunk, lets now read the next one
	err = r.readNextChunk()
	if r.eof {
		return 0, io.EOF
	}
	if err != nil {
		return 0, err
	}

	if len(p) < len(r.chunk) {
		// Request range was in range of the new chunk
		copy(p, r.chunk)
		r.chunkReadOffset = len(p)
		return len(p), nil
	}

	// Requested more bytes than the next chunk can offer, lets return what we can return
	copy(p, r.chunk)
	r.chunkReadOffset = len(r.chunk)
	return len(r.chunk), nil
}

// readNextChunk reads the next encrypted chunk and decrypts it
func (r *EncryptReader) readNextChunk() (err error) {
	encryptedChunkSizeBytes, err := r.mustReadFromSource(4)
	if err != nil {
		if err == io.EOF {
			r.eof = true
			return nil
		}
		return err
	}

	encryptedChunkSize, err := numbers.BytesToUint(encryptedChunkSizeBytes)
	if err != nil {
		return err
	}

	encryptedBytes, err := r.mustReadFromSource(int(encryptedChunkSize))
	if err != nil {
		return err
	}

	r.chunk, err = r.gcm.Open(r.chunk[:0], r.nonce, encryptedBytes, nil)
	if err != nil {
		return err
	}
	r.chunkReadOffset = 0

	r.genNextNonce()
	return nil
}

// mustReadFromSource reads exactly n bytes from the source and returns error if they could not be read
func (r *EncryptReader) mustReadFromSource(n int) ([]byte, error) {
	dst := make([]byte, n)
	bytesRead, err := io.ReadFull(r.source, dst)
	if bytesRead == 0 {
		return dst, io.EOF
	}
	return dst, err
}

// genNextNonce generates the next nonce in the chain
func (r *EncryptReader) genNextNonce() {
	genNextNonce(r.nonceHasher, r.nonce)
}
