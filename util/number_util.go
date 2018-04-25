package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/**
去除末尾0后，精度不够2位的，保留两位；
去除末尾0后，精度超过2位的，去除末尾的0
*/
func FormatFloatString(input string) (string, error) {
	if fVal, err := strconv.ParseFloat(input, 64); err != nil {
		return "", err
	} else {
		if !strings.Contains(input, ".") {
			return input, nil
		}
		fields := strings.Split(input, ".")
		if len(fields) != 2 {
			return "", errors.New("input number is invalid")
		}
		if len(strings.TrimRight(fields[1], "0")) <= 2 {
			return fmt.Sprintf("%.2f", fVal), nil
		}
		return strings.TrimRight(input, "0"), nil
	}
}
