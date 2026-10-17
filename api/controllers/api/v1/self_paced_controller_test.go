package v1

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
)

func TestSelfPacedSessionError(t *testing.T) {
	cases := []struct {
		name    string
		session models.ActiveQuiz
		want    string
	}{
		{"live session", models.ActiveQuiz{Mode: constants.ModeLive, IsActive: true}, constants.ErrSessionNotSelfPaced},
		{"ended self-paced", models.ActiveQuiz{Mode: constants.ModeSelfPaced, IsActive: false}, constants.ErrSelfPacedQuizEnded},
		{"active self-paced", models.ActiveQuiz{Mode: constants.ModeSelfPaced, IsActive: true}, ""},
	}
	for _, tc := range cases {
		if got := selfPacedSessionError(tc.session); got != tc.want {
			t.Fatalf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestFirstUnanswered(t *testing.T) {
	if got := firstUnanswered(nil); got != -1 {
		t.Fatalf("empty rows: got %d, want -1", got)
	}
	rows := []models.SelfPacedResponseRow{
		{Answered: true},
		{Answered: false, DeliveredAt: sql.NullTime{Valid: true, Time: time.Now()}},
		{Answered: false},
	}
	if got := firstUnanswered(rows); got != 1 {
		t.Fatalf("got %d, want 1", got)
	}
	rows[1].Answered = true
	rows[2].Answered = true
	if got := firstUnanswered(rows); got != -1 {
		t.Fatalf("all answered: got %d, want -1", got)
	}
}
