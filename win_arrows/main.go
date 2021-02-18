package main

import "os"

func main() {
	err := SetRaw()
	if err != nil {
		panic(err)
	}

	buf := make([]uint16, 20)
	Read(buf)
	os.Stdout.Write([]byte("Read"))
	Read(buf)
	os.Stdout.Write([]byte("Read"))
	Read(buf)
	os.Stdout.Write([]byte("Read"))
	Read(buf)
	os.Stdout.Write([]byte("Read"))
	Read(buf)
	os.Stdout.Write([]byte("Read"))
	Read(buf)
	os.Stdout.Write([]byte("Read"))

	Close()
}
