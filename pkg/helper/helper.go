package helper

import (
	"context"
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

func CreateRequestContext(userId string, age int) context.Context{
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", userId)
	ctx = context.WithValue(ctx, "age", age)
	return ctx
}