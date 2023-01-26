package pkg

import (
	"fmt"
)

// This function takes in a int representing bytes and returns a string of the input in the appropriate unit
func ByteToUnit(byteCount float64) string {
	units := []string{"B", "KiB", "MB", "GB", "TB"}
	unit := 0
	for byteCount > 1024 && unit < 4 {
		byteCount /= 1024
		unit++
	}
	return fmt.Sprintf("%.2f", byteCount) + units[unit]
}

func MegOrGig(bytes float64) string {
	if bytes >= (1 << 30) {
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1<<30))
	} else {
		fmt.Println(float64(bytes)/(1<<20))
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1<<20))
	}
}
