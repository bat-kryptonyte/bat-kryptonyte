package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const username = "bat-kryptonyte"

type Event struct {
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		panic(err)
	}

	currentStreak, longestStreak := calculateStreak(events)

	badgeURL := fmt.Sprintf("https://img.shields.io/badge/Current_Streak-%d-brightgreen", currentStreak)
	badgeResp, err := http.Get(badgeURL)
	if err != nil {
		panic(err)
	}
	defer badgeResp.Body.Close()

	badge, err := ioutil.ReadAll(badgeResp.Body)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("streak_badge.svg", badge, 0644)

	updateReadme()
}

func calculateStreak(events []Event) (int, int) {
	var currentStreak, longestStreak int
	streakMap := make(map[string]bool)

	for _, event := range events {
		date := event.CreatedAt.Format("2006-01-02")
		streakMap[date] = true
	}

	var prevDate string
	for date := range streakMap {
		if prevDate == "" {
			prevDate = date
			currentStreak = 1
		} else {
			tPrev, _ := time.Parse("2006-01-02", prevDate)
			tCurr, _ := time.Parse("2006-01-02", date)
			if tCurr.Sub(tPrev).Hours() == 24 {
				currentStreak++
			} else {
				if currentStreak > longestStreak {
					longestStreak = currentStreak
				}
				currentStreak = 1
			}
			prevDate = date
		}
	}

	if currentStreak > longestStreak {
		longestStreak = currentStreak
	}

	return currentStreak, longestStreak
}

func updateReadme() {
	content, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	readme := string(content)
	badgeMarkdown := "![Streak](streak_badge.svg)"

	if !strings.Contains(readme, badgeMarkdown) {
		readme += "\n\n" + badgeMarkdown
	}

	ioutil.WriteFile("README.md", []byte(readme), 0644)
}
