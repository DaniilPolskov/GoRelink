# GoRelink – Go URL Shortener

A simple URL shortener written in Go with both a console interface and a basic HTML web UI.

## Features

* Accepts a URL from the terminal or via web form.
* Automatically prepends `https://` if missing.
* Generates a random short identifier.
* Stores short → original URL mappings in memory.
* Exposes an `/api/shorten` endpoint for web clients.
* Redirects users from short URLs using `/gorelink/{id}`.
* Includes a basic HTML frontend with a logo.

---

## Usage (Web UI)

1. Open `index.html` in your browser or host it via Go.
2. Enter a URL into the form.
3. The shortened link will appear below the form.
4. It uses `POST /api/shorten` to generate short URLs.

---

## API

### `POST /api/shorten`

**Request Body:**

```json
{
  "url": "https://example.com"
}
```

**Response:**

```json
{
  "shortURL": "http://localhost:8080/gorelink/XyZ123"
}
```

CORS is enabled for local testing.

---

## Storage

Currently uses in-memory map. Data will be lost when the program stops.

## Possible Improvements

* Store links in a file or database.
* Add a web form for submitting URLs.
* Generate collision-free short IDs.
* Support expiration dates for short links.

## Testing

This project includes unit tests for the HTTP handlers and storage logic.

To run tests, use:

```bash
go test ./...
```

This will run all tests and show verbose output.

---

## Run Instructions

```bash
git clone https://github.com/DaniilPolskov/GoRelink.git
cd GoRelink
go run main.go
```

