package gowebtoolkit

import "crypto/rand"

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

type Tools struct {}

// GenerateRandomString returns a random string of the given size using the characters of randomStringSource
func (t *Tools) GenerateRandomString(size int) string {
	s, r := make([]rune, size), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x % y]
	}

	return string(s)
}
