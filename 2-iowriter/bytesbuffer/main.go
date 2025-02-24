package main

import (
	"bytes"
	"fmt"
)

func main() {
	var buffer bytes.Buffer
	//buffer.Write([]byte("bytes.Buffer example\n"))
	buffer.WriteString("bytes.Buffer example\n")
	fmt.Println(buffer.String())
}
