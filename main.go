package main

import (
	common "github.com/matehaxor03/holistic_common/common"
)

func main() {
	common.EscapeString("a'", true, "'")
}
