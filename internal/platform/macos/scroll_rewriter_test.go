package macos

import (
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/config"
)

func TestScrollMatcherOnlyConsumesRecentMatchingWheelTicks(t *testing.T) {
	matcher := newScrollMatcher(50 * time.Millisecond)
	matcher.Record("mx-master-4", "wheel_up", config.ScrollConfig{
		Direction:    "standard",
		SmoothScroll: true,
	}, time.Unix(1, 0))

	plan, ok := matcher.Match(nativeScrollEvent{
		VerticalLine:  1,
		VerticalPoint: 12,
		At:            time.Unix(1, 10*1e6),
	})
	if !ok {
		t.Fatal("Match returned false, want matched wheel event")
	}
	if !plan.Settings.SmoothScroll || plan.Settings.Direction != "standard" {
		t.Fatalf("plan.Settings = %#v, want standard + smooth true", plan.Settings)
	}
	if got, want := len(plan.Emissions), 4; got != want {
		t.Fatalf("len(plan.Emissions) = %d, want %d", got, want)
	}
	if plan.Emissions[0].Unit != scrollUnitPixel {
		t.Fatalf("plan.Emissions[0].Unit = %v, want pixel", plan.Emissions[0].Unit)
	}
	if plan.Emissions[0].Vertical >= 0 {
		t.Fatalf("plan.Emissions[0].Vertical = %d, want inverted negative delta", plan.Emissions[0].Vertical)
	}
}

func TestScrollMatcherRejectsExpiredWheelTicks(t *testing.T) {
	matcher := newScrollMatcher(50 * time.Millisecond)
	matcher.Record("mx-master-4", "wheel_up", config.ScrollConfig{
		Direction: "standard",
	}, time.Unix(1, 0))

	if _, ok := matcher.Match(nativeScrollEvent{
		VerticalLine: 1,
		At:           time.Unix(1, 100*1e6),
	}); ok {
		t.Fatal("Match returned true, want expired wheel tick to be ignored")
	}
}

func TestScrollMatcherKeepsHorizontalAndVerticalIndependent(t *testing.T) {
	matcher := newScrollMatcher(50 * time.Millisecond)
	matcher.Record("mx-master-4", "thumb_wheel_left", config.ScrollConfig{
		Direction:    "standard",
		SmoothScroll: false,
	}, time.Unix(1, 0))

	if _, ok := matcher.Match(nativeScrollEvent{
		VerticalLine: 1,
		At:           time.Unix(1, 10*1e6),
	}); ok {
		t.Fatal("Match returned true, want vertical native scroll to ignore horizontal HID tick")
	}

	plan, ok := matcher.Match(nativeScrollEvent{
		HorizontalLine: 1,
		At:             time.Unix(1, 20*1e6),
	})
	if !ok {
		t.Fatal("Match returned false, want horizontal native scroll to match")
	}
	if len(plan.Emissions) != 1 {
		t.Fatalf("len(plan.Emissions) = %d, want 1", len(plan.Emissions))
	}
	if plan.Emissions[0].Horizontal == 0 {
		t.Fatalf("plan.Emissions[0] = %#v, want non-zero horizontal delta", plan.Emissions[0])
	}
}
