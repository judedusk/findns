package main

import (
	"os"
	"syscall"
	"unsafe"
)

const enableVirtualTerminalProcessing = 0x0004

func init() {
	enableANSI(os.Stderr)
	enableANSI(os.Stdout)
}

func enableANSI(f *os.File) {
	var mode uint32
	h := syscall.Handle(f.Fd())
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getMode := kernel32.NewProc("GetConsoleMode")
	setMode := kernel32.NewProc("SetConsoleMode")

	r, _, _ := getMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return
	}
	setMode.Call(uintptr(h), uintptr(mode|enableVirtualTerminalProcessing))
}
