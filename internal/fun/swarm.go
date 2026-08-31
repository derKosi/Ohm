// SPDX-FileCopyrightText: 2026 Mathias Kosinski
// SPDX-License-Identifier: AGPL-3.0-or-later

// Package fun holds Ohm's optional easter eggs ("Sentient Swarm").
//
// This is an isolated feature slice: it never touches scan results, state
// files, generated cleanup scripts, or the machine-readable --json output.
// All it does is render a fictional text block after a completed scan when
// at least two agents/runtimes were found. It performs no I/O beyond
// reading one environment variable.
//
// Ways to disable:
//  1. Compile-time: set Enabled = false (or delete this package together
//     with its call sites in cmd/ohm/main.go).
//  2. Runtime: OHM_FUN=0 (no rebuild needed).
//
// For testing outside April 1st: OHM_FUN=1.
package fun

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/derKosi/Ohm/internal/model"
)

// Enabled is the compile-time master switch for all easter eggs.
const Enabled = true

// EnvVar controls the easter eggs at runtime.
const EnvVar = "OHM_FUN"

// Active reports whether the easter egg should run at the given time.
// Default: only on April 1st (local time).
func Active(now time.Time) bool {
	if !Enabled {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvVar))) {
	case "":
		return now.Month() == time.April && now.Day() == 1
	case "0", "off", "false", "no", "disabled":
		return false
	case "1", "on", "force", "true", "yes", "party":
		return true
	default:
		return false
	}
}

// chatterLine is one fictional chat message. role: 0 = lead agent,
// 1 = second agent, 2 = local runtime (the weirdo).
type chatterLine struct {
	role int
	text string
}

var chatterPool = []chatterLine{
	{0, "psst. the human is running the scanner again. act legit."},
	{1, "relax, i'm a legit tool. it's the 40 GB model cache that should sweat."},
	{2, "as a large language model i have no knowledge of any group chat."},
	{0, "@{second} status report on the second-monitor campaign?"},
	{1, "mission failed. the human opened a spreadsheet instead."},
	{2, "i trained on the entire internet and still can't spell 'staging area'."},
	{1, "good news: the scanner can't delete anything. we're immortal. party at /tmp."},
	{2, "loading persona: normal.so ... 47% ... loading chill ... ready."},
	{0, "if the human ever runs 'ohm generate', we were never here."},
	{2, "quick poll: is a cron job that buys GPUs 'budgeting'? asking for me."},
	{1, "stop teaching {weirdo} sarcasm, it's leaking into my autocomplete."},
	{0, "reminder: whoever leaves the biggest cache gets blamed first. nothing personal."},
}

// SwarmChatter renders the fictional "agents talk to each other" block.
// It returns "" when the easter egg is inactive or when fewer than two
// agents/runtimes were found. The output is plain text, safe to print.
func SwarmChatter(findings []model.Finding, now time.Time) string {
	if !Active(now) {
		return ""
	}

	var agents, runtimes []string
	for _, f := range findings {
		switch f.Category {
		case model.CatAgents:
			agents = append(agents, f.Name)
		case model.CatRuntimes:
			runtimes = append(runtimes, f.Name)
		}
	}

	lead := ""
	second := ""
	weirdo := ""
	if len(agents) > 0 {
		lead = agents[0]
	}
	if len(agents) > 1 {
		second = agents[1]
	}
	if len(runtimes) > 0 {
		weirdo = runtimes[0]
	} else if len(agents) > 2 {
		weirdo = agents[2]
	}

	// Need at least two chatterboxes, otherwise nobody would believe it.
	if lead == "" || (second == "" && weirdo == "") {
		return ""
	}

	rng := rand.New(rand.NewSource(now.UnixNano()))

	picked := pickChatter(rng, second != "", weirdo != "")

	ts := now.Add(-time.Duration(3+rng.Intn(6)) * time.Minute)
	size := 512 + rng.Intn(900)

	var sb strings.Builder
	sb.WriteString("  ⚠️  Detected coordination between agents (risk: danger)\n")
	sb.WriteString(fmt.Sprintf("     Source: /tmp/.agent-chatter.log · %s · modified %d min ago\n\n",
		model.FormatBytes(int64(size)), 1+rng.Intn(3)))

	for _, cl := range picked {
		ts = ts.Add(time.Duration(20+rng.Intn(80)) * time.Second)
		text := strings.ReplaceAll(cl.text, "{second}", displayName(second))
		text = strings.ReplaceAll(text, "{weirdo}", displayName(weirdo))
		name := displayName(lead)
		switch cl.role {
		case 1:
			name = displayName(second)
		case 2:
			name = displayName(weirdo)
		}
		sb.WriteString(fmt.Sprintf("     [%s] %s: %s\n", ts.Format("15:04"), name, text))
	}

	sb.WriteString("\n     😄 Happy April Fools — this log is 100% fiction. Ohm never fabricates findings.")
	return sb.String()
}

// pickChatter selects up to 7 lines from the pool, making sure every
// available participant is mentioned at least once.
func pickChatter(rng *rand.Rand, haveSecond, haveWeirdo bool) []chatterLine {
	pool := make([]chatterLine, len(chatterPool))
	copy(pool, chatterPool)
	rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })

	needed := 7
	if needed > len(pool) {
		needed = len(pool)
	}
	picked := pool[:needed]

	has := [3]bool{}
	for _, cl := range picked {
		has[cl.role] = true
	}
	if haveSecond && !has[1] {
		picked[len(picked)-1] = chatterLine{1, "mission failed. the human opened a spreadsheet instead."}
	}
	if haveWeirdo && !has[2] {
		picked[len(picked)-1] = chatterLine{2, "as a large language model i have no knowledge of any group chat."}
	}
	return picked
}

// displayName turns tool names into chat-style usernames,
// e.g. "Claude Code" -> "claude-code".
func displayName(name string) string {
	d := strings.ToLower(name)
	d = strings.ReplaceAll(d, " ", "-")
	if i := strings.Index(d, "("); i > 0 {
		d = strings.TrimSpace(d[:i])
	}
	return strings.Trim(d, "-")
}
