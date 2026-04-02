package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/realtong/logi-cli/internal/events"
	"github.com/realtong/logi-cli/internal/hidapi"
	"github.com/spf13/cobra"
)

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
		Short: "Print raw HID events",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !raw {
				return errors.New("normalized event output is not implemented yet; rerun with --raw")
			}

			resolvedPath, err := resolveEventDevicePath(hidClient, path)
			if err != nil {
				return err
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			return streamRawReports(ctx, openSource(resolvedPath), cmd.OutOrStdout(), output)
		},
	}

	cmd.Flags().BoolVar(&raw, "raw", true, "print raw HID reports")
	cmd.Flags().StringVar(&path, "path", "", "optional HID device path to capture")
	cmd.Flags().StringVar(&output, "output", "", "optional path for captured reports")
	return cmd
}

func resolveEventDevicePath(hidClient hidapi.Client, explicitPath string) (string, error) {
	if explicitPath != "" {
		return explicitPath, nil
	}

	devices, err := hidClient.ListDevices()
	if err != nil {
		return "", err
	}

	switch len(devices) {
	case 0:
		return "", errors.New("no HID devices available")
	case 1:
		if devices[0].Path == "" {
			return "", errors.New("selected HID device is missing a path")
		}
		return devices[0].Path, nil
	default:
		return "", errors.New("multiple HID devices found; rerun with --path")
	}
}

func streamRawReports(ctx context.Context, source rawSource, out io.Writer, outputPath string) error {
	writer := out
	if writer == nil {
		writer = io.Discard
	}

	var captureFile *os.File
	if outputPath != "" {
		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		captureFile = file
		defer func() {
			_ = captureFile.Close()
		}()
		writer = io.MultiWriter(writer, captureFile)
	}

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
