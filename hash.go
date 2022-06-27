package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash/crc32"
	"strings"

	t "github.com/rivo/tview"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

var hashScreen = struct {
	view   *t.Flex
	hashes []*hashInfo
}{
	hashes: []*hashInfo{
		{
			// TODO: Check for more variants at https://crccalc.com/
			name: "CRC32",
			calc: func(b []byte) []byte {
				c := crc32.ChecksumIEEE(b)
				return []byte{
					byte((c >> 24) & 255),
					byte((c >> 16) & 255),
					byte((c >> 8) & 255),
					byte(c & 255),
				}
			},
		},
		{
			name: "MD5",
			calc: func(b []byte) []byte {
				return hashWith(b, md5.Sum)
			},
		},
		{
			name: "SHA1",
			calc: func(b []byte) []byte {
				return hashWith(b, sha1.Sum)
			},
		},
		{
			name: "SHA224",
			calc: func(b []byte) []byte {
				return hashWith(b, sha256.Sum224)
			},
		},
		{
			name: "SHA256",
			calc: func(b []byte) []byte {
				return hashWith(b, sha256.Sum256)
			},
		},
		{
			name: "SHA384",
			calc: func(b []byte) []byte {
				return hashWith(b, sha512.Sum384)
			},
		},
		{
			name: "SHA512",
			calc: func(b []byte) []byte {
				return hashWith(b, sha512.Sum512)
			},
		},
		{
			name: "SHA512/224",
			calc: func(b []byte) []byte {
				return hashWith(b, sha512.Sum512_224)
			},
		},
		{
			name: "SHA512/256",
			calc: func(b []byte) []byte {
				return hashWith(b, sha512.Sum512_256)
			},
		},
		{
			name: "SHA3-224",
			calc: func(b []byte) []byte {
				return hashWith(b, sha3.Sum224)
			},
		},
		{
			name: "SHA3-256",
			calc: func(b []byte) []byte {
				return hashWith(b, sha3.Sum256)
			},
		},
		{
			name: "SHA3-384",
			calc: func(b []byte) []byte {
				return hashWith(b, sha3.Sum384)
			},
		},
		{
			name: "SHA3-512",
			calc: func(b []byte) []byte {
				return hashWith(b, sha3.Sum512)
			},
		},
		{
			name: "SHAKE128-256",
			calc: func(b []byte) []byte {
				h := make([]byte, 256/8)
				sha3.ShakeSum128(h, b)
				return h
			},
		},
		{
			name: "SHAKE256-512",
			calc: func(b []byte) []byte {
				h := make([]byte, 512/8)
				sha3.ShakeSum256(h, b)
				return h
			},
		},
		{
			name: "BLAKE2b-256",
			calc: func(b []byte) []byte {
				return hashWith(b, blake2b.Sum256)
			},
		},
		{
			name: "BLAKE2b-384",
			calc: func(b []byte) []byte {
				return hashWith(b, blake2b.Sum384)
			},
		},
		{
			name: "BLAKE2b-512",
			calc: func(b []byte) []byte {
				return hashWith(b, blake2b.Sum512)
			},
		},
	},
}

func hashWith[R []byte | [16]byte | [20]byte | [28]byte | [32]byte | [48]byte | [64]byte](b []byte, hash func([]byte) R) []byte {
	h := hash(b)
	c := make([]byte, len(h))
	for i := 0; i < len(h); i++ {
		c[i] = h[i]
	}
	return c
}

type hashInfo struct {
	name string
	calc func([]byte) []byte
	view *t.TextView
}

func makeHashScreen() t.Primitive {
	hashScreen.view = NewFlexColumn().
		AddItem(t.NewInputField().SetLabel("Input: ").SetChangedFunc(updateHashes), 1, 0, true).
		AddItem(t.NewBox(), 1, 0, false)

	maxNameLength := 0
	for _, h := range hashScreen.hashes {
		if len(h.name) > maxNameLength {
			maxNameLength = len(h.name)
		}
	}

	for _, h := range hashScreen.hashes {
		h.view = t.NewTextView()
		name := h.name
		if len(name) < maxNameLength {
			name = strings.Repeat(" ", maxNameLength-len(name)) + name
		}
		hashScreen.view.
			AddItem(wrapWithLabel(h.view, name+": "), 1, 0, false)
	}

	hashScreen.view.AddItem(t.NewBox(), 0, 1, false)

	updateHashes("")

	return hashScreen.view
}

func updateHashes(input string) {
	for _, h := range hashScreen.hashes {
		h.view.SetText(hex.EncodeToString(h.calc([]byte(input))))
	}
}
