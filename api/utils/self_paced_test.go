package utils

import (
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
)

func TestParseSessionMode(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"", constants.ModeLive, false},
		{"live", constants.ModeLive, false},
		{"self_paced", constants.ModeSelfPaced, false},
		{"SELF_PACED", "", true},
		{"async", "", true},
	}
	for _, tc := range cases {
		got, err := ParseSessionMode(tc.in)
		if (err != nil) != tc.wantErr {
			t.Fatalf("ParseSessionMode(%q) err = %v, wantErr %v", tc.in, err, tc.wantErr)
		}
		if got != tc.want {
			t.Fatalf("ParseSessionMode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestIsSelfPacedTimedOut(t *testing.T) {
	delivered := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	duration := 30

	if IsSelfPacedTimedOut(delivered.Add(29*time.Second), delivered, duration) {
		t.Fatal("answer within duration must not time out")
	}
	if IsSelfPacedTimedOut(delivered.Add(31*time.Second), delivered, duration) {
		t.Fatal("answer within grace period must not time out")
	}
	if IsSelfPacedTimedOut(delivered.Add(32*time.Second), delivered, duration) {
		t.Fatal("answer exactly at the deadline must not time out")
	}
	if !IsSelfPacedTimedOut(delivered.Add(32*time.Second+time.Millisecond), delivered, duration) {
		t.Fatal("answer after the deadline must time out")
	}
	want := delivered.Add(32 * time.Second)
	if got := SelfPacedDeadline(delivered, duration); !got.Equal(want) {
		t.Fatalf("SelfPacedDeadline = %v, want %v", got, want)
	}
}

func TestSelfPacedResponseTimeMs(t *testing.T) {
	delivered := time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		now  time.Time
		want int
	}{
		{"clock skew before delivery", delivered.Add(-time.Second), 0},
		{"instant", delivered, 0},
		{"mid", delivered.Add(12345 * time.Millisecond), 12345},
		{"in grace period", delivered.Add(31 * time.Second), 30000},
	}
	for _, tc := range cases {
		if got := SelfPacedResponseTimeMs(tc.now, delivered, 30); got != tc.want {
			t.Fatalf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestIsAnswerCorrect(t *testing.T) {
	if IsAnswerCorrect(nil, []int{1}, constants.SingleAnswer) {
		t.Fatal("empty selection is never correct")
	}
	if IsAnswerCorrect([]int{}, []int{1, 2}, constants.Survey) {
		t.Fatal("empty survey selection is not correct")
	}
	if !IsAnswerCorrect([]int{3}, []int{1, 2, 3}, constants.Survey) {
		t.Fatal("any survey selection is correct")
	}
	if !IsAnswerCorrect([]int{2}, []int{2}, constants.SingleAnswer) {
		t.Fatal("matching single answer must be correct")
	}
	if IsAnswerCorrect([]int{1}, []int{2}, constants.SingleAnswer) {
		t.Fatal("wrong single answer must not be correct")
	}
}
