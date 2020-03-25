package privacy

import (
	"strings"
	"testing"
)

func tospace(s string) string {

	return strings.ToUpper(s)
}

func TestGCM(t *testing.T) {

	/*u := binary.LittleEndian.Uint32([]byte{1, 0, 0, 0})
	t.Log(u)

	hello := NewMethodWithName("AES-GCM")
	hello.SetIV([]byte("123456789012"))

	key := []byte("12345678901234567890123456789012")
	t.Log("key = ", tospace(hex.EncodeToString(key)))
	ori := []byte("Im a secret message!")
	t.Log("ori = ", tospace(hex.EncodeToString(ori)))
	fuck, _ := hello.Compress(ori, key)
	t.Log("fuck = ", tospace(hex.EncodeToString(fuck)))

	one, _ := hello.Uncompress(fuck, key)

	t.Log("one = ", tospace(hex.EncodeToString(one)))*/

	//key := []byte("12345678901234567890123456789012")
	//ori := []byte("Im a secret message!")
	//kk := MakeCompressKey(string(ori))
	hello := NewMethodWithName("CHACHA-Ploy1305")
	hello.SetIV([]byte("123456789012123456789012"))
	t.Log(hello.ToBytes())

}

//12340C313233343536373839303132
//12340C313233343536373839303132
//C99FDD8F2232A48297774BA4BC701B2F7FD973D534DA9B9B08BA9C86F5F7167C21AAECF5
//C99FDD8F2232A48297774BA4BC701B2F7FD973D534DA9B9B08BA9C86F5F7167C21AAECF5
