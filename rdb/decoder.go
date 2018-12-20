package rdb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/SwanSpouse/redis_go/loggers"
)

type ValueType byte

const (
	TypeString ValueType = 0
	TypeList   ValueType = 1
	TypeSet    ValueType = 2
	TypeZSet   ValueType = 3
	TypeHash   ValueType = 4

	TypeHashZipmap    ValueType = 9
	TypeListZiplist   ValueType = 10
	TypeSetIntset     ValueType = 11
	TypeZSetZiplist   ValueType = 12
	TypeHashZiplist   ValueType = 13
	TypeListQuicklist ValueType = 14
)

const (
	rdb6bitLen     = 0
	rdb14bitLen    = 1
	rdb32bitLen    = 2
	rdbEncodingVal = 3

	rdbFlagAux      = 0xfa
	rdbFlagResizeDB = 0xfb
	rdbFlagExpiryMS = 0xfc
	rdbFlagExpiry   = 0xfd
	rdbFlagSelectDB = 0xfe
	rdbFlagEOF      = 0xff

	rdbEncInt8  = 0
	rdbEncInt16 = 1
	rdbEncInt32 = 2
	rdbEncLZF   = 3

	rdbZiplist6bitlenString  = 0
	rdbZiplist14bitlenString = 1
	rdbZiplist32bitlenString = 2

	rdbZiplistInt16 = 0xc0
	rdbZiplistInt32 = 0xd0
	rdbZiplistInt64 = 0xe0
	rdbZiplistInt24 = 0xf0
	rdbZiplistInt8  = 0xfe
	rdbZiplistInt4  = 15
)

type byteReader interface {
	io.Reader
	io.ByteReader
}

// A IDecoder must be implemented to parse a RDB file.
type IDecoder interface {
	// StartRDB is called when parsing of a valid RDB file starts.
	StartRDB()
	// StartDatabase is called when database n starts.
	// Once a database starts, another database will not start until EndDatabase is called.
	StartDatabase(n int)
	// AUX field
	Aux(key, value []byte)
	// ResizeDB hint
	ResizeDatabase(dbSize, expiresSize uint32)
	// Set is called once for each string key.
	Set(key, value []byte, expiry int64)
	// StartHash is called at the beginning of a hash.
	// Hset will be called exactly length times before EndHash.
	StartHash(key []byte, length, expiry int64)
	// Hset is called once for each field=value pair in a hash.
	Hset(key, field, value []byte)
	// EndHash is called when there are no more fields in a hash.
	EndHash(key []byte)
	// StartSet is called at the beginning of a set.
	// Sadd will be called exactly cardinality times before EndSet.
	StartSet(key []byte, cardinality, expiry int64)
	// Sadd is called once for each member of a set.
	Sadd(key, member []byte)
	// EndSet is called when there are no more fields in a set.
	EndSet(key []byte)
	// StartList is called at the beginning of a list.
	// Rpush will be called exactly length times before EndList.
	// If length of the list is not known, then length is -1
	StartList(key []byte, length, expiry int64)
	// Rpush is called once for each value in a list.
	Rpush(key, value []byte)
	// EndList is called when there are no more values in a list.
	EndList(key []byte)
	// StartZSet is called at the beginning of a sorted set.
	// Zadd will be called exactly cardinality times before EndZSet.
	StartZSet(key []byte, cardinality, expiry int64)
	// Zadd is called once for each member of a sorted set.
	Zadd(key []byte, score float64, member []byte)
	// EndZSet is called when there are no more members in a sorted set.
	EndZSet(key []byte)
	// EndDatabase is called at the end of a database.
	EndDatabase(n int)
	// EndRDB is called when parsing of the RDB file is complete.
	EndRDB()
}

type Decoder struct {
	r      byteReader
	server IDecoder
}

func NewDecoder(filename string, event IDecoder) (*Decoder, error) {
	if f, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		return &Decoder{
			r:      bufio.NewReader(f),
			server: event}, nil
	}
}

