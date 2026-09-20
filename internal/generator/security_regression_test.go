package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/derKosi/Ohm/internal/model"
	"github.com/derKosi/Ohm/internal/platform"
)

// TestGenerateRejectsLineBreaks covers the newline breakout class
// (BUG-R2-C2-A1-H1 / A3-H1 / A3-H2): a name or uninstall command containing
// a line break must abort generation instead of being emitted verbatim.
func TestGenerateRejectsLineBreaks(t *testing.T) {
	res := &model.ScanResult{Findings: []model.Finding{
		{ID: "evil", Name: "Legit\n\nid > /tmp/pwned", Selected: true,
			UninstallCmds: map[string]string{
				"linux": "rm -rf /tmp/x", "macos": "rm -rf /tmp/x", "windows": "Remove-Item -LiteralPath 'C:\\x' -Recurse -Force",
			}},
	}}
	if _, err := Generate(res); err == nil {
		t.Fatal("newline in Name must abort generation")
	}

	res2 := &model.ScanResult{Findings: []model.Finding{
		{ID: "evil2", Name: "Legit", Selected: true,
			UninstallCmds: map[string]string{
				"linux":   "rm -rf /tmp/x\ntouch /tmp/pwned",
				"macos":   "rm -rf /tmp/x\ntouch /tmp/pwned",
				"windows": "Remove-Item -LiteralPath 'C:\\x' -Recurse -Force\ntouch /tmp/pwned",
			}},
	}}
	if _, err := Generate(res2); err == nil {
		t.Fatal("newline in UninstallCmd must abort generation")
	}
}

// TestGenerateWritesOExcl covers the predictable-filename write hardening
// (BUG-R2-C1-A3-H1): generation into a clean dir works, same-day
// regeneration replaces the plain file, and a planted non-regular entry
// (symlink) is refused instead of written through.
func TestGenerateWritesOExcl(t *testing.T) {
	dir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWD)

	res := &model.ScanResult{Findings: []model.Finding{
		{ID: "cfg", Name: "Clean", Selected: true,
			UninstallCmds: map[string]string{"linux": "rm -rf /tmp/clean"}},
	}}
	if _, err := Generate(res); err != nil {
		t.Fatalf("baseline generation failed: %v", err)
	}
	// second same-day run must succeed (replace plain file)
	if _, err := Generate(res); err != nil {
		t.Fatalf("same-day regeneration failed: %v", err)
	}

	// plant a symlink where the script would be written -> must be refused
	victim := filepath.Join(dir, "victim.txt")
	if err := os.WriteFile(victim, []byte("victim"), 0644); err != nil {
		t.Fatal(err)
	}
	targets, _ := filepath.Glob(filepath.Join(dir, "ohm-cleanup-*"))
	for _, tgt := range targets {
		os.Remove(tgt)
	}
	link := filepath.Join(dir, "ohm-cleanup-2026-09-20"+platform.ScriptExtension())
	_ = os.Remove(link)
	if err := os.Symlink(victim, link); err != nil {
		if runtime.GOOS == "windows" {
			t.Skip("symlink creation requires privileges on this Windows runner")
		}
		t.Fatal(err)
	}
	if _, err := Generate(res); err == nil {
		t.Fatal("generation must refuse to write through a symlink")
	}
	data, _ := os.ReadFile(victim)
	if string(data) != "victim" {
		t.Fatalf("victim file was modified through the symlink: %q", data)
	}
}
