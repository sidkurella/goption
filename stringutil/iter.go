package stringutil

import "github.com/sidkurella/goption/option"

type stringIter struct {
	s []rune
	i int
}

// Returns an iterator of runes in the string.
// For a UTF-8 encoded string, certain code points may be encoded with multiple bytes.
// This will return each code point, not each individual byte.
// If you wish to iterate by raw bytes, use ByteIter() instead.
func Iter(s string) *stringIter {
	return &stringIter{
		s: []rune(s),
		i: 0,
	}
}

func (s *stringIter) Next() option.Option[rune] {
	if s.i >= len(s.s) {
		return option.Nothing[rune]()
	}
	ret := option.Some(s.s[s.i])
	s.i++
	return ret
}

type stringByteIter struct {
	s string
	i int
}

// Returns an iterator of bytes in the string.
// For a UTF-8 encoded string, certain code points may be encoded with multiple bytes.
// This will return each byte individually.
// If you wish to iterate by code points, use Iter() instead, which returns runes.
func ByteIter(s string) *stringByteIter {
	return &stringByteIter{
		s: s,
		i: 0,
	}
}

func (s *stringByteIter) Next() option.Option[byte] {
	if s.i >= len(s.s) {
		return option.Nothing[byte]()
	}
	ret := option.Some(s.s[s.i])
	s.i++
	return ret
}
