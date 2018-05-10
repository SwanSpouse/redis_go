package rdb

import (
	"encoding/binary"
	"fmt"
	"github.com/cupcake/rdb/crc64"
	"hash"
	"io"
	"math"
	"os"
	"strconv"
)

const Version = 6

type Encoder struct {
	w   io.Writer
	crc hash.Hash
}

func NewEncoder(filename string) (*Encoder, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &Encoder{w: io.MultiWriter(f), crc: crc64.New()}, nil
}

func (e *Encoder) EncodeHeader() error {
	_, err := fmt.Fprintf(e.w, "REDIS%04d", Version)
	return err
}

func (e *Encoder) EncodeFooter() error {
	e.w.Write([]byte{rdbFlagEOF})
	_, err := e.w.Write(e.crc.Sum(nil))
	return err
}

func (e *Encoder) EncodeDatabase(n int) error {
	e.w.Write([]byte{rdbFlagSelectDB})
	return e.EncodeLength(uint32(n))
}

func (e *Encoder) EncodeLength(length uint32) (err error) {
	switch {
	case length < 1<<6:
		// 如果6个bit能够存储下，就是用6个bit
		_, err = e.w.Write([]byte{byte(length)})
	case length < 1<<14:
		// 如果14个bit能够存储下，就是用14个bit
		_, err = e.w.Write([]byte{byte(length>>8) | rdb14bitLen<<6, byte(length)})
	default:
		// 否则是用4个字节来进行存储
		b := make([]byte, 5)
		b[0] = rdb32bitLen << 6
		binary.BigEndian.PutUint32(b[1:], length)
		_, err = e.w.Write(b)
	}
	return
}

func (e *Encoder) EncodeType(typo ValueType) error {
	_, err := e.w.Write([]byte{byte(typo)})
	return err
}

func (e *Encoder) EncodeRawString(s string) error {
	return e.EncodeString([]byte(s))
}

func (e *Encoder) EncodeString(s []byte) error {
	// TODO lmj redis object string的编码方式现在只有string一种。没有int
	//written, err := e.encodeIntString(s)
	//if written {
	//	return err
	//}
	e.EncodeLength(uint32(len(s)))
	_, err := e.w.Write(s)
	return err
}

func (e *Encoder) encodeIntString(b []byte) (written bool, err error) {
	s := string(b)
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return
	}
	// if the stringified parsed int isn't exactly the same, we can't encode it as an int
	if s != strconv.FormatInt(i, 10) {
		return
	}
	switch {
	case i >= math.MinInt8 && i <= math.MaxInt8:
		_, err = e.w.Write([]byte{rdbEncodingVal << 6, byte(int8(i))})
	case i >= math.MinInt16 && i <= math.MaxInt16:
		b := make([]byte, 3)
		b[0] = rdbEncodingVal<<6 | rdbEncInt16
		binary.LittleEndian.PutUint16(b[1:], uint16(int16(i)))
		_, err = e.w.Write(b)
	case i >= math.MinInt32 && i <= math.MaxInt32:
		b := make([]byte, 5)
		b[0] = rdbEncodingVal<<6 | rdbEncInt32
		binary.LittleEndian.PutUint32(b[1:], uint32(int32(i)))
		_, err = e.w.Write(b)
	default:
		return
	}
	return true, err
}
