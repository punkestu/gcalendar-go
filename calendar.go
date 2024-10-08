package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Calendar struct {
	Summary string
	Events  []*calendar.Event
}

type Calendars struct {
	Calendars   []Calendar
	eventsCount uint
}

func GetCalendars(client *http.Client, verbose bool) Calendars {
	ctx := context.Background()

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calList := getCalendars(srv)

	var calendars []Calendar
	eventsCount := 0
	if verbose {
		fmt.Printf("\033[31m=== Reading %d calendars ===\n\033[0m", len(calList.Items))
	}
	for _, cal := range calList.Items {
		events := getEvents(srv, cal)
		if events == nil {
			continue
		}
		calendars = append(calendars, *events)
		eventsCount += len(events.Events)
	}

	return Calendars{Calendars: calendars, eventsCount: uint(eventsCount)}
}

func getCalendars(srv *calendar.Service) *calendar.CalendarList {
	calList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of calendars: %v", err)
	}
	return calList
}

func getEvents(srv *calendar.Service, cal *calendar.CalendarListEntry) *Calendar {
	tmin := time.Now()
	tmax := tmin.Add(time.Duration(1 * 24 * 60 * 60 * 1000 * 1000 * 1000))
	events, err := srv.Events.List(cal.Id).ShowDeleted(false).
		SingleEvents(true).TimeMin(tmin.Format(time.RFC3339)).TimeMax(tmax.Format(time.RFC3339)).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	if len(events.Items) == 0 {
		return nil
	}
	return &Calendar{Summary: cal.Summary, Events: events.Items}
}
