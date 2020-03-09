package privacy

import "testing"

func TestGCM(t *testing.T) {

	hello := NewMethodWithName("RAW")
	fuck, err := hello.Compress([]byte("zhangweihua"), MakeCompressKey("kdkzhangweihua"))
	t.Log(fuck, err)
	one, err := hello.Uncompress(fuck, MakeCompressKey("kdkzhangweihua"))
	t.Log(string(one), err)

}
