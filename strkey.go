package main

import (
	"bytes"
	"fmt"
	"encoding/base32"
	"github.com/xdrpp/goxdr/xdr"
)

type StrKeyError string
func (e StrKeyError) Error() string { return string(e) }

type StrKeyVersionByte byte

var b32	= base32.StdEncoding.WithPadding(base32.NoPadding)
// var b32	= base32.StdEncoding

const (
	STRKEY_ALG_ED25519 = 0
)

const (
	STRKEY_PUBKEY         StrKeyVersionByte = 6<<3  // 'G'
	STRKEY_MUXED          StrKeyVersionByte = 12<<3 // 'M'
	STRKEY_PRIVKEY        StrKeyVersionByte = 18<<3 // 'S'
	STRKEY_PRE_AUTH_TX    StrKeyVersionByte = 19<<3 // 'T',
	STRKEY_HASH_X         StrKeyVersionByte = 23<<3 // 'X'
	STRKEY_ERROR          StrKeyVersionByte = 255
)

var payloadLen = map[StrKeyVersionByte]int {
	STRKEY_PUBKEY|STRKEY_ALG_ED25519: 32,
	STRKEY_MUXED|STRKEY_ALG_ED25519: 40,
	STRKEY_PRIVKEY|STRKEY_ALG_ED25519: 32,
	STRKEY_PRE_AUTH_TX: 32,
	STRKEY_HASH_X: 32,
}

var crc16table [256]uint16

func init() {
	const poly = 0x1021
	for i := 0; i < 256; i++ {
		crc := uint16(i) << 8
		for j := 0; j < 8; j++ {
			if crc&0x8000 != 0 {
				crc = crc<<1 ^ poly
			} else {
				crc <<= 1
			}
		}
		crc16table[i] = crc
	}
}

func crc16(data []byte) (crc uint16) {
	for _, b := range data {
		temp := b ^ byte(crc>>8)
		crc = crc16table[temp] ^ (crc << 8)
	}
	return
}

// ToStrKey converts the raw bytes of a key to ASCII strkey format.
func ToStrKey(ver StrKeyVersionByte, bin []byte) string {
	var out bytes.Buffer
	out.WriteByte(byte(ver))
	out.Write(bin)
	sum := crc16(out.Bytes())
	out.WriteByte(byte(sum))
	out.WriteByte(byte(sum >> 8))
	return b32.EncodeToString(out.Bytes())
}

// FromStrKey decodes a strkey-format string into the raw bytes of the
// key and the type of key.  Returns the reserved StrKeyVersionByte
// STRKEY_ERROR if it fails to decode the string.
func FromStrKey(in []byte) ([]byte, StrKeyVersionByte) {
	if rem := len(in) % 8; rem == 1 || rem == 3 || rem == 6 {
		return nil, STRKEY_ERROR
	}
	bin := make([]byte, b32.DecodedLen(len(in)))
	n, err := b32.Decode(bin, in)
	if err != nil || n != len(bin) || n < 3 {
		return nil, STRKEY_ERROR
	}
	if targetlen, ok := payloadLen[StrKeyVersionByte(bin[0])]; !ok ||
		targetlen != n - 3 {
		return nil, STRKEY_ERROR
	}
	want := uint16(bin[len(bin)-2]) | uint16(bin[len(bin)-1])<<8
	if want != crc16(bin[:len(bin)-2]) {
		return nil, STRKEY_ERROR
	}
	return bin[1 : len(bin)-2], StrKeyVersionByte(bin[0])
}

func XdrToBytes(t xdr.XdrType) []byte {
        out := bytes.Buffer{}
        t.XdrMarshal(&xdr.XdrOut{&out}, "")
        return out.Bytes()
}

func XdrFromBytes(t xdr.XdrType, input []byte) (err error) {
	defer func() {
		if i := recover(); i != nil {
			if xe, ok := i.(error); ok {
				err = xe
				return
			}
			panic(i)
		}
	}()
	in := bytes.NewReader(input)
	t.XdrMarshal(&xdr.XdrIn{in}, "")
	return
}

// Renders a PublicKey in strkey format.
func (pk PublicKey) String() string {
	switch pk.Type {
	case PUBLIC_KEY_TYPE_ED25519:
		return ToStrKey(STRKEY_PUBKEY|STRKEY_ALG_ED25519, pk.Ed25519()[:])
	default:
		return fmt.Sprintf("PublicKey.Type#%d", int32(pk.Type))
	}
}

