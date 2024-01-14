package golibs

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ConvertDatetimeFromString(strDatetime string) (time.Time, error) {
	retTime, err := time.Parse("2006/01/02 15:04:05", strDatetime)
	if err != nil {
		return retTime, err
	}
	return retTime, nil
}

func ToIntByRemoveString(str string) int {
	n := 0
	for _, r := range str {
		if ('0' <= r && r <= '9') || (r == '-') {
			n = n*10 + int(r-'0')
		}
	}
	return n
}

func ToFloatByRemoveString(str string) (float64, error) {
	strNumber := ""
	slice := strings.Split(str, "")
	for _, s := range slice {
		if s == "." || s == "-" {
			strNumber += s
			continue
		}

		_, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		strNumber += s
	}
	f, err := strconv.ParseFloat(strNumber, 0)
	if err != nil {
		return 0.0, err
	}
	f, err = strconv.ParseFloat(fmt.Sprintf("%.2f", f), 0)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

func ToStringByFloat64(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}

func ToStringByInt(number int) string {
	return strconv.Itoa(number)
}
