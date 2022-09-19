package schedule

import (
	"fmt"
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/config"
	"github.com/kaanaktas/go-slm/policy"
	"path/filepath"
	"sort"
	"time"
)

var cacheIn = cache.NewInMemory()

const (
	key        = "schedule_rule"
	timeLayout = "2006-01-02T15:04:05"
	dateLayout = "2006-01-02"
)

type schedule struct {
	ScheduleName string   `yaml:"scheduleName"`
	Days         []string `yaml:"days"`
	Start        string   `yaml:"start"`
	Duration     int      `yaml:"duration"`
	Message      string   `yaml:"message"`
}

type Executor struct {
	Actions []policy.Action
}

func (e *Executor) Apply() {
	cachedSchedule, ok := cacheIn.Get(key)
	if !ok {
		panic("schedule doesn't exist in the cache")
	}

	schedules := cachedSchedule.([]schedule)
	sort.Slice(e.Actions, func(i, j int) bool {
		return e.Actions[i].Order < e.Actions[j].Order
	})

	for _, action := range e.Actions {
		if action.Active {
			for _, sc := range schedules {
				if sc.ScheduleName == action.Name {
					if ok := isScheduleMatchWithPolicy(sc); ok {
						panic(sc.Message)
					}
				}
			}
		}
	}
}

func Load(scheduleStatementPath string) {
	if scheduleStatementPath == "" {
		panic("GO_SLM_SCHEDULE_STATEMENT_PATH hasn't been set")
	}
	content := config.MustReadFile(filepath.Join(config.RootDirectory, scheduleStatementPath))
	var schedules []schedule
	config.MustUnmarshalYaml(scheduleStatementPath, content, &schedules)

	cacheIn.Set(key, schedules, cache.NoExpiration)
}

func isScheduleMatchWithPolicy(sc schedule) bool {
	if isScheduledDayActive(sc.Days) {
		return isCurrentTimeInScheduledTime(generateStartTime(sc.Start), time.Duration(sc.Duration))
	}
	return false
}

func generateStartTime(start string) string {
	return time.Now().Format(dateLayout) + "T" + start
}

func isScheduledDayActive(days []string) bool {
	dayOfTheWeek := time.Now().Weekday().String()
	for _, day := range days {
		if dayOfTheWeek == day {
			return true
		}
	}
	return false
}

func isCurrentTimeInScheduledTime(startTime string, duration time.Duration) bool {
	start, err := time.ParseInLocation(timeLayout, startTime, time.Local)
	if err != nil {
		panic(fmt.Sprintf("Error during parsing the time %s", err))
	}
	end := start.Add(time.Minute * duration)
	current := time.Now()
	return current.After(start) && current.Before(end)
}
