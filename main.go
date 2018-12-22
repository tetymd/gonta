package main

import (
	"fmt"
    "os/exec"
    "os"
    "io"
    "syscall"
    "github.com/kr/pty"
)

func main() {
    pid := os.Getpid()
    fmt.Println(pid)

    err := syscall.Unshare(syscall.CLONE_NEWPID)
    if err != nil {
        fmt.Println(err)
        return
    } else {
        fmt.Println("No Error")
    }

    ptmx, pts, err := pty.Open()
    if err != nil {
        fmt.Println(err)
    }
    defer func() { _ = ptmx.Close()  }() // Best effort.
    defer func() { _ = pts.Close()  }() // Best effort.


    c := exec.Command("/bin/bash")

    c.Stdin = pts
    c.Stdout = pts
    c.Stderr = pts
    c.Dir = "/"

    if c.SysProcAttr == nil {
        c.SysProcAttr = &syscall.SysProcAttr{}
    }
    c.SysProcAttr.Setctty = true
    c.SysProcAttr.Setsid = true
    //c.SysProcAttr.Cloneflags = syscall.CLONE_NEWNS
    //c.SysProcAttr.Unshareflags = syscall.CLONE_NEWPID

    err = syscall.Chroot("./container")
    if err != nil{
        fmt.Println(err)
    }


    fmt.Println("run bash")
    err = c.Start()
    if err != nil {
        ptmx.Close()
        fmt.Println(err)
    }

    go func() { _, _ = io.Copy(ptmx, os.Stdin)  }()
    _, _ = io.Copy(os.Stdout, ptmx)

    pid = os.Getpid()
    fmt.Println(pid)
}
