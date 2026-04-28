package launcher

import (
	"reflect"
	"testing"
)

func TestSortByRecentThenName(t *testing.T) {
	items := []item{
		{title: "Zoom"},
		{title: "Calendar"},
		{title: "Arc"},
		{title: "Docker"},
		{title: "Finder"},
		{title: "GitHub"},
		{title: "Notes"},
		{title: "Slack"},
	}
	usageByTitle := map[string]usage{
		"Zoom":     {Last: 10},
		"Calendar": {Last: 70},
		"Arc":      {Last: 20},
		"Finder":   {Last: 50},
		"GitHub":   {Last: 40},
		"Notes":    {Last: 60},
	}

	got := sortByRecentThenName(items, usageByTitle)
	gotTitles := titles(got)

	want := []string{
		"Calendar",
		"Notes",
		"Finder",
		"GitHub",
		"Arc",
		"Docker",
		"Slack",
		"Zoom",
	}

	if !reflect.DeepEqual(gotTitles, want) {
		t.Fatalf("sortByRecentThenName() = %v, want %v", gotTitles, want)
	}
}

func titles(items []item) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.title)
	}
	return result
}
