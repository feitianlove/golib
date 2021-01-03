package icmp_tools

import (
	"fmt"
	"testing"
)

func TestDoPing(t *testing.T) {
	fmt.Println(DoPing("127.0.0.1"))
	fmt.Println(DoPing("10.49.134.224"))
	fmt.Println(DoPing("9.49.134.224"))
}
