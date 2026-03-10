package tui

import (
	"time"

	"github.com/SamNet-dev/findns/internal/scanner"
)

type progressMsg struct {
	stepIndex int
	done      int
	total     int
	passed    int
	failed    int
}

type scanDoneMsg struct {
	report   scanner.ChainReport
	elapsed  time.Duration
	err      error
	writeErr error
}

type scanStartedMsg struct {
	progressCh chan progressMsg
	doneCh     chan scanDoneMsg
}

type inputLoadedMsg struct {
	ips []string
	err error
}
