// +build windows

// Platform
package main

import (
    "fmt"
    "log"
    "strings"
    "syscall"
    "unsafe"
)

const (
    VER_PLATFORM_WIN32_WINDOWS = 1
    VER_PLATFORM_WIN32_NT      = 2
)

const (
    FALSE = 0
    TRUE  = 1
)

type OSVERSIONINFO struct {
    dwOSVersionInfoSize uint32
    dwMajorVersion      uint32
    dwMinorVersion      uint32
    dwBuildNumber       uint32
    dwPlatformId        uint32
    szCSDVersion        [128]uint16
}

func getOSVersion() (maj, min, buildno, plat int, csd string) {
    var os OSVERSIONINFO
    os.dwOSVersionInfoSize = uint32(unsafe.Sizeof(os))

    dll := syscall.MustLoadDLL("kernel32.dll")
    p := dll.MustFindProc("GetVersionExW")

    v, _, err := p.Call(uintptr(unsafe.Pointer(&os)))
    if v == FALSE {
        log.Fatal(err)
    }

    maj = int(os.dwMajorVersion)
    min = int(os.dwMinorVersion)
    buildno = int(os.dwBuildNumber)
    plat = int(os.dwPlatformId)
    csd = syscall.UTF16ToString(os.szCSDVersion[:])
    return
}

func regQueryStringValue(root syscall.Handle, path, key string) (value string) {
    var err error
    var h syscall.Handle
    pathp, _ := syscall.UTF16PtrFromString(path)
    err = syscall.RegOpenKeyEx(root, pathp, 0, syscall.KEY_READ, &h)
    if err != nil {
        log.Fatal(err)
    }
    defer syscall.RegCloseKey(h)

    var valtype uint32
    var buf [1 << 10]uint16
    buflen := uint32(len(buf) * 2)
    keyp, _ := syscall.UTF16PtrFromString("ProductName")
    err = syscall.RegQueryValueEx(h, keyp, nil, &valtype, (*byte)(unsafe.Pointer(&buf[0])), &buflen)
    if err != nil {
        log.Fatal(err)
    }
    if valtype != syscall.REG_SZ {
        return
    }
    value = syscall.UTF16ToString(buf[:])
    return
}

func Platform() (platform string) {
    maj, min, buildno, plat, csd := getOSVersion()
    switch plat {
    case VER_PLATFORM_WIN32_NT:
        name := regQueryStringValue(syscall.HKEY_LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", "ProductName")
        platform = fmt.Sprintf("%s-%d.%d.%d", strings.Replace(name, " ", "-", -1), maj, min, buildno)
        if csd != "" {
            platform = platform + "-" + strings.Replace(csd, " ", "-", -1)
        }
    }
    return
}

func main() {
    fmt.Println(Platform())
}
