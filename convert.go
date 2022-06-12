package golibs

import "time"

func ConvertDatetimeFromString(strDatetime string) (time.Time, error) {
	retTime, err := time.Parse("2006/01/02 15:04:05", strDatetime)
	if err != nil {
		return retTime, err
	}
	return retTime, nil
}