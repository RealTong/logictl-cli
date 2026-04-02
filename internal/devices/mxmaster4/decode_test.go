package mxmaster4

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/events"
)

func TestDecodeThumbButtonDownFixture(t *testing.T) {
	reports := mustLoadFixtureReports(t, "thumb-button-down.txt")

	event, err := Decode(reports[0])
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if event.Kind != events.ButtonDown || event.Control != "thumb_button" {
		t.Fatalf("event = %#v, want thumb_button down", event)
	}
	if !event.At.Equal(reports[0].At) {
		t.Fatalf("event.At = %v, want %v", event.At, reports[0].At)
	}
}

func TestAdapterDecodeThumbButtonDownFixtureEmitsThumbButtonUpOnStableIdle(t *testing.T) {
	reports := mustLoadFixtureReports(t, "thumb-button-down.txt")
	adapter := Adapter{}

	var releaseCount int
	var releaseAt time.Time

	for index, report := range reports {
		event, err := adapter.Decode(report)
		if err != nil {
			if errors.Is(err, ErrUnsupportedReport) {
				continue
			}
			t.Fatalf("Decode returned error: %v", err)
		}
		if event.Kind != events.ButtonUp || event.Control != "thumb_button" {
			continue
		}
		releaseCount++
		releaseAt = event.At
		if index != 1 {
			t.Fatalf("Decode() emitted thumb_button up at report %d, want stable idle report 1", index)
		}
	}

	if releaseCount != 1 {
		t.Fatalf("releaseCount = %d, want 1", releaseCount)
	}
	if !releaseAt.Equal(reports[1].At) {
		t.Fatalf("releaseAt = %v, want %v", releaseAt, reports[1].At)
	}
}

func TestAdapterDecodeHoldMoveDownFixture(t *testing.T) {
	reports := mustLoadFixtureReports(t, "thumb-button-hold-move-down.txt")

	adapter := Adapter{}

	var sawDown bool
	var maxDeltaY int

	for _, report := range reports {
		event, err := adapter.Decode(report)
		if err != nil {
			if errors.Is(err, ErrUnsupportedReport) {
				continue
			}
			t.Fatalf("Decode returned error: %v", err)
		}
		if event.Kind == events.ButtonDown && event.Control == "thumb_button" {
			sawDown = true
		}
		if event.Kind == events.PointerMove && event.Control == "pointer" && event.DeltaY > maxDeltaY {
			maxDeltaY = event.DeltaY
		}
	}

	if !sawDown {
		t.Fatal("decoded stream did not contain a thumb_button down event")
	}
	if maxDeltaY <= 0 {
		t.Fatalf("decoded stream max DeltaY = %d, want positive downward motion", maxDeltaY)
	}
}

func TestAdapterDecodeHoldMoveDownFixtureDoesNotEmitThumbButtonUp(t *testing.T) {
	reports := mustLoadFixtureReports(t, "thumb-button-hold-move-down.txt")
	adapter := Adapter{}

	for index, report := range reports {
		event, err := adapter.Decode(report)
		if err != nil {
			if errors.Is(err, ErrUnsupportedReport) {
				continue
			}
			t.Fatalf("Decode returned error: %v", err)
		}
		if event.Kind == events.ButtonUp && event.Control == "thumb_button" {
			t.Fatalf("Decode() emitted thumb_button up at report %d: %#v", index, event)
		}
	}
}

func TestAdapterDecodeRejectsUnsupportedState(t *testing.T) {
	adapter := Adapter{}

	_, err := adapter.Decode(events.RawReport{
		DeviceID: "mx-master-4",
		Bytes:    []byte{0x02, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	})
	if !errors.Is(err, ErrUnsupportedReport) {
		t.Fatalf("Decode() error = %v, want ErrUnsupportedReport", err)
	}
}

func mustLoadFixtureReports(t *testing.T, name string) []events.RawReport {
	t.Helper()

	path := filepath.Join("..", "..", "..", "testdata", "mxmaster4", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", path, err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	reports := make([]events.RawReport, 0, len(lines))
	for _, line := range lines {
		report, err := parseFixtureReportLine(line)
		if err != nil {
			t.Fatalf("parseFixtureReportLine(%q) returned error: %v", line, err)
		}
		reports = append(reports, report)
	}

	return reports
}

func parseFixtureReportLine(line string) (events.RawReport, error) {
	fields := strings.Fields(line)
	if len(fields) != 9 {
		return events.RawReport{}, strconv.ErrSyntax
	}

	at, err := time.Parse(time.RFC3339Nano, fields[0])
	if err != nil {
		return events.RawReport{}, err
	}

	bytes := make([]byte, 0, len(fields)-1)
	for _, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 16, 8)
		if err != nil {
			return events.RawReport{}, err
		}
		bytes = append(bytes, byte(value))
	}

	return events.RawReport{
		DeviceID: "mx-master-4",
		At:       at,
		Bytes:    bytes,
	}, nil
}
