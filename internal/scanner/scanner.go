package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/platform"
)

// Options controls what the scanner looks for.
type Options struct {
	ScanPATH     bool
	ScanENV      bool
	ScanShell    bool
	ScanDeep     bool
	CustomSigDir string
}

// Scanner discovers AI software on the system.
type Scanner struct {
	os       platform.OS
	home     string
	opts     Options
	findings []model.Finding
	mu       sync.Mutex
}

// New creates a new Scanner.
func New(opts Options) *Scanner {
	return &Scanner{
		os:   platform.Detect(),
		home: platform.HomeDir(),
		opts: opts,
	}
}

// Scan runs all scanners and returns findings.
func (s *Scanner) Scan() *model.ScanResult {
	s.findings = nil

	hostname, _ := os.Hostname()

	// Run category scanners
	s.scanAgents()
	s.scanEditors()
	s.scanRuntimes()
	s.scanComfyUI()
	s.scanSDKs()
	s.scanModelCaches()
	s.scanInstructions()
	s.scanMemory()
	s.scanMCP()
	s.scanPlugins()
	s.scanVSCodeExtensions()
	s.scanConfigDirs()
	s.scanDocker()
	s.scanStragglers()

	// Opt-in scanners
	if s.opts.ScanPATH {
		s.scanPATH()
	}
	if s.opts.ScanENV {
		s.scanENV()
	}
	if s.opts.ScanShell {
		s.scanShellProfiles()
	}

	return &model.ScanResult{
		Findings:  s.findings,
		ScannedAt: time.Now(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		Hostname:  hostname,
	}
}

// addFinding safely adds a finding.
func (s *Scanner) addFinding(f model.Finding) {
	s.mu.Lock()
	s.findings = append(s.findings, f)
	s.mu.Unlock()
}

// finding creates a basic finding with common defaults.
func (s *Scanner) finding(id, name string, category model.Category) model.Finding {
	return model.Finding{
		ID:       id,
		Name:     name,
		Category: category,
	}
}

// dirExists checks if a directory exists.
func (s *Scanner) dirExists(path string) bool {
	return platform.IsDir(path)
}

// fileExists checks if a file exists.
func (s *Scanner) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// dirSize gets directory size.
func (s *Scanner) dirSize(path string) int64 {
	return platform.DirSize(path)
}

// fileSize gets file size.
func (s *Scanner) fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// expandPath resolves ~ to home directory.
func (s *Scanner) expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		return filepath.Join(s.home, path[1:])
	}
	return path
}

// hasCommand checks if a command exists.
func (s *Scanner) hasCommand(name string) bool {
	return platform.HasCommand(name)
}

// runCommand runs a command and returns output.
func (s *Scanner) runCommand(name string, args ...string) (string, error) {
	return platform.RunCommand(name, args...)
}
