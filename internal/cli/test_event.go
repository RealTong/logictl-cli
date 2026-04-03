package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/realtong/logi-cli/internal/devices/mxmaster4"
	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
	"github.com/spf13/cobra"
)

const logitechVendorID = 0x046d

type RawReport = events.RawReport

type rawSource interface {
	Stream(context.Context) (<-chan RawReport, <-chan error)
}

type rawSourceFactory func(path string) rawSource

func newTestCmd(hidClient hidapi.Client, openSource rawSourceFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Inspect live device events",
	}

	cmd.AddCommand(newTestEventDeviceCmd(hidClient, openSource))
	return cmd
}

func newTestEventCmd(source rawSource, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event",
		Short: "Print raw HID events",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return streamRawReports(cmd.Context(), source, out, "")
		},
	}

	return cmd
}

func newTestEventDeviceCmd(hidClient hidapi.Client, openSource rawSourceFactory) *cobra.Command {
	var path string
	var output string
	var raw bool

	cmd := &cobra.Command{
		Use:   "event",
		Short: "Print live device events",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if raw {
				resolvedPath, err := resolveEventDevicePath(hidClient, path)
				if err != nil {
					return err
				}

				ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
				defer stop()

				return streamRawReports(ctx, openSource(resolvedPath), cmd.OutOrStdout(), output)
			}

			resolvedPath, err := resolveSemanticEventDevicePath(hidClient, path)
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			return streamSemanticEvents(ctx, openSource(resolvedPath), cmd.OutOrStdout(), output)
		},
	}

	cmd.Flags().BoolVar(&raw, "raw", false, "print raw HID reports")
	cmd.Flags().StringVar(&path, "path", "", "optional HID device path to capture")
	cmd.Flags().StringVar(&output, "output", "", "optional path for captured reports")
	return cmd
}

