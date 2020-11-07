package publicfunc

import (
	"fmt"
	"testing"
)

func Test_StringToDestinationName(t *testing.T) {
	str := "@fuck asdasdasdasd"
	res := StringToDestinationName(str)
	fmt.Println(res)
}

func Test_StringToDestinationAddr(t *testing.T) {
	str := "@fuck asdasdasdasd"
	res := StringToDestinationAddr(str)
	fmt.Println(res)
	fmt.Println(len(res))
}
