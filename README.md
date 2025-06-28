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

It's recommended to add the server using the `claude` command-line tool.

Run the following command in your terminal:
```bash
claude mcp add voice-notify go run github.com/kyong0612/voice-notify-mcp@latest
```

This will register the server under the name `voice-notify`.

Alternatively, you can manually add the following to your MCP server configuration file (e.g., `.claude.json` or `.mcp.json`):

```json
{
  "voice-notify": {
    "command": "go",
    "args": ["run", "github.com/kyong0612/voice-notify-mcp@latest"]
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

The AI is designed to use voice notifications autonomously to keep you informed without you needing to constantly check its progress. Here are some scenarios where you can expect a notification:

- **Long-running task completions**: When tasks like builds, tests, deployments, file downloads, or data analysis take more than a few seconds, the AI will notify you upon completion.
  - *Voice: "Build complete."*
- **User input required**: If the AI needs your approval, input, or a decision to proceed, it will alert you.
  - *Voice: "Approval required. Shall I proceed?"*
- **Errors or issues**: You'll be immediately notified if a process fails or an important warning occurs.
  - *Voice: "Deployment failed. Please check the logs."*
- **Key milestones**: For multi-step tasks, the AI will announce when it reaches an important checkpoint.
  - *Voice: "Staging deployment complete. Ready for final review."*

The AI may also proactively ask if you'd like a voice notification for a long-running task it's about to start.

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

### Debug Mode
Enable debug mode to see detailed logs:

```json
{
  "voice-notify": {
    "command": "go",
    "args": ["run", "github.com/kyong0612/voice-notify-mcp@latest"],
    "env": {
      "VOICE_NOTIFY_DEBUG": "true"
    }
  }
}
```

Debug logs include:
- Environment configuration at startup
- MCP request/response details  
- Voice selection process
- Language detection results
- Rate limiting decisions
- Command execution details

### No voice output
- Ensure your Mac's volume is not muted
- Check if the specified voice is installed
- Verify the `say` command works: `say "test"`
- Enable debug mode to see detailed error messages

### Voice not found
- The server will fall back to the system default voice
- Install additional voices in System Preferences → Accessibility → Spoken Content
- Use debug mode to see which voices are available

### Notifications during quiet hours
- Check your `VOICE_NOTIFY_QUIET_HOURS` setting
- Format should be "HH:MM-HH:MM" (24-hour format)
- Debug mode will show quiet hour calculations

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