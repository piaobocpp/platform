// +build linux

// Platform
package main

import (
    "fmt"
    "log"
    "strings"
    "syscall"
)

func charsToString(ca []int8) string {
    s := make([]byte, len(ca))
    var lens int
    for ; lens < len(ca); lens++ {
        if ca[lens] == 0 {
            break
        }
        s[lens] = uint8(ca[lens])
    }
    return string(s[0:lens])
}

func Platform() (platform string) {
    var buf syscall.Utsname
    err := syscall.Uname(&buf)
    if err != nil {
        log.Fatal(err)
    }
    s := []string{charsToString(buf.Sysname[:]), charsToString(buf.Release[:]), charsToString(buf.Version[:])}
    platform = strings.Join(s, "-")
    return
}

func main() {
    fmt.Println(Platform())
}
