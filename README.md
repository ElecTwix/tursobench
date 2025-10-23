# Turso Go Concurrency Bug

Repository: https://github.com/ElecTwix/tursobench

## Issue

Turso Go driver v0.2.2 produces "database is locked" errors with 4+ concurrent goroutines, despite Turso supporting concurrent writes.

## Error Details

- 1-2 workers: Works but slow
- 4+ workers: "database is locked" errors
- Performance degrades instead of improving
- Timeouts with multiple concurrent operations

## Environment

- Turso Go Driver: v0.2.2
- Go: 1.21+
- OS: Linux
- Concurrency: 4+ goroutines triggers the issue

## Reproduction

```bash
git clone https://github.com/ElecTwix/tursobench
cd tursobench
go run simple_concurrent.go
```
