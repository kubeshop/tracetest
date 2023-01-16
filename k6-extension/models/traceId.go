package models

import (
	cr "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	K6Prefix     = 0756 // Being 075 the ASCII code for 'K' :)
	K6Code_Cloud = 12   // To ingest and process the related spans in k6 Cloud.
	K6Code_Local = 33   // To not ingest and process the related spans, b/c they are part of a non-cloud run.
)

type TraceID struct {
	Prefix            int16
	Code              int8
	UnixTimestampNano uint64
}

func (t *TraceID) IsValid() bool {
	return t.Prefix == K6Prefix && (t.Code == K6Code_Cloud || t.Code == K6Code_Local)
}

func (t *TraceID) IsValidCloud() bool {
	return t.Prefix == K6Prefix && t.Code == K6Code_Cloud
}

func (t *TraceID) Encode() (string, []byte, error) {
	if !t.IsValid() {
		return "", nil, fmt.Errorf("failed to encode traceID: %v", t)
	}

	buf := make([]byte, 16)

	n := binary.PutVarint(buf, int64(t.Prefix))
	n += binary.PutVarint(buf[n:], int64(t.Code))
	n += binary.PutUvarint(buf[n:], t.UnixTimestampNano)

	randomness := make([]byte, 16-len(buf[:n]))
	err := binary.Read(cr.Reader, binary.BigEndian, randomness)
	if err != nil {
		return "", nil, err
	}

	buf = append(buf[:n], randomness[:]...)
	hx := hex.EncodeToString(buf)
	return hx, buf, nil
}

func DecodeTraceID(buf []byte) *TraceID {
	pre, preLen := binary.Varint(buf)
	code, codeLen := binary.Varint(buf[preLen:])
	ts, _ := binary.Uvarint(buf[preLen+codeLen:])

	return &TraceID{
		Prefix:            int16(pre),
		Code:              int8(code),
		UnixTimestampNano: uint64(ts),
	}
}
