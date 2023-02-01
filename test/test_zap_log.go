package main

import (
	"github.com/xulehexuwei/scikits"
	"go.uber.org/zap"
)

func division(x, y int) (result int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	result = x / y
	return result, nil
}

func main() {
	defer scikits.Logger.Sync()
	for i := 0; i < 40; i++ {
		scikits.Logger.Info("info...", zap.Int("message", i), zap.String("url", "url"))
		scikits.SugarLogger.Errorf("this is error message %d", i)
	}

	_, err := division(1, 0)

	scikits.SugarLogger.Error(err)
}
