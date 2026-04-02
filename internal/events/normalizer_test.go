package events_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/realtong/logi-cli/internal/devices/mxmaster4"
	"github.com/realtong/logi-cli/internal/events"
)

func TestNormalizeHoldMoveDownFixture(t *testing.T) {
	reports := mustLoadFixtureReports(t, "thumb-button-hold-move-down.txt")
	stream := mustDecodeFixtureStream(t, reports)

	got := events.Normalize(stream, events.NormalizeConfig{GestureDistance: 1000})

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonDown && event.Control == "thumb_button"
	}) {
		t.Fatalf("normalized stream = %#v, want thumb_button down", got)
	}

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Kind == events.ButtonHold && event.Control == "thumb_button"
	}) {
		t.Fatalf("normalized stream = %#v, want thumb_button hold", got)
	}

	if !containsEvent(got, func(event events.DeviceEvent) bool {
		return event.Gesture == "hold(thumb_button)+move(down)"
	}) {
		t.Fatalf("normalized stream = %#v, want hold(thumb_button)+move(down)", got)
	}
}

func mustDecodeFixtureStream(t *testing.T, reports []events.RawReport) []events.DeviceEvent {
	t.Helper()

	adapter := mxmaster4.Adapter{}
	decoded := make([]events.DeviceEvent, 0, len(reports))
	for _, report := range reports {
		event, err := adapter.Decode(report)
		if err != nil {
			t.Fatalf("Decode returned error: %v", err)
		}
		if event.Kind == "" && event.Gesture == "" && event.Control == "" {
			continue
		}
		decoded = append(decoded, event)
	}
	return decoded
}

func containsEvent(eventsList []events.DeviceEvent, predicate func(events.DeviceEvent) bool) bool {
	for _, event := range eventsList {
		if predicate(event) {
			return true
		}
	}
	return false
}

func mustLoadFixtureReports(t *testing.T, name string) []events.RawReport {
	t.Helper()

	path := filepath.Join("..", "..", "testdata", "mxmaster4", name)
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
