// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package platform

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// OS represents the operating system.
type OS int

const (
	OSLinux OS = iota
	OSMacOS
	OSWindows
)

func (o OS) String() string {
	names := []string{"linux", "macos", "windows"}
	if int(o) < len(names) {
		return names[o]
	}
	return "unknown"
}

// Detect returns the current OS.
func Detect() OS {
	switch runtime.GOOS {
	case "linux":
		return OSLinux
	case "darwin":
		return OSMacOS
	case "windows":
		return OSWindows
	default:
		return OSLinux
	}
}

// HomeDir returns the user's home directory.
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

// ConfigDir returns the platform-specific config directory.
func ConfigDir() string {
	home := HomeDir()
	switch Detect() {
	case OSWindows:
		if appData := os.Getenv("APPDATA"); appData != "" {
			return appData
		}
		return filepath.Join(home, "AppData", "Roaming")
	case OSMacOS:
		return filepath.Join(home, ".config")
	default:
		return filepath.Join(home, ".config")
	}
}

// LocalAppData returns the Windows LocalAppData path, or empty string on other OSes.
func LocalAppData() string {
	if Detect() == OSWindows {
		return os.Getenv("LOCALAPPDATA")
	}
	return ""
}

// HasCommand checks if a command exists in PATH.
func HasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// RunCommand executes a command and returns its output.
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// DirSize calculates the total size of a directory.
func DirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// Exists checks if a path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir checks if a path is a directory.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// ScriptExtension returns the appropriate script extension for the current OS.
func ScriptExtension() string {
	switch Detect() {
	case OSWindows:
		return ".ps1"
	default:
		return ".sh"
	}
}

// ScriptHeader returns the script header for the current OS.
func ScriptHeader() string {
	switch Detect() {
	case OSWindows:
		return "# Ohm Cleanup Script (PowerShell)\n# Review before running!\n$ErrorActionPreference = 'Stop'\n"
	default:
		return "#!/usr/bin/env bash\n# Ohm Cleanup Script\n# Review before running!\nset -euo pipefail\n"
	}
}
