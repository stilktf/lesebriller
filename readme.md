# Lesebriller (Reading Glasses)

Go-based alternative server for syncing book progress with KOReader

## Features
- Improved security (uses a more secure password hashing algorithm)
- Lightweight
- Easy to use in containers
  - Image works with arm64 as well as amd64
  - 15MB container (Alpine image)
- Uses sqlite
- Only 1 binary

## Set up using Docker

```
docker run -v ./db:/app/db -p 7200:7200 ghcr.io/stilktf/lesebriller -d
```