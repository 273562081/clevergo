// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package stringutil

import (
	"crypto/rand"
	"io"
)

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// Return part of a string.
func SubString(s string, start, length int) string {
	rs := []rune(s)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// Make a string's first character is uppercase.
func UpperFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	c := s[0]
	if ('a' <= c) && (c <= 'z') {
		return string(rune(int(c)-32)) + SubString(s, 1, len(s)-1)
	}
	return s
}

// Make a string's first character is lowercase.
func LowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	c := s[0]
	if ('A' <= c) && (c <= 'Z') {
		return string(rune(int(c)+32)) + SubString(s, 1, len(s)-1)
	}
	return s
}

// Generate random string.
// @length length of random string.
func GenerateRandomString(length int) string {
	return string(GenerateRandomByte(length))
}

// Generate random []byte.
// @length length of []byte.
func GenerateRandomByte(length int) []byte {
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil
	}
	return b
}

func RandomChars(length int, args ...[]byte) string {
	if length <= 0 {
		return ""
	}

	var chars []byte
	if len(args) > 0 {
		chars = args[0]
	} else {
		chars = StdChars
	}

	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for NewLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
