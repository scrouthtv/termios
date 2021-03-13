package main

import "fmt"
import "reflect"
import "syscall"

import "golang.org/x/sys/windows"

func main() {
	fmt.Println("starting analysis...")

	in, err := windows.Open("CONOUT$", windows.O_RDWR, 0)
	if err != nil {
		fmt.Println("error opening the console:", err)
		return
	}
	defer windows.Close(in)

	ok, err := testVTSupport(in)
	if err != nil {
		fmt.Println("error determining support:", err)
		return
	}
	if ok {
		fmt.Println("your terminal does support vt emulation")
	} else {
		fmt.Println("your terminal does not support vt emulation")
	}
}

func testVTSupport(con windows.Handle) (bool, error) {

	var oldMode uint32
	err := windows.GetConsoleMode(con, &oldMode)
	if err != nil {
		fmt.Println("error getting console mode:", err)
		return false, err
	}

	mode := oldMode
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = windows.SetConsoleMode(con, mode)
	fmt.Println(reflect.TypeOf(err))
	defer windows.SetConsoleMode(con, oldMode)

	serr, ok := err.(syscall.Errno)
	if ok {
		ierr := uint64(uintptr(serr))
		if ierr == 87 {
			fmt.Println("error code is 87")
			return false, nil
		} else {
			fmt.Println("invalid error code:", uintptr(serr))
			return false, err
		}
	} else {
		fmt.Println("is not a syscall.Errno")
	}

	return true, nil
}
