package uuid

import (
	"crypto/rand"
	"fmt"
	"log"
)

// A UUID representation compliant with specification in
// RFC 4122 document.
type UUID [16]byte

// V4 generates a v4 UUID.
func V4() (*UUID, error) {
	u := new(UUID)
	_, err := rand.Read(u[:])
	if err != nil {
		return nil, err
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (4 << 4)
	return u, nil
}

// MustV4 generates a v4 UUID and return uuid string
func MustV4() string {
	u, err := V4()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

// Returns unparsed version of the generated UUID sequence.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