// Renders a PublicKey in strkey format.
func (pk MuxedAccount) String() string {
	switch pk.Type {
	case KEY_TYPE_ED25519:
		return ToStrKey(STRKEY_PUBKEY|STRKEY_ALG_ED25519, pk.Ed25519()[:])
	case KEY_TYPE_MUXED_ED25519:
		return ToStrKey(STRKEY_MUXED|STRKEY_ALG_ED25519,
			XdrToBytes(pk.Med25519()))
	default:
		return fmt.Sprintf("PublicKey.Type#%d", int32(pk.Type))
	}
}

// Renders a SignerKey in strkey format.
func (pk SignerKey) String() string {
	switch pk.Type {
	case SIGNER_KEY_TYPE_ED25519:
		return ToStrKey(STRKEY_PUBKEY|STRKEY_ALG_ED25519, pk.Ed25519()[:])
	case SIGNER_KEY_TYPE_PRE_AUTH_TX:
		return ToStrKey(STRKEY_PRE_AUTH_TX, pk.PreAuthTx()[:])
	case SIGNER_KEY_TYPE_HASH_X:
		return ToStrKey(STRKEY_HASH_X, pk.HashX()[:])
	default:
		return fmt.Sprintf("SignerKey.Type#%d", int32(pk.Type))
	}
}

// Returns true if c is a valid character in a strkey formatted key.
func IsStrKeyChar(c rune) bool {
	return c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
}

// Parses a public key in strkey format.
func (pk *PublicKey) Scan(ss fmt.ScanState, _ rune) error {
	bs, err := ss.Token(true, IsStrKeyChar)
	if err != nil {
		return err
	}
	return pk.UnmarshalText(bs)
}

// Parses a public key in strkey format.
func (pk *MuxedAccount) Scan(ss fmt.ScanState, _ rune) error {
	bs, err := ss.Token(true, IsStrKeyChar)
	if err != nil {
		return err
	}
	return pk.UnmarshalText(bs)
}

// Parses a signer in strkey format.
func (pk *SignerKey) Scan(ss fmt.ScanState, _ rune) error {
	bs, err := ss.Token(true, IsStrKeyChar)
	if err != nil {
		return err
	}
	return pk.UnmarshalText(bs)
}

// Parses a public key in strkey format.
func (pk *PublicKey) UnmarshalText(bs []byte) error {
	key, vers := FromStrKey(bs)
	switch vers {
	case STRKEY_PUBKEY|STRKEY_ALG_ED25519:
		pk.Type = PUBLIC_KEY_TYPE_ED25519
		copy(pk.Ed25519()[:], key)
		return nil
	default:
		return StrKeyError("Invalid public key type")
	}
}

// Parses a public key in strkey format.
func (pk *MuxedAccount) UnmarshalText(bs []byte) error {
	key, vers := FromStrKey(bs)
	switch vers {
	case STRKEY_PUBKEY|STRKEY_ALG_ED25519:
		pk.Type = KEY_TYPE_ED25519
		copy(pk.Ed25519()[:], key)
		return nil
	case STRKEY_MUXED|STRKEY_ALG_ED25519:
		pk.Type = KEY_TYPE_MUXED_ED25519
		return XdrFromBytes(pk.Med25519(), key)
	default:
		return StrKeyError("Invalid public key type")
	}
}

// Parses a signer in strkey format.
func (pk *SignerKey) UnmarshalText(bs []byte) error {
	key, vers := FromStrKey(bs)
	switch vers {
	case STRKEY_PUBKEY|STRKEY_ALG_ED25519:
		pk.Type = SIGNER_KEY_TYPE_ED25519
		copy(pk.Ed25519()[:], key)
	case STRKEY_PRE_AUTH_TX:
		pk.Type = SIGNER_KEY_TYPE_PRE_AUTH_TX
		copy(pk.PreAuthTx()[:], key)
	case STRKEY_HASH_X:
		pk.Type = SIGNER_KEY_TYPE_HASH_X
		copy(pk.HashX()[:], key)
	default:
		return StrKeyError("Invalid signer key string")
	}
	return nil
}
