package golibs

import (
	"fmt"
	"time"
)

func WaitTimeMillSecond(millSecond int) {
	time.Sleep(time.Duration(millSecond) * time.Millisecond)
}

func WaitTimeSecond(second int) {
	time.Sleep(time.Duration(second) * time.Second)
}

// WaitForCondition 与えた引数の関数がtrueを返すまで待つ
func WaitForCondition(waitTimeMillSecond int, intervalMillSecond int, conditionFunc func(limitTimeMillSecond int) (bool, error)) error {
	for i := waitTimeMillSecond; i > 0; i = i - intervalMillSecond {
		WaitTimeMillSecond(intervalMillSecond)
		is, err := conditionFunc(i)
		if err != nil {
			return err
		}
		if is {
			return nil
		}
	}
	return fmt.Errorf("timeout")
}
