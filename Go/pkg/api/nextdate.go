package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Format = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if now.Format(Format) == "" {
		now = time.Now().Truncate(24 * time.Hour)
	}

	if dstart == "" {
		dstart = now.Format(Format)
	}

	dstartParse, err := time.Parse(Format, dstart)
	if err != nil {
		return "", err
	}

	switch {
	case repeat == "y":
		for {
			dstartParse = dstartParse.AddDate(1, 0, 0)
			if dstartParse.After(now) {
				return dstartParse.Format(Format), nil
			}

		}
	case strings.HasPrefix(repeat, "d "):

		repeatSplit := strings.Split(repeat, " ")
		if len(repeatSplit) < 2 {
			return "", fmt.Errorf("wrong day format")
		}

		dayNum, err := strconv.Atoi(repeatSplit[1])
		if err != nil {
			return "", err
		}
		if dayNum > 400 {
			return "", fmt.Errorf("Wrong day number")
		}

		for {
			dstartParse = dstartParse.AddDate(0, 0, dayNum)
			if dstartParse.After(now) {
				return dstartParse.Format(Format), nil
			}
		}
	case repeat == "":
		return "", nil
	default:
		return "", fmt.Errorf("wrong data input")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	if now == "" {
		now = time.Now().Format(Format)
	}

	nowTime, err := time.Parse(Format, now)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if date == "" {
		date = now
		return
	}

	result, err := NextDate(nowTime, date, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}
