package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// AuthResult holds the data for an auth check response.
type AuthResult struct {
	Status             string `json:"status"`
	User               string `json:"user"`
	RateLimitRemaining int    `json:"rate_limit_remaining"`
}

// PrintAuth writes the auth check result to w in text or JSON format.
func PrintAuth(w io.Writer, r AuthResult, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	}
	fmt.Fprintf(w, "auth: %s\n", r.Status)
	fmt.Fprintf(w, "user: %s\n", r.User)
	fmt.Fprintf(w, "rate_limit_remaining: %d\n", r.RateLimitRemaining)
	return nil
}

// PrintError writes an error message to w in text or JSON format.
func PrintError(w io.Writer, msg string, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]string{"error": msg})
		return
	}
	fmt.Fprintf(w, "error: %s\n", msg)
}

// EntryData represents a content entry for output formatting.
type EntryData struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"download_url,omitempty"`
}

// PrintEntry writes a single entry to w in text or JSON format.
func PrintEntry(w io.Writer, e EntryData, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(e)
	}
	fmt.Fprintf(w, "type: %s\n", e.Type)
	fmt.Fprintf(w, "path: %s\n", e.Path)
	fmt.Fprintf(w, "sha: %s\n", e.SHA)
	fmt.Fprintf(w, "size: %d\n", e.Size)
	if e.DownloadURL != "" {
		fmt.Fprintf(w, "download_url: %s\n", e.DownloadURL)
	}
	return nil
}

// PrintEntries writes a list of entries to w in text or JSON format.
func PrintEntries(w io.Writer, entries []EntryData, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\n", e.Type, e.Path)
	}
	return nil
}

// MutationResultData represents the result of a write or delete operation.
type MutationResultData struct {
	Action string `json:"action"` // "created", "updated", "deleted"
	Path   string `json:"path"`
	SHA    string `json:"sha"`    // commit SHA
	Branch string `json:"branch,omitempty"`
}

// PrintMutationResult writes a mutation result to w in text or JSON format.
func PrintMutationResult(w io.Writer, r MutationResultData, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	}
	fmt.Fprintf(w, "action: %s\n", r.Action)
	fmt.Fprintf(w, "path: %s\n", r.Path)
	fmt.Fprintf(w, "sha: %s\n", r.SHA)
	if r.Branch != "" {
		fmt.Fprintf(w, "branch: %s\n", r.Branch)
	}
	return nil
}
