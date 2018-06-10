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

// 最高保留6位精度
func FloatToSimpleString(input float64) string {
	ret := strings.TrimRight(fmt.Sprintf("%.6f", input), "0")
	if len(ret) > 0 && ret[len(ret)-1] == '.' {
		ret += "0"
	}
	return ret
}
