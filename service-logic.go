// Логика сервиса бронирования номеров

package main

import (
	"fmt"
	"log"
	"time"
)

func daysBetween(from time.Time, to time.Time, uemail string) []time.Time {
	if from.After(to) {
		log.Println("From before To --- should be") // Проверка дат. Check from < to
		return nil
	}

	// Проверка наличия дат в мапе. Checks dates of user
	for _, v := range cache {
		if v.From == from && v.To == to && uemail == v.UserEmail {
			log.Printf("Checks the order is booked: %v\n", v)
			log.Println("Hotel room is booked the user for selected dates")
			return nil
		}
	}

	days := make([]time.Time, 0)
	for d := toDay(from); !d.After(toDay(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func toDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

var logger = log.Default()

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
