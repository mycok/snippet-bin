package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T)  {
	testCases := []struct{
		name string
		input time.Time
		output string
	}{
		{
			name: "UTC",
			input: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			output: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			input: time.Time{},
			output: "",
		},
		{
			name: "CET",
			input: time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			output: "17 Dec 2020 at 09:00",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hd := humanReadableDate(tc.input)

			if hd != tc.output {
				t.Errorf("expected %q; got %q", tc.output, hd)
			}
		})
	}
}