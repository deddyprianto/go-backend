package helper

import (
	"fmt"
	"time"
)

func Converter(bytes []uint8) (time.Time, error){
	str := string(bytes)
	t ,err := time.ParseInLocation("2006-01-02 15:04:05",str,time.Local)
	if err != nil{
		return time.Time{}, fmt.Errorf("error parsing time: %v", err)
	}
	return t, nil
}