package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage(w io.Writer) {
	fmt.Println(`NAME:
   org2gcal - Convert time-log to json format for gcal

Usage:
   org2gcal [date]

DATE:
   format	2006-1-2
   note		This argument is optional. if you do not specify this, date is used today.`)
}

var (
	fileName   string
	calendarID string
	nIO, errIO io.Writer
)

func main() {
	nIO, errIO = os.Stdin, os.Stdout

	fileName = "events.json"
	calendarID = os.Getenv("CALENDAR_ID")
	if calendarID == "" {
		calendarID = "primary"
	}

	var date string
	switch len(os.Args) {
	case 1:
		date = time.Now().Format("2006-1-2")
	case 2:
		if os.Args[1] == "--help" {
			usage(nIO)
			os.Exit(0)
		}
		date = os.Args[1]
	default:
		failed(fmt.Sprintf("Invalid argument number. got=%v, want=%v", len(os.Args), "1 or 2"))
	}

	t, err := time.ParseInLocation("2006-1-2", date, time.Local)
	if err != nil {
		failed(err.Error())
	}

	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	w := bufio.NewWriter(f)
	defer func() {
		w.Flush()
		f.Close()
	}()

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(body), "\n")

	var events []Event
	var startTime *time.Time
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		log := strings.SplitN(line, " ", 3)[1:]
		clock := strings.Split(log[0], ":")
		h, err := strconv.Atoi(clock[0])
		if err != nil {
			failed(err.Error())
		}
		m, err := strconv.Atoi(clock[1])
		if err != nil {
			failed(err.Error())
		}
		endTime := t.Add(time.Hour*time.Duration(h) + time.Minute*time.Duration(m))
		if startTime == nil {
			startTime = &endTime
		}
		events = append(events, NewEvent(log[1], *startTime, endTime))
		startTime = &endTime
	}

	if d, err := json.Marshal(events); err != nil {
		failed(err.Error())
	} else {
		w.Write(d)
	}
}

func failed(message string) {
	fmt.Fprintln(errIO, message)
	usage(nIO)
	os.Exit(1)
}

func NewEvent(summary string, start, end time.Time) Event {
	return Event{
		CalendarID: calendarID,
		Resource: EventResource{
			Summary: summary,
			Start:   NewEventTime(start),
			End:     NewEventTime(end),
		},
	}
}

func NewEventTime(t time.Time) EventTime {
	return EventTime{
		DateTime: t.Format("2006-01-02T15:04:05-0700"),
	}
}

type Event struct {
	CalendarID string        `json:"calendarId"`
	Resource   EventResource `json:"resource"`
}

type EventResource struct {
	Summary     string    `json:"summary,omitempty"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	Start       EventTime `json:"start,omitempty"`
	End         EventTime `json:"end,omitempty"`
}

type EventTime struct {
	DateTime string `json:"dateTime,omitempty"`
}
