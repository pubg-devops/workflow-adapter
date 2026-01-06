package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// Printer handles formatted output for CLI commands
type Printer struct {
	out io.Writer
	err io.Writer
}

// NewPrinter creates a new Printer with stdout and stderr
func NewPrinter() *Printer {
	return &Printer{
		out: os.Stdout,
		err: os.Stderr,
	}
}

// NewPrinterWithWriters creates a Printer with custom writers (useful for testing)
func NewPrinterWithWriters(out, errOut io.Writer) *Printer {
	return &Printer{
		out: out,
		err: errOut,
	}
}

// Table prints data in a formatted table
func (p *Printer) Table(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(p.out, 0, 0, 2, ' ', 0)

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	separators := make([]string, len(headers))
	for i, h := range headers {
		separators[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
}

// JSON prints data as formatted JSON
func (p *Printer) JSON(data interface{}) error {
	encoder := json.NewEncoder(p.out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Printf prints formatted output
func (p *Printer) Printf(format string, args ...interface{}) {
	fmt.Fprintf(p.out, format, args...)
}

// Println prints a line
func (p *Printer) Println(args ...interface{}) {
	fmt.Fprintln(p.out, args...)
}

// Errorf prints formatted error output
func (p *Printer) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(p.err, "Error: "+format+"\n", args...)
}

// Warnf prints formatted warning output
func (p *Printer) Warnf(format string, args ...interface{}) {
	fmt.Fprintf(p.err, "Warning: "+format+"\n", args...)
}

// FormatTimestamp converts a Unix timestamp string (milliseconds) to a readable format
func FormatTimestamp(timestampMs string) string {
	var ms int64
	if _, err := fmt.Sscanf(timestampMs, "%d", &ms); err != nil {
		return timestampMs
	}

	t := time.UnixMilli(ms)
	return t.Format("2006-01-02 15:04:05")
}

// TruncateString truncates a string to the specified length, adding "..." if truncated
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// FormatMessagePreview creates a preview of message content for display
func FormatMessagePreview(body string, maxLen int) string {
	// Remove newlines and extra whitespace for preview
	preview := strings.ReplaceAll(body, "\n", " ")
	preview = strings.ReplaceAll(preview, "\r", "")
	preview = strings.ReplaceAll(preview, "\t", " ")

	// Collapse multiple spaces
	for strings.Contains(preview, "  ") {
		preview = strings.ReplaceAll(preview, "  ", " ")
	}

	preview = strings.TrimSpace(preview)

	return TruncateString(preview, maxLen)
}

// PrettyJSON formats JSON string with indentation
func PrettyJSON(jsonStr string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// Not valid JSON, return as-is
		return jsonStr, nil
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return jsonStr, err
	}

	return string(pretty), nil
}
