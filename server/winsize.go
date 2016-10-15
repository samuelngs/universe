package server

import (
	"syscall"
	"unsafe"
)

// Winsize
type Winsize struct {
	Height, Width, x, y uint16
}

// SetWinsize sets the size of the given pty.
func setWinsize(fd uintptr, w, h uint32) {
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}
