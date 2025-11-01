# Contributing Guidelines

## Development Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Copy `.env.example` to `.env`
4. Start PostgreSQL: `docker-compose up -d postgres`
5. Run migrations: The app will auto-migrate on startup
6. Start the server: `go run cmd/server/main.go`

## Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for public functions
- Keep functions focused and small

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Commit Messages

- Use clear and descriptive messages
- Follow conventional commits format
- Reference issues when applicable

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request with a clear description




