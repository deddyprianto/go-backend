package helper

import (
	"context"
	"fmt"
	"os"
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

func ParseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",           // MySQL DATETIME
		"2006-01-02T15:04:05Z",          // ISO 8601 UTC
		"2006-01-02T15:04:05.000Z",      // ISO 8601 dengan milliseconds
		"2006-01-02T15:04:05-07:00",     // ISO 8601 dengan timezone
		"2006-01-02T15:04:05.000-07:00", // ISO 8601 dengan milliseconds dan timezone
		time.RFC3339,                    // RFC3339
		time.RFC3339Nano,
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateStr)
}

// Function to save profile picture and return the file path
func SaveProfilePicture(profilePicture string) (string, error) {
    // Create a directory for profile pictures if it doesn't exist
    dir := "profile_pictures"
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.Mkdir(dir, 0755)
    }

    // Generate a unique filename
    fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), "profile.jpg")
    filePath := fmt.Sprintf("%s/%s", dir, fileName)

    // Save the file
    if err := os.WriteFile(filePath, []byte(profilePicture), 0644); err != nil {
        return "", fmt.Errorf("failed to save profile picture: %v", err)
    }

    return filePath, nil
}
