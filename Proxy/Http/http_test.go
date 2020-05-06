package Http

import (
	"fmt"
	"testing"
)

func Test_http(t *testing.T) {

}
func Test_byte(t *testing.T) {
	//_ := []byte("123456789")
	// byte类型存的ascii码
	asciiQQ := []byte{0x38, 0x36, 0x31, 0x37, 0x32, 0x39, 0x30, 0x39, 0x31}
	fmt.Println(string(asciiQQ))
}
