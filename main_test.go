package main

import "testing"

func Test_parsePoint(t *testing.T) {
	cases := []struct {
		name    string
		s       string
		want    Point
		wantErr bool
	}{
		{
			name:    "valid",
			s:       "1,2",
			want:    Point{X: 1, Y: 2},
			wantErr: false,
		},
		{
			name:    "no-comma",
			s:       "12",
			want:    Point{},
			wantErr: true,
		},
		{
			name:    "not-a-number",
			s:       "1,a",
			want:    Point{},
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parsePoint(c.s)
			if (err != nil) != c.wantErr {
				t.Errorf("parsePoint() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("parsePoint() = %v, want %v", got, c.want)
			}
		})
	}
}
