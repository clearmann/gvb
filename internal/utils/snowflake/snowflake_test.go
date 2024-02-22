package sonwflake

import (
	"fmt"
	"testing"
)

func TestGenID(t *testing.T) {
	Init("2003-01-23", 2)
	fmt.Println(GenID())
}
