package utils

import (
	"strconv"
	"strings"
)

type SplittedAddress struct {
	Host string
	Port uint16
}

func SplitAddress(address string) *SplittedAddress {
	parts := strings.Split(address, ":")
	if len(parts) != 2 {
		return nil
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	return &SplittedAddress{
		Host: host,
		Port: uint16(port),
	}
}
