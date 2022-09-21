package schedule

import (
	"github.com/kaanaktas/go-slm/cache"
	"github.com/kaanaktas/go-slm/policy"
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
			name: "weekend_schedule",
			args: args{
				scheduleName: "weekend",
				daysSize:     2,
			},
		},
		{
			name: "weekdays_schedule",
			args: args{
				scheduleName: "weekdays_between_8_18",
				daysSize:     5,
			},
		},
		{
			name: "weekdays_schedule_all_day",
			args: args{
				scheduleName: "weekdays_all_day",
				daysSize:     5,
			},
		},
		{
			name: "no_days_match",
			args: args{
				scheduleName: "no_days_match",
				daysSize:     0,
			},
		},
	}

	expectedCacheSize := 4

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

func TestExecutor_Apply(t *testing.T) {
	cacheIn := cache.NewInMemory()
	cacheIn.Flush()

	Load("/testdata/schedule.yaml")

	type fields struct {
		Actions []policy.Action
	}
	tests := []struct {
		name   string
		panic  bool
		fields fields
	}{
		{
			name:  "test_schedule_not_permitted",
			panic: true,
			fields: fields{Actions: []policy.Action{
				{Name: "weekdays_all_day", Active: true, Order: 10},
			}}},
		{
			name:  "test_schedule_not_active",
			panic: false,
			fields: fields{Actions: []policy.Action{
				{Name: "weekdays_all_day", Active: false, Order: 10},
			}}},
		{
			name:  "empty_schedule",
			panic: false,
			fields: fields{Actions: []policy.Action{
				{Name: "no_days_match", Active: true, Order: 10},
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) && tt.panic == false {
					t.Errorf("%s did panic", tt.name)
				} else if (r == nil) && tt.panic == true {
					t.Errorf("%s didn't panic", tt.name)
				}
			}()

			e := &Executor{
				Actions: tt.fields.Actions,
			}
			e.Apply()
		})
	}
}