func (d *Decoder) checkRDBFileHeader() error {
	header := make([]byte, 9)
	_, err := io.ReadFull(d.r, header)
	if err != nil {
		return err
	}
	if !bytes.Equal(header[:5], []byte("REDIS")) {
		return errors.New("rbd: invalid file format")
	}

	version, _ := strconv.ParseInt(string(header[5:]), 10, 64)
	if version < 1 || version > 7 {
		return errors.New(fmt.Sprintf("rdb: invalid RDB version number %d", version))
	}
	return nil
}

func (d *Decoder) Decode() error {
	if err := d.checkRDBFileHeader(); err != nil {
		loggers.Errorf("check rdb file header error %+v", err)
		return err
	}
	d.server.StartRDB()

	firstDB := true
	var db uint32
	var expiry int64
	for {
		objType, err := d.r.ReadByte()
		fmt.Printf("current objType ==>%+v %U %s\n", objType, objType, string(objType))

		if err != nil {
			return err
		}
		switch objType {
		case rdbFlagSelectDB:
			if !firstDB {
				d.server.EndDatabase(int(db))
			}
			db, _, err = d.readLength()
			if err != nil {
				return err
			}
			d.server.StartDatabase(int(db))
		case rdbFlagEOF:
			fmt.Printf("reach EOF\n")
			d.server.EndDatabase(int(db))
			d.server.EndRDB()
			return nil
		case rdbFlagExpiryMS:
			// TODO 接下来将要读入8个字节长度、毫秒为单位的过期时间
		case rdbFlagExpiry:
			// TODO 接下来将要读入4个字节长度、秒为单位的过期时间
		case rdbFlagResizeDB:
			// TODO
		case rdbFlagAux:
			// TODO
		default:
			key, err := d.readString()
			if err != nil {
				return err
			}
			err = d.readObject(key, ValueType(objType), expiry)
			if err != nil {
				return err
			}
			expiry = 0
		}
	}
	panic("should not reached")
}

func (d *Decoder) readString() ([]byte, error) {
	length, isEncoded, err := d.readLength()
	if err != nil {
		return nil, err
	}
	if isEncoded {
		switch length {
		case rdbEncInt8:
			i, err := d.readUint8()
			return []byte(strconv.FormatInt(int64(int8(i)), 10)), err
		case rdbEncInt16:
			i, err := d.readUint16()
			return []byte(strconv.FormatInt(int64(int16(i)), 10)), err
		case rdbEncInt32:
			i, err := d.readUint32()
			return []byte(strconv.FormatInt(int64(int32(i)), 10)), err
		}
	}
	key := make([]byte, length)
	_, err = io.ReadFull(d.r, key)
	return key, err
}

func (d *Decoder) readObject(key []byte, typo ValueType, expiry int64) error {
	switch typo {
	case TypeString:
		value, err := d.readString()
		if err != nil {
			return err
		}
		loggers.Info("we get a string object key:%s, value:%s, expiry:%d", key, value, expiry)
		d.server.Set(key, value, expiry)
	case TypeList:
		length, _, err := d.readLength()
		if err != nil {
			return err
		}
		d.server.StartList(key, int64(length), expiry)
		for i := uint32(0); i < length; i++ {
			value, err := d.readString()
			if err != nil {
				return err
			}
			loggers.Info("we get a list object key:%s, item:%s, expiry:%d", key, value, expiry)
			d.server.Rpush(key, value)
		}
		d.server.EndList(key)
	case TypeSet:
		cardinality, _, err := d.readLength()
		if err != nil {
			return err
		}
		d.server.StartSet(key, int64(cardinality), expiry)
		for i := uint32(0); i < cardinality; i++ {
			member, err := d.readString()
			if err != nil {
				return err
			}
			loggers.Info("we get a set object key:%s, item:%s, expiry:%d", key, member, expiry)
			d.server.Sadd(key, member)
		}
		d.server.EndSet(key)
	case TypeHash:
		length, _, err := d.readLength()
		if err != nil {
			return err
		}
		d.server.StartHash(key, int64(length), expiry)
		for i := uint32(0); i < length; i++ {
			field, err := d.readString()
			if err != nil {
				return err
			}
			value, err := d.readString()
			if err != nil {
				return err
			}
			loggers.Info("we get a hash object key:%s, field:%s, value:%s, expiry:%d", key, field, value, expiry)
			d.server.Hset(key, field, value)
		}
		d.server.EndHash(key)
	case TypeZSet:
		cardinality, _, err := d.readLength()
		if err != nil {
			return err
		}
		d.server.StartZSet(key, int64(cardinality), expiry)
		for i := uint32(0); i < cardinality; i++ {
			member, err := d.readString()
			if err != nil {
				return err
			}
			score, err := d.readFloat64()
			if err != nil {
				return err
			}
			loggers.Info("we get a zset object  key:%s, member:%s, score:%.2f, expiry:%d", key, member, score, expiry)
			d.server.Zadd(key, score, member)
		}
		d.server.EndZSet(key)
	default:
		return fmt.Errorf("rdb: unknown object type %d for key %s", typo, key)
	}
	return nil
}

