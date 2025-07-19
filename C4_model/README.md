# 1. Уровень 1: Контекст

```mermaid
graph LR
    A[Пользователь] -->|CLI/TUI| B[GophKeeper Client]
    B -->|HTTPS/gRPC| C[GophKeeper Server]
    C -->|SQL| D[(PostgreSQL)]
```