# Voice Notify MCP Server

A Model Context Protocol (MCP) server that enables AI assistants to send voice notifications on macOS. The AI can autonomously decide when to notify users about task completions, errors, or when attention is needed.

## Features

- 🎙️ Voice notifications using macOS `say` command
- 🌍 Automatic language detection for appropriate voice selection
- 🤖 Autonomous AI notifications (no explicit user instruction needed)
- 🔕 Quiet hours support
- 🎯 Priority-based notifications
- 🚀 Easy installation without cloning the repository

## Requirements

- macOS (uses the built-in `say` command)
- Go 1.21 or later
- Claude Desktop, Claude Code, Cursor, or Windsurf

## Installation

### Quick Start (Recommended)

You can run the server directly without cloning the repository:

```bash
go run github.com/kyong0612/voice-notify-mcp@latest
```

### Local Installation

1. Clone the repository:
```bash
git clone https://github.com/kyong0612/voice-notify-mcp.git
cd voice-notify-mcp
```

2. Run the server:
```bash
go run main.go
```

## Configuration

### Claude Desktop (Desktop Extensions)

1. Install the extension by adding the `dxt.json` file to your Claude Desktop extensions
2. Or manually configure in Claude Desktop settings

### Claude Code / Cursor / Windsurf

Add to your MCP server configuration:

```json
{
  "voice-notify": {
    "command": "go",
    "args": ["run", "github.com/kyong0612/voice-notify-mcp@latest"],
    "env": {
      "VOICE_NOTIFY_DEFAULT_VOICE": "Samantha",
      "VOICE_NOTIFY_DEFAULT_LANGUAGE": "en",
      "VOICE_NOTIFY_AUTO_NOTIFY": "true",
      "VOICE_NOTIFY_MIN_TASK_DURATION": "3",
      "VOICE_NOTIFY_QUIET_HOURS": "22:00-07:00"
    }
  }
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VOICE_NOTIFY_DEFAULT_VOICE` | Default voice name (e.g., "Samantha", "Kyoko") | System default |
| `VOICE_NOTIFY_DEFAULT_LANGUAGE` | Default language code (e.g., "en", "ja") | "en" |
| `VOICE_NOTIFY_AUTO_DETECT_LANGUAGE` | Enable automatic language detection | "true" |
| `VOICE_NOTIFY_AUTO_NOTIFY` | Enable autonomous AI notifications | "true" |
| `VOICE_NOTIFY_MIN_TASK_DURATION` | Minimum task duration (seconds) for auto-notification | "3" |
| `VOICE_NOTIFY_QUIET_HOURS` | Quiet hours range (e.g., "22:00-07:00") | None |

## Usage Examples

### Autonomous Notifications

The AI will automatically notify you in situations like:

- Long-running task completions (> 3 seconds)
- When user approval or input is needed
- Error occurrences
- Important milestones in multi-step processes

### Manual Notifications

You can also explicitly ask for voice notifications:

```
User: "Notify me when the analysis is complete"
AI: "I'll send a voice notification when the analysis finishes."
[Later] *Voice notification: "Analysis completed"*
```

### Language Support

The server automatically detects the language of the notification message and selects an appropriate voice:

- English: "Task completed" → English voice
- Japanese: "タスクが完了しました" → Japanese voice
- French: "Tâche terminée" → French voice

## Available Voices

To see available voices on your system:

```bash
say -v '?'
```

Common voices include:
- English: Alex, Samantha, Daniel
- Japanese: Kyoko, Otoya
- French: Amelie, Thomas
- Spanish: Monica, Jorge

## Troubleshooting

### No voice output
- Ensure your Mac's volume is not muted
- Check if the specified voice is installed
- Verify the `say` command works: `say "test"`

### Voice not found
- The server will fall back to the system default voice
- Install additional voices in System Preferences → Accessibility → Spoken Content

### Notifications during quiet hours
- Check your `VOICE_NOTIFY_QUIET_HOURS` setting
- Format should be "HH:MM-HH:MM" (24-hour format)

## Development

### Building from source

```bash
go build -o voice-notify-mcp
./voice-notify-mcp
```

### Running tests

```bash
go test ./...
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Built with [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- Inspired by the Model Context Protocol specification