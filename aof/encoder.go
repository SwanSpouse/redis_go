package aof

import (
	"io"
	"os"
	"redis_go/database"
	"redis_go/encodings"
	"redis_go/loggers"
	"redis_go/util"
)

type Encoder struct {
	f *os.File
}

func NewEncoder(filename string) (*Encoder, error) {
	var f *os.File
	var err error
	if util.FileExists(filename) {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0600)
		f.Seek(0, io.SeekEnd)
	} else {
		f, err = os.Create(filename)
	}
	if err != nil {
		return nil, err
	}
	return &Encoder{f: f}, nil
}

func (encoder *Encoder) Write(buf []byte) (int, error) {
	return encoder.f.Write(buf)
}

func (encoder *Encoder) rewriteStringObject(obj database.TBase) {
	if obj.GetObjectType() != encodings.RedisTypeString {
		loggers.Errorf("obj type is not string %s", obj.GetObjectType())
	}
	switch obj.GetEncoding() {
	case encodings.RedisEncodingInt:
	case encodings.RedisEncodingEmbStr:
	case encodings.RedisEncodingRaw:
	}
}

func (encoder *Encoder) rewriteListObject(obj database.TBase) {
	if obj.GetObjectType() != encodings.RedisTypeList {
		loggers.Errorf("obj type:%s is not list ", obj.GetObjectType())
	}
}

func (encoder *Encoder) rewriteHashObject(obj database.TBase) {
	if obj.GetObjectType() != encodings.RedisTypeHash {
		loggers.Errorf("obj type:%s is not hash ", obj.GetObjectType())
	}
}

func (encoder *Encoder) rewriteSetObject(obj database.TBase) {
	if obj.GetObjectType() != encodings.RedisTypeSet {
		loggers.Errorf("obj type:%s is not set ", obj.GetObjectType())
	}
}

func (encoder *Encoder) rewriteZSetObject(obj database.TBase) {
	if obj.GetObjectType() != encodings.RedisTypeZSet {
		loggers.Errorf("obj type:%s is not zset ", obj.GetObjectType())
	}
}