func streamSemanticEvents(ctx context.Context, source rawSource, out io.Writer, outputPath string) error {
	writer, closeWriter, err := newEventWriter(out, outputPath)
	if err != nil {
		return err
	}
	defer closeWriter()

	adapter := mxmaster4.Adapter{}
	normalizer := events.NewNormalizer(events.NormalizeConfig{})
	reportedWarnings := map[string]struct{}{}

	reports, errs := source.Stream(ctx)
	for reports != nil || errs != nil {
		select {
		case report, ok := <-reports:
			if !ok {
				reports = nil
				continue
			}

			decoded, err := adapter.Decode(report)
			if err != nil {
				if errors.Is(err, mxmaster4.ErrUnsupportedReport) {
					if writeErr := writeWarningOnce(writer, reportedWarnings, semanticWarningKey("unsupported", err), formatUnsupportedSemanticReport(report, err)); writeErr != nil {
						return writeErr
					}
					continue
				}
				if writeErr := writeWarningOnce(writer, reportedWarnings, semanticWarningKey("ignored", err), formatIgnoredSemanticReport(report, err)); writeErr != nil {
					return writeErr
				}
				continue
			}
			for _, event := range decoded {
				for _, normalized := range normalizer.Push(event) {
					if _, err := fmt.Fprintln(writer, events.FormatDeviceEvent(normalized)); err != nil {
						return err
					}
				}
			}
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func writeWarningOnce(writer io.Writer, seen map[string]struct{}, key string, line string) error {
	if _, ok := seen[key]; ok {
		return nil
	}
	seen[key] = struct{}{}
	_, err := fmt.Fprintln(writer, line)
	return err
}

func semanticWarningKey(kind string, err error) string {
	return fmt.Sprintf("%s:%s", kind, err.Error())
}

func formatUnsupportedSemanticReport(report events.RawReport, err error) string {
	return fmt.Sprintf("unsupported_report %s (%v)", events.FormatRawReport(report), err)
}

func formatIgnoredSemanticReport(report events.RawReport, err error) string {
	return fmt.Sprintf("ignored_report %s (%v)", events.FormatRawReport(report), err)
}

func newEventWriter(out io.Writer, outputPath string) (io.Writer, func(), error) {
	writer := out
	if writer == nil {
		writer = io.Discard
	}

	if outputPath == "" {
		return writer, func() {}, nil
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return nil, nil, err
	}

	return io.MultiWriter(writer, file), func() {
		_ = file.Close()
	}, nil
}

func streamRawReports(ctx context.Context, source rawSource, out io.Writer, outputPath string) error {
	writer, closeWriter, err := newEventWriter(out, outputPath)
	if err != nil {
		return err
	}
	defer closeWriter()

	reports, errs := source.Stream(ctx)
	for reports != nil || errs != nil {
		select {
		case report, ok := <-reports:
			if !ok {
				reports = nil
				continue
			}
			if _, err := fmt.Fprintln(writer, events.FormatRawReport(report)); err != nil {
				return err
			}
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func resolveEventDevicePath(hidClient hidapi.Client, explicitPath string) (string, error) {
	if explicitPath != "" {
		return explicitPath, nil
	}

	devices, err := hidClient.ListDevices()
	if err != nil {
		return "", err
	}

	switch candidates := collapseSupportedEventCandidates(supportedEventCandidates(devices)); len(candidates) {
	case 0:
		if len(devices) == 0 {
			return "", errors.New("no HID devices available")
		}
		return "", errors.New("no supported Logitech HID devices available; rerun with --path")
	case 1:
		if candidates[0].Path == "" {
			return "", errors.New("selected HID device is missing a path")
		}
		return candidates[0].Path, nil
	default:
		return "", errors.New("multiple supported Logitech HID devices found; rerun with --path")
	}
}

func resolveSemanticEventDevicePath(hidClient hidapi.Client, explicitPath string) (string, error) {
	devices, err := hidClient.ListDevices()
	if err != nil {
		return "", err
	}

	adapter := mxmaster4.Adapter{}
	if explicitPath != "" {
		for _, device := range devices {
			if device.Path != explicitPath {
				continue
			}
			if !adapter.Matches(device) {
				return "", fmt.Errorf("unsupported semantic capture for HID path %q: only MX Master 4 is supported; rerun with --raw", explicitPath)
			}
			return explicitPath, nil
		}
		return "", fmt.Errorf("HID path %q is not currently available for semantic capture; connect an MX Master 4 or rerun with --raw", explicitPath)
	}

	switch candidates := collapseSupportedEventCandidates(supportedSemanticEventCandidates(devices, adapter)); len(candidates) {
	case 0:
		if len(devices) == 0 {
			return "", errors.New("no HID devices available")
		}
		return "", errors.New("no supported MX Master 4 HID devices available; rerun with --raw or --path")
	case 1:
		if candidates[0].Path == "" {
			return "", errors.New("selected HID device is missing a path")
		}
		return candidates[0].Path, nil
	default:
		return "", errors.New("multiple supported MX Master 4 HID devices found; rerun with --path")
	}
}

func supportedEventCandidates(devices []hidapi.DeviceInfo) []hidapi.DeviceInfo {
	candidates := make([]hidapi.DeviceInfo, 0, len(devices))
	for _, device := range devices {
		if !isSupportedEventCandidate(device) {
			continue
		}
		candidates = append(candidates, device)
	}
	return candidates
}

func supportedSemanticEventCandidates(devices []hidapi.DeviceInfo, adapter mxmaster4.Adapter) []hidapi.DeviceInfo {
	candidates := make([]hidapi.DeviceInfo, 0, len(devices))
	for _, device := range devices {
		if !adapter.Matches(device) {
			continue
		}
		candidates = append(candidates, device)
	}
	return candidates
}

func collapseSupportedEventCandidates(devices []hidapi.DeviceInfo) []hidapi.DeviceInfo {
	collapsed := make([]hidapi.DeviceInfo, 0, len(devices))
	seen := make(map[string]struct{}, len(devices))
	for _, device := range devices {
		key := supportedEventCandidateKey(device)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		collapsed = append(collapsed, device)
	}
	return collapsed
}

func supportedEventCandidateKey(device hidapi.DeviceInfo) string {
	if device.SerialNumber != "" {
		return fmt.Sprintf("serial:%04x:%04x:%s", device.VendorID, device.ProductID, device.SerialNumber)
	}
	return fmt.Sprintf("path:%s", device.Path)
}

func isSupportedEventCandidate(device hidapi.DeviceInfo) bool {
	if device.VendorID != logitechVendorID {
		return false
	}
	return strings.Contains(strings.ToLower(device.Product), "mx master")
}
