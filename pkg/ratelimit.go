package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GetRateLimit returns the rate limit in int64 and checks if the rate limit is valid
func GetRateLimit(rateLimit string) (int64, error) {
	if rateLimit == "0" || rateLimit == "" {
		return 0, nil
	}
	unit := strings.ToLower(rateLimit[len(rateLimit)-1:])
	value, err := strconv.Atoi(rateLimit[:len(rateLimit)-1])
	if err != nil {
		fmt.Println("invalid rate limit, use k or M")
		return 0, err
	}
	if unit != "k" && unit != "m" {
		return 0, errors.New("invalid unit of measurement, use k or M")
	}
	if unit == "k" {
		return int64(value * 1000), nil
	}
	return int64(value * 1000000), nil
}
