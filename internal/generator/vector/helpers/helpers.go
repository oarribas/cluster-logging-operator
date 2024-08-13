package helpers

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

var (
	Replacer         = strings.NewReplacer(" ", "_", "-", "_", ".", "_")
	listenAllAddress string
	listenAllOnce    sync.Once
)

func MakeInputs(in ...string) string {
	out := make([]string, len(in))
	for i, o := range in {
		if strings.HasPrefix(o, "\"") && strings.HasSuffix(o, "\"") {
			out[i] = o
		} else {
			out[i] = fmt.Sprintf("%q", o)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(out, ","))
}

func TrimSpaces(in []string) []string {
	o := make([]string, len(in))
	for i, s := range in {
		o[i] = strings.TrimSpace(s)
	}
	return o
}

func FormatComponentID(name string) string {
	return strings.ToLower(Replacer.Replace(name))
}

func ListenOnAllLocalInterfacesAddress() string {
	f := func() {
		listenAllAddress = func() string {
			if fd, err := unix.Socket(unix.AF_INET6, unix.SOCK_STREAM, unix.IPPROTO_IP); err != nil {
				return `0.0.0.0`
			} else {
				unix.Close(fd)
				return `[::]`
			}
		}()
	}
	listenAllOnce.Do(f)
	return listenAllAddress
}

func EscapeDollarSigns(s string) string {
	return strings.ReplaceAll(s, "$", "$$")
}