/*
长度编码用于存储流中接下来对象的长度。长度编码是一个可变字节编码，为尽可能少用字节而设计:
从流中读取一个字节，最高 2 bit 被读取。
	如果开始 bit 是 00，接下来 6 bit 表示长度。
	如果开始 bit 是 01，从流中再读取额外一个字节。这组合的的 14 bit 表示长度。
	如果开始 bit 是 10，那么剩余的 6bit 丢弃，从流中读取额外的 4 字节，这 4 个字节表示长度。
	如果开始 bit 是 11，那么接下来的对象是以特殊格式编码的。剩余 6 bit 指示格式。这种编码通常用于把数字作为字符串存储或存储编码后的字符串。
*/

func (d *Decoder) readLength() (uint32, bool, error) {
	b, err := d.r.ReadByte()
	if err != nil {
		return 0, false, err
	}
	switch (b & 0xc0) >> 6 {
	case rdb6bitLen:
		// when the first two bits are 00, the next 6 bits are the length.
		return uint32(b & 0x3f), false, nil
	case rdb14bitLen:
		// when the first two bits are 01, the next 14 bits are the length.
		nextByte, err := d.r.ReadByte()
		if err != nil {
			return 0, false, err
		}
		return (uint32(b&0x3f) << 8) | uint32(nextByte), false, nil
	case rdb32bitLen:
		// when the first two bits are 10, the next 6 bits are discarded.
		// The next 4 bytes are the length
		length, err := d.readUint32()
		return length, false, err
	case rdbEncodingVal:
		// when the first two bits are 11, the next object is encoded.
		// the next 6 bits indicate the encoding type
		return uint32(b & 0x3f), true, nil
	}
	panic("should not reached")
}

func (d *Decoder) readUint8() (uint8, error) {
	b, err := d.r.ReadByte()
	return uint8(b), err
}

func (d *Decoder) readUint16() (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(d.r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func (d *Decoder) readUint32() (uint32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(d.r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

func (d *Decoder) readUint64() (uint64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(d.r, buf)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// Doubles are saved as strings prefixed by an unsigned
// 8 bit integer specifying the length of the representation.
// This 8 bit integer has special values in order to specify the following
// conditions:
// 253: not a number
// 254: + inf
// 255: - inf
// TODO lmj copy from /Users/LiMingji/Documents/go/go_codings/src/github.com/cupcake/rdb/decoder.go
func (d *Decoder) readFloat64() (float64, error) {
	// TODO lmj 这里有个疑问，为啥float64的长度一定是1个字节呢？
	length, err := d.readUint8()
	if err != nil {
		return 0, err
	}
	switch length {
	case 253:
		return math.NaN(), nil
	case 254:
		return math.Inf(0), nil
	case 255:
		return math.Inf(-1), nil
	default:
		floatBytes := make([]byte, length)
		_, err := io.ReadFull(d.r, floatBytes)
		if err != nil {
			return 0, err
		}
		f, err := strconv.ParseFloat(string(floatBytes), 64)
		return f, err
	}

	panic("not reached")
}
