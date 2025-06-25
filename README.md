# Go URL Shortener

A simple console-based URL shortener with a redirect server built using Go.

## Features

- Accepts a URL from the user via the console.
- Automatically adds `https://` if missing.
- Generates a random short identifier.
- Stores the mapping between short and original URLs in memory.
- Starts a local HTTP server to handle redirection.
- Redirects users to the original URL when they visit the shortened one.

## Example

```

Enter URL to shorten: google.com
Old URL: [https://google.com](https://google.com)
Short URL: [http://localhost:8080/gorelink/1234567890](http://localhost:8080/gorelink/1234567890)

```

Visiting the short URL in your browser redirects you to the original link.

## How to Run

1. Make sure you have Go installed.
2. Clone the repository.
3. Run the application:

```bash
go run main.go
```

4. Enter a URL when prompted in the terminal.
5. Open the printed short URL in your browser to test the redirect.

## Notes

* This version uses in-memory storage (`map[string]string`). All data is lost when the program stops.
* It’s intended for demonstration or learning purposes.
* The short ID is randomly generated and not guaranteed to be unique. In production, you'd want to check for duplicates or use a UUID or hash.

## Possible Improvements

* Store links in a file or database.
* Add a web form for submitting URLs.
* Generate collision-free short IDs.
* Support expiration dates for short links.
