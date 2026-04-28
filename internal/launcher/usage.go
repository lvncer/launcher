package launcher

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
)

const recentLimit = 5

func loadUsage() map[string]usage {
	m := map[string]usage{}
	data, err := os.ReadFile(usageFile)
	if err != nil {
		return m
	}
	_ = json.Unmarshal(data, &m)
	return m
}

func saveUsage(m map[string]usage) {
	data, _ := json.Marshal(m)
	_ = os.WriteFile(usageFile, data, 0644)
}

func (m *model) updateUsage(key string) {
	u := m.usage[key]
	u.Count++
	u.Last = time.Now().Unix()
	m.usage[key] = u
	saveUsage(m.usage)
}

func sortByRecentThenName(items []item, usageByTitle map[string]usage) []item {
	sort.SliceStable(items, func(i, j int) bool {
		return compareTitle(items[i].title, items[j].title)
	})

	recent := recentItems(items, usageByTitle)
	sort.SliceStable(recent, func(i, j int) bool {
		left := usageByTitle[recent[i].title]
		right := usageByTitle[recent[j].title]

		if left.Last != right.Last {
			return left.Last > right.Last
		}
		return compareTitle(recent[i].title, recent[j].title)
	})

	if len(recent) > recentLimit {
		recent = recent[:recentLimit]
	}

	result := make([]item, 0, len(recent)+len(items))
	result = append(result, recent...)
	result = append(result, items...)
	return result
}

func recentItems(items []item, usageByTitle map[string]usage) []item {
	var recent []item

	for _, item := range items {
		if usageByTitle[item.title].Last > 0 {
			recent = append(recent, item)
		}
	}
	return recent
}

func compareTitle(left, right string) bool {
	l := strings.ToLower(left)
	r := strings.ToLower(right)
	if l != r {
		return l < r
	}
	return left < right
}
