package schedule

import (
	"github.com/kaanaktas/go-slm/cache"
	"strconv"
	"testing"
	"time"
)

func Test_isCurrentTimeInScheduledTime(t *testing.T) {
	type args struct {
		startTime string
		duration  time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "shouldReturnTrue_24_Hour_duration",
			args: args{
				startTime: time.Now().Format("2006-01-02") + "T00:00:00",
				duration:  1440,
			},
			want: true,
		},
		{
			name: "shouldReturnTrue_12_hour_duration",
			args: args{
				startTime: time.Now().Format("2006-01-02") + "T" + prepareTime(time.Now().Hour(), 0),
				duration:  720,
			},
			want: true,
		},
		{
			name: "shouldReturnTrue_5_minute_duration",
			args: args{
				startTime: time.Now().Format("2006-01-02") + "T" + prepareTime(time.Now().Hour(), 0),
				duration:  time.Duration(time.Now().Minute() + 5),
			},
			want: true,
		},
		{
			name: "shouldReturnFalse_5_minute_duration",
			args: args{
				startTime: time.Now().Format("2006-01-02") + "T" + prepareTime(time.Now().Hour(), 0),
				duration:  time.Duration(time.Now().Minute() - 5),
			},
			want: false,
		},
		{
			name: "shouldReturnFalse_exceed_the_same_minute",
			args: args{
				startTime: time.Now().Format("2006-01-02") + "T" + prepareTime(time.Now().Hour(), 0),
				duration:  time.Duration(time.Now().Minute()),
			},
			want: false,
		},
		{
			name: "shouldReturnFalse_before_the_1_minute_start",
			args: args{
				startTime: time.Now().Format("2006-01-02") + "T" + prepareTime(time.Now().Hour(), time.Now().Minute()+1),
				duration:  1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCurrentTimeInScheduledTime(tt.args.startTime, tt.args.duration); got != tt.want {
				t.Errorf("isCurrentTimeInScheduledTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func prepareTime(hour int, min int) string {
	hourStr := strconv.Itoa(hour)
	if hour < 10 {
		hourStr = "0" + hourStr
	}

	minStr := strconv.Itoa(min)
	if min < 10 {
		minStr = "0" + minStr
	}

	return hourStr + ":" + minStr + ":00"
}

func TestLoad(t *testing.T) {
	cacheIn := cache.NewInMemory()
	cacheIn.Flush()

	Load("/testdata/schedule.yaml")

	type args struct {
		scheduleName string
		daysSize     int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "weekdays_schedule",
			args: args{
				scheduleName: "weekend",
				daysSize:     2,
			},
		},
		{
			name: "weekend_schedule",
			args: args{
				scheduleName: "weekdays",
				daysSize:     5,
			},
		},
	}

	expectedCacheSize := 2

	if cachedData, ok := cacheIn.Get(key); ok {
		scheduleCache := cachedData.([]schedule)
		if len(scheduleCache) != expectedCacheSize {
			t.Errorf("cached data size didn't match up. Expected: %d, got:%d", expectedCacheSize,
				len(scheduleCache))
		}

		for i, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				data := scheduleCache[i]
				if data.ScheduleName != tt.args.scheduleName {
					t.Errorf("scheduleName didn't match up. Expected: %s, got:%s",
						tt.args.scheduleName, data.ScheduleName)
				}

				if len(data.Days) != tt.args.daysSize {
					t.Errorf("daysSize didn't match up. Expected: %d, got:%d",
						tt.args.daysSize, len(data.Days))
				}
			})
		}
	} else {
		t.Error("schedule Statements is not in the cache")
	}
}
