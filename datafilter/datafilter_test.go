package datafilter

import (
	"testing"
)

func TestExecute(t *testing.T) {
	type args struct {
		data        string
		serviceName string
	}
	tests := []struct {
		name  string
		panic bool
		args  args
	}{
		{
			name:  "test_sqli_filter",
			panic: true,
			args: args{
				data:        "admin' AND 1=1 --",
				serviceName: "test",
			}},
		{
			name:  "test_xss_filter",
			panic: true,
			args: args{
				data:        "http://testing.com/book.html?default=<script>alert(document.cookie)</script>",
				serviceName: "test",
			}},
		{
			name:  "test_pan_filter",
			panic: true,
			args: args{
				data:        "44044333322221111swfkjbfjksjkf4444333322221111dedeefefefe",
				serviceName: "test",
			}},
		{
			name:  "test_no_match",
			panic: false,
			args: args{
				data:        "test data",
				serviceName: "test",
			}},
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
			Execute(tt.args.data, tt.args.serviceName)
		})
	}
}
