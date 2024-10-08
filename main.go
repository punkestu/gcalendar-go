package main

import (
	"fmt"
)

func main() {
	client := getClient()

	fmt.Println("Listing the upcoming events:")
	calendars := GetCalendars(client, true)
	if calendars.eventsCount == 0 {
		fmt.Println("No upcoming events found.")
		return
	}
	for _, cal := range calendars.Calendars {
		fmt.Printf("\033[34m=== %s ===\n\033[0m", cal.Summary)
		for _, item := range cal.Events {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			dateEnd := item.End.DateTime
			if dateEnd == "" {
				dateEnd = item.End.Date
			}
			fmt.Printf("\033[31m%v  -  %v\033[0m >> %v\n", date, dateEnd, item.Summary)
		}
	}
}
