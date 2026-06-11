package stream

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type FFmpegRunner struct {
	ffmpegPath    string
	inputURL      string
	rtspTransport string
	videoCodec    string
	output        io.Writer
}

func NewFFmpegRunner(ffmpegPath string, inputURL string, rtspTransport string, videoCodec string, output io.Writer) *FFmpegRunner {
	return &FFmpegRunner{
		ffmpegPath:    ffmpegPath,
		inputURL:      inputURL,
		rtspTransport: rtspTransport,
		videoCodec:    strings.TrimSpace(videoCodec),
		output:        output,
	}
}

func (r *FFmpegRunner) RunLoop(ctx context.Context) {
	for {
		if err := r.runOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("ffmpeg stopped: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
			log.Println("restarting ffmpeg")
		}
	}
}

func (r *FFmpegRunner) runOnce(ctx context.Context) error {
	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-rtsp_transport", r.rtspTransport,
		"-i", r.inputURL,
		"-an",
	}

	args = append(args, r.videoArgs()...)
	args = append(args,
		"-f", "flv",
		"pipe:1",
	)

	cmd := exec.CommandContext(ctx, r.ffmpegPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	log.Println("ffmpeg started")
	_, copyErr := io.Copy(r.output, stdout)
	waitErr := cmd.Wait()

	if copyErr != nil {
		return copyErr
	}
	if waitErr != nil {
		return fmt.Errorf("%w: %s", waitErr, strings.TrimSpace(stderr.String()))
	}

	return nil
}

func (r *FFmpegRunner) videoArgs() []string {
	if strings.EqualFold(r.videoCodec, "copy") {
		return []string{"-c:v", "copy"}
	}

	codec := r.videoCodec
	if codec == "" {
		codec = "libx264"
	}

	return []string{
		"-c:v", codec,
		"-preset", "veryfast",
		"-tune", "zerolatency",
		"-pix_fmt", "yuv420p",
		"-g", "30",
		"-bf", "0",
	}
}
