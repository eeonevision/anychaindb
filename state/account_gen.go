package state

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Account) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "_id":
			z.ID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "pubkey":
			z.PubKey, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Account) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "_id"
	err = en.Append(0x82, 0xa3, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		return
	}
	// write "pubkey"
	err = en.Append(0xa6, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.PubKey)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Account) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "_id"
	o = append(o, 0x82, 0xa3, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "pubkey"
	o = append(o, 0xa6, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79)
	o = msgp.AppendString(o, z.PubKey)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Account) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "_id":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "pubkey":
			z.PubKey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Account) Msgsize() (s int) {
	s = 1 + 4 + msgp.StringPrefixSize + len(z.ID) + 7 + msgp.StringPrefixSize + len(z.PubKey)
	return
}
