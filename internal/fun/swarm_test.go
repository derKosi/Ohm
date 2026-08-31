// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

package fun

import (
	"strings"
	"testing"
	"time"

	"github.com/derKosi/Ohm/internal/model"
)

func agent(name string) model.Finding {
	return model.Finding{ID: name, Name: name, Category: model.CatAgents}
}

func runtime(name string) model.Finding {
	return model.Finding{ID: name, Name: name, Category: model.CatRuntimes}
}

func TestInactiveOutsideApril1(t *testing.T) {
	findings := []model.Finding{agent("Claude Code"), agent("Codex CLI (OpenAI)")}
	now := time.Date(2026, 8, 31, 15, 0, 0, 0, time.Local)
	if got := SwarmChatter(findings, now); got != "" {
		t.Fatalf("expected empty output outside April 1st, got:\n%s", got)
	}
}

func TestActiveOnApril1(t *testing.T) {
	findings := []model.Finding{agent("Claude Code"), agent("Codex CLI (OpenAI)")}
	now := time.Date(2026, 4, 1, 2, 0, 0, 0, time.Local)
	got := SwarmChatter(findings, now)
	if got == "" {
		t.Fatal("expected output on April 1st, got none")
	}
	if !strings.Contains(got, "claude-code") || !strings.Contains(got, "codex-cli") {
		t.Fatalf("expected both agents in chat, got:\n%s", got)
	}
	if !strings.Contains(got, "fiction") {
		t.Fatalf("expected fictional disclaimer, got:\n%s", got)
	}
}

func TestForceOnViaEnv(t *testing.T) {
	t.Setenv(EnvVar, "1")
	findings := []model.Finding{agent("Claude Code"), runtime("Ollama")}
	now := time.Date(2026, 8, 31, 15, 0, 0, 0, time.Local)
	got := SwarmChatter(findings, now)
	if got == "" {
		t.Fatal("expected forced output, got none")
	}
	if !strings.Contains(got, "ollama") {
		t.Fatalf("expected runtime in chat, got:\n%s", got)
	}
}

func TestForceOffViaEnv(t *testing.T) {
	t.Setenv(EnvVar, "0")
	findings := []model.Finding{agent("Claude Code"), agent("Codex CLI (OpenAI)")}
	now := time.Date(2026, 4, 1, 2, 0, 0, 0, time.Local)
	if got := SwarmChatter(findings, now); got != "" {
		t.Fatalf("expected env kill switch to win, got:\n%s", got)
	}
}

func TestNeedsTwoChatterboxes(t *testing.T) {
	t.Setenv(EnvVar, "1")
	now := time.Date(2026, 4, 1, 2, 0, 0, 0, time.Local)
	if got := SwarmChatter([]model.Finding{agent("Claude Code")}, now); got != "" {
		t.Fatalf("expected no chatter with a single agent, got:\n%s", got)
	}
	if got := SwarmChatter([]model.Finding{agent("Claude Code"), agent("Claude Flow"), agent("Claude Squad")}, now); got == "" {
		t.Fatal("expected chatter with three agents")
	}
}

func TestOutputStableForFixedTime(t *testing.T) {
	t.Setenv(EnvVar, "1")
	findings := []model.Finding{agent("Claude Code"), agent("Codex CLI (OpenAI)"), runtime("Ollama")}
	now := time.Date(2026, 8, 31, 15, 0, 0, 0, time.Local)
	a := SwarmChatter(findings, now)
	b := SwarmChatter(findings, now)
	if a != b {
		t.Fatalf("output should be deterministic for a fixed scan time:\n---\n%s\n---\n%s", a, b)
	}
}
