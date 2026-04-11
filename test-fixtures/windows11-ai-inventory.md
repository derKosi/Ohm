# Ohm — Windows 11 AI Inventory (Test Fixture)

Source: Ohm scan on Windows 11 development machine
Date: 2026-04-11
Purpose: Test case for Ohm scan target detection (Windows platform)

---

## Summary

Total AI software footprint: **~6.6 GB** across 31 items

## Ohm Scan Output (--no-tui)

```
C:\AI> ohm scan --no-tui

🤖 Agents & Harnesses (9 found, 2.6 GB)
   ⚠️  pi (Coding Agent)              15.3 MB    C:\Users\user\.pi
     - @mariozechner/pi-coding-agent
     - pi-session-search
     - @thesethrose/pi-zai-provider
     - @alexanderfortin/pi-zai-usage
     - @samfp/pi-memory
   🔑 Claude Code                    407.2 MB   C:\Users\user\.claude, C:\Users\user\.claude.json
     - @anthropic-ai/claude-code
      Aider                          0 B        (binary found in PATH)
   🔑 Codex CLI (OpenAI)             138.6 MB   C:\Users\user\.codex
     - @openai/codex
   🔑 Gemini CLI (Google)            725.3 MB   C:\Users\user\.gemini
     - @google/gemini-cli
      Mistral Vibe                   1.9 MB     C:\Users\user\.vibe
   ⚠️  PaperclipAI                    0 B        C:\Users\user\.paperclip
     - paperclipai
      OpenCode                       0 B        (binary found in PATH)
     - opencode-ai
   🔑 Continue                       1.4 GB     C:\Users\user\.continue

🖥️ AI Editors & IDEs (1 found, 82.4 KB)
   ⚠️  Cursor IDE                     82.4 KB    C:\Users\user\.cursor

⚙️ Model Runtimes (1 found, 0 B)
      LM Studio                      0 B

📦 SDKs & Frameworks (3 found, 0 B)
      HuggingFace Transformers       0 B
     - transformers, datasets, tokenizers
      OpenAI SDK                     0 B
   ⚠️  Playwright (AI-adjacent)       0 B
     - playwright

💾 Model Caches (1 found, 80 B)
      HuggingFace Hub Cache          80 B       C:\Users\user\.cache\huggingface

📄 Agent Config & Instructions (1 found, 10.7 KB)
   ⚠️  Agent Instruction Files        10.7 KB    (scattered across projects)
     - C:\Users\user\.claude\CLAUDE.md
     - C:\Users\user\.gemini\GEMINI.md
     - C:\Users\user\.pi\agent\AGENTS.md
     - C:\Users\user\.vscode\extensions\publisher.extension-v1.0.0\CLAUDE.md
     - C:\Users\user\.agent\worktrees\feature-1234\.cursorrules

🧠 Agent Memory & Sessions (4 found, 1.3 GB)
   🔑 Claude Code Memory             407.1 MB   C:\Users\user\.claude
      Mistral Vibe History           6.6 KB     C:\Users\user\.vibe\vibehistory
   🔑 PaperclipAI Context            213.3 MB   C:\Users\user\.paperclip\context.json, C:\Users\user\.paperclip\instances
   ⚠️  Gemini CLI History             725.3 MB   C:\Users\user\.gemini

🔌 MCP Configurations (1 found, 22 B)
   🔑 MCP Configuration Files        22 B
     - C:\Users\user\.lmstudio\mcp.json

🧩 Plugins & Extensions (1 found)
     - pi skill: paperclip
     - pi skill: paperclip-create-agent
     - pi skill: paperclip-create-plugin
     - pi skill: para-memory-files

📁 Config & Data Dirs (9 found, 2.6 GB)
   🔑 Claude Config                  407.1 MB
   🔑 Claude Config (JSON)           53.6 KB
   ⚠️  pi Config                      15.3 MB
   ⚠️  Codex CLI Config               138.6 MB
   ⚠️  Cursor Config                  82.4 KB
   ⚠️  Gemini CLI Config              725.3 MB
   🔑 Continue Config                1.4 GB
      Mistral Vibe Config            1.9 MB
   🔑 PaperclipAI Config             0 B

Total: 31 items (6.6 GB)
```

## Notes

- Aider found via PATH binary but no config dir (clean install)
- OpenCode found via PATH binary but no config dir
- LM Studio detected but size calculation returned 0 (Windows path issue — research needed)
- Continue (VS Code extension) stores 1.4 GB in its config dir
- Agent detected via .cursorrules in worktree
- Extension falsely matched for CLAUDE.md (contains CLAUDE in path)
