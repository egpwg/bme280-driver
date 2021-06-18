package common

import (
	"log"
	"os"

	"golang.org/x/sys/unix"
)

// Ioctl read/write and other opertion from kernel
func ioctl(fd, op, data uintptr) (err unix.Errno) {
	_, _, err = unix.Syscall(unix.SYS_IOCTL, fd, op, data)
	return err
}

type Ioctler interface {
	Ioctl(op, data uintptr) (err unix.Errno)
}

type File struct {
	*os.File
}

func (f *File) Open(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}

	f.File = file

	return nil
}

func (f *File) Close() (err error) {
	return f.File.Close()
}

func (f *File) Ioctl(op, data uintptr) (err unix.Errno) {
	return ioctl(f.Fd(), op, data)
}

type FileOpen func(path string) (file *File, err error)
