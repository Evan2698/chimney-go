package privacy

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGCM(t *testing.T) {

	hello := NewMethodWithName("AES-GCM")
	hello.SetIV([]byte("123456789012"))
	fuck, err := hello.Compress([]byte("Im a secret message!"), []byte("12345678901234567890123456789012"))

	t.Log(strings.ToUpper(hex.EncodeToString(fuck)), err, len(fuck), strings.ToUpper(hex.EncodeToString(hello.ToBytes())))
	one, err := hello.Uncompress(fuck, []byte("12345678901234567890123456789012"))
	t.Log(string(one), err)

}

//12340C313233343536373839303132
//12340C313233343536373839303132
//C99FDD8F2232A48297774BA4BC701B2F7FD973D534DA9B9B08BA9C86F5F7167C21AAECF5
//C99FDD8F2232A48297774BA4BC701B2F7FD973D534DA9B9B08BA9C86F5F7167C21AAECF5
