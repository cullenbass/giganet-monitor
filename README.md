# giganet-monitor
Simple internet monitoring in golang.

Requires [go-httpstat](github.com/tcnksm/go-httpstat)

Running with no arguments fetches google.com once a minute. URL and interval can be configured with arguments.

Usage: `giganet-monitor <url to hit> <minutes between tests>`

Example: `giganet-monitor http://www.example.com 5` 