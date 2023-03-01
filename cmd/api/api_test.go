package api

import "testing"

func TestGetURL(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"foo path", args{path: "/foo"}, "https://rebrickable.com/api/v3/foo"},
		{"foo path", args{path: "foo"}, "https://rebrickable.com/api/v3/foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetURL(tt.args.path); got != tt.want {
				t.Errorf("GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
