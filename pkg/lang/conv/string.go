package conv

import "unsafe"

// StringToBytes Unsafe string to []byte
func StringToBytes(val string) []byte {
	sh := (*[2]uintptr)(unsafe.Pointer(&val))
	bh := [3]uintptr{sh[0], sh[1], sh[1]}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

// BytesToString Unsafe []byte to string
func BytesToString(val []byte) string {
	bh := (*[3]uintptr)(unsafe.Pointer(&val))
	sh := [2]uintptr{bh[0], bh[1]}

	return *(*string)(unsafe.Pointer(&sh))
}

func BytesToStr(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	if len(b) < 64 {
		return string(b)
	}

	for _, v := range b {
		if v > 127 {
			return string(b)
		}
	}

	return BytesToString(b)
}
