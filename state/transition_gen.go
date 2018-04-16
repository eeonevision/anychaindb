package state

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Transition) DecodeMsg(dc *msgp.Reader) (err error) {
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
		case "advertiser_account_id":
			z.AdvertiserAccountID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "affiliate_account_id":
			z.AffiliateAccountID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "click_id":
			z.ClickID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "stream_id":
			z.StreamID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "offer_id":
			z.OfferID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "created_at":
			z.CreatedAt, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "expires_in":
			z.ExpiresIn, err = dc.ReadInt64()
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
func (z *Transition) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 8
	// write "_id"
	err = en.Append(0x88, 0xa3, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		return
	}
	// write "advertiser_account_id"
	err = en.Append(0xb5, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.AdvertiserAccountID)
	if err != nil {
		return
	}
	// write "affiliate_account_id"
	err = en.Append(0xb4, 0x61, 0x66, 0x66, 0x69, 0x6c, 0x69, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.AffiliateAccountID)
	if err != nil {
		return
	}
	// write "click_id"
	err = en.Append(0xa8, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ClickID)
	if err != nil {
		return
	}
	// write "stream_id"
	err = en.Append(0xa9, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.StreamID)
	if err != nil {
		return
	}
	// write "offer_id"
	err = en.Append(0xa8, 0x6f, 0x66, 0x66, 0x65, 0x72, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.OfferID)
	if err != nil {
		return
	}
	// write "created_at"
	err = en.Append(0xaa, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.CreatedAt)
	if err != nil {
		return
	}
	// write "expires_in"
	err = en.Append(0xaa, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.ExpiresIn)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Transition) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 8
	// string "_id"
	o = append(o, 0x88, 0xa3, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ID)
	// string "advertiser_account_id"
	o = append(o, 0xb5, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.AdvertiserAccountID)
	// string "affiliate_account_id"
	o = append(o, 0xb4, 0x61, 0x66, 0x66, 0x69, 0x6c, 0x69, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.AffiliateAccountID)
	// string "click_id"
	o = append(o, 0xa8, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.ClickID)
	// string "stream_id"
	o = append(o, 0xa9, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.StreamID)
	// string "offer_id"
	o = append(o, 0xa8, 0x6f, 0x66, 0x66, 0x65, 0x72, 0x5f, 0x69, 0x64)
	o = msgp.AppendString(o, z.OfferID)
	// string "created_at"
	o = append(o, 0xaa, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74)
	o = msgp.AppendFloat64(o, z.CreatedAt)
	// string "expires_in"
	o = append(o, 0xaa, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e)
	o = msgp.AppendInt64(o, z.ExpiresIn)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Transition) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "advertiser_account_id":
			z.AdvertiserAccountID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "affiliate_account_id":
			z.AffiliateAccountID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "click_id":
			z.ClickID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "stream_id":
			z.StreamID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "offer_id":
			z.OfferID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "created_at":
			z.CreatedAt, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "expires_in":
			z.ExpiresIn, bts, err = msgp.ReadInt64Bytes(bts)
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
func (z *Transition) Msgsize() (s int) {
	s = 1 + 4 + msgp.StringPrefixSize + len(z.ID) + 22 + msgp.StringPrefixSize + len(z.AdvertiserAccountID) + 21 + msgp.StringPrefixSize + len(z.AffiliateAccountID) + 9 + msgp.StringPrefixSize + len(z.ClickID) + 10 + msgp.StringPrefixSize + len(z.StreamID) + 9 + msgp.StringPrefixSize + len(z.OfferID) + 11 + msgp.Float64Size + 11 + msgp.Int64Size
	return
}
