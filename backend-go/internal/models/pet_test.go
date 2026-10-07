package models_test

import (
	"testing"
	"time"

	"github.com/RealKiro/learnstar-planet/backend-go/internal/models"
)

func TestLevelFromScore(t *testing.T) {
	cases := []struct {
		score int
		want  int
	}{
		{0, 1}, {1, 1}, {14, 1},
		{15, 2}, {34, 2},
		{35, 3}, {59, 3},
		{60, 4}, {89, 4},
		{90, 5}, {124, 5},
		{125, 6}, {164, 6},
		{165, 7}, {209, 7},
		{210, 8}, {259, 8},
		{260, 9}, {314, 9},
		{315, 10}, {374, 10},
		{375, 11}, {449, 11},
		{450, 12}, {9999, 12},
	}
	for _, c := range cases {
		if got := models.LevelFromScore(c.score); got != c.want {
			t.Errorf("LevelFromScore(%d) = %d, want %d", c.score, got, c.want)
		}
	}
}

func TestCurrentStage(t *testing.T) {
	cases := []struct {
		level int
		stage string
		emoji string
	}{
		{1, "egg", "🥚"},
		{2, "egg", "🥚"},
		{3, "baby", "🐣"},
		{4, "baby", "🐣"},
		{5, "growing", "🌱"},
		{6, "growing", "🌱"},
		{7, "mature", "🌿"},
		{8, "mature", "🌿"},
		{9, "legendary", "🌟"},
		{10, "legendary", "🌟"},
		{11, "transcendent", "👑"},
		{12, "transcendent", "👑"},
	}
	for _, c := range cases {
		p := models.Pet{Level: c.level}
		got := p.CurrentStage()
		if got.Stage != c.stage || got.Emoji != c.emoji {
			t.Errorf("CurrentStage(level %d) = (%s, %s), want (%s, %s)",
				c.level, got.Stage, got.Emoji, c.stage, c.emoji)
		}
	}
}

func TestAddExperienceLevelsUp(t *testing.T) {
	p := models.Pet{Level: 1, Experience: 0}
	p.AddExperience(25) // Lv1 需 20，剩 5 升到 Lv2
	if p.Level != 2 || p.Experience != 5 {
		t.Fatalf("AddExperience(25) = level %d exp %d, want level 2 exp 5", p.Level, p.Experience)
	}
}

func TestAddExperienceCapsAtMaxLevel(t *testing.T) {
	p := models.Pet{Level: 12, Experience: 0}
	p.AddExperience(1000)
	if p.Level != models.MaxLevel {
		t.Fatalf("level = %d, want %d", p.Level, models.MaxLevel)
	}
}

func TestRemoveExperienceLevelsDown(t *testing.T) {
	p := models.Pet{Level: 2, Experience: 0}
	p.RemoveExperience(1)
	if p.Level != 1 || p.Experience != 19 {
		t.Fatalf("RemoveExperience(1) = level %d exp %d, want level 1 exp 19", p.Level, p.Experience)
	}
}

func TestRemoveExperienceFloorsAtZero(t *testing.T) {
	p := models.Pet{Level: 0, Experience: 0}
	p.RemoveExperience(10)
	if p.Level != 0 || p.Experience != 0 {
		t.Fatalf("RemoveExperience(10) = level %d exp %d, want level 0 exp 0", p.Level, p.Experience)
	}
}

func TestFeedCapsMood(t *testing.T) {
	p := models.Pet{Mood: 90}
	now := time.Now()
	p.Feed(now)
	if p.Mood != 100 {
		t.Fatalf("mood = %d, want 100", p.Mood)
	}
	if p.LastFedAt == nil || !p.LastFedAt.Equal(now) {
		t.Fatalf("LastFedAt not recorded correctly")
	}
}

func TestDecayMood(t *testing.T) {
	now := time.Now()

	// 满 24 小时衰减 10。
	p := models.Pet{Mood: 80, LastFedAt: timePtr(now.Add(-25 * time.Hour))}
	p.DecayMood(now)
	if p.Mood != 70 {
		t.Fatalf("mood = %d, want 70", p.Mood)
	}

	// 未满 24 小时不衰减。
	p2 := models.Pet{Mood: 80, LastFedAt: timePtr(now.Add(-23 * time.Hour))}
	p2.DecayMood(now)
	if p2.Mood != 80 {
		t.Fatalf("mood = %d, want 80", p2.Mood)
	}

	// 衰减后不为负。
	p3 := models.Pet{Mood: 5, LastFedAt: timePtr(now.Add(-25 * time.Hour))}
	p3.DecayMood(now)
	if p3.Mood != 0 {
		t.Fatalf("mood = %d, want 0", p3.Mood)
	}
}

func TestSwitchCost(t *testing.T) {
	cases := []struct {
		level int
		want  int
	}{
		{0, 5}, {1, 5}, {2, 10}, {12, 60},
	}
	for _, c := range cases {
		if got := models.SwitchCost(c.level); got != c.want {
			t.Errorf("SwitchCost(%d) = %d, want %d", c.level, got, c.want)
		}
	}
}

func timePtr(t time.Time) *time.Time { return &t }
