package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPrintAuth_Text(t *testing.T) {
	var buf bytes.Buffer
	r := AuthResult{Status: "ok", User: "octocat", RateLimitRemaining: 4999}
	if err := PrintAuth(&buf, r, false); err != nil {
		t.Fatal(err)
	}
	want := "auth: ok\nuser: octocat\nrate_limit_remaining: 4999\n"
	if buf.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestPrintAuth_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := AuthResult{Status: "ok", User: "octocat", RateLimitRemaining: 4999}
	if err := PrintAuth(&buf, r, true); err != nil {
		t.Fatal(err)
	}
	var decoded AuthResult
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Status != "ok" || decoded.User != "octocat" || decoded.RateLimitRemaining != 4999 {
		t.Errorf("unexpected decoded result: %+v", decoded)
	}
}

func TestPrintError_Text(t *testing.T) {
	var buf bytes.Buffer
	PrintError(&buf, "something broke", false)
	want := "error: something broke\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}

func TestPrintError_JSON(t *testing.T) {
	var buf bytes.Buffer
	PrintError(&buf, "something broke", true)
	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded["error"] != "something broke" {
		t.Errorf("unexpected error field: %q", decoded["error"])
	}
}

func TestPrintEntry_Text(t *testing.T) {
	var buf bytes.Buffer
	e := EntryData{
		Type:        "file",
		Path:        "README.md",
		SHA:         "abc123",
		Size:        42,
		DownloadURL: "https://example.com/README.md",
	}
	if err := PrintEntry(&buf, e, false); err != nil {
		t.Fatal(err)
	}
	want := "type: file\npath: README.md\nsha: abc123\nsize: 42\ndownload_url: https://example.com/README.md\n"
	if buf.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestPrintEntry_JSON(t *testing.T) {
	var buf bytes.Buffer
	e := EntryData{
		Type: "file",
		Path: "README.md",
		SHA:  "abc123",
		Size: 42,
	}
	if err := PrintEntry(&buf, e, true); err != nil {
		t.Fatal(err)
	}
	var decoded EntryData
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Type != "file" || decoded.Path != "README.md" || decoded.SHA != "abc123" || decoded.Size != 42 {
		t.Errorf("unexpected decoded result: %+v", decoded)
	}
}

func TestPrintEntries_Text(t *testing.T) {
	var buf bytes.Buffer
	entries := []EntryData{
		{Type: "file", Path: "a.md"},
		{Type: "dir", Path: "sub"},
	}
	if err := PrintEntries(&buf, entries, false); err != nil {
		t.Fatal(err)
	}
	want := "file\ta.md\ndir\tsub\n"
	if buf.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestPrintEntries_JSON(t *testing.T) {
	var buf bytes.Buffer
	entries := []EntryData{
		{Type: "file", Path: "a.md", SHA: "aaa", Size: 10},
		{Type: "dir", Path: "sub", SHA: "bbb", Size: 0},
	}
	if err := PrintEntries(&buf, entries, true); err != nil {
		t.Fatal(err)
	}
	var decoded []EntryData
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(decoded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(decoded))
	}
	if decoded[0].Type != "file" || decoded[1].Type != "dir" {
		t.Errorf("unexpected types: %+v", decoded)
	}
}

func TestPrintMutationResult_Text(t *testing.T) {
	var buf bytes.Buffer
	r := MutationResultData{
		Action: "created",
		Path:   "new-file.txt",
		SHA:    "abc123",
		Branch: "main",
	}
	if err := PrintMutationResult(&buf, r, false); err != nil {
		t.Fatal(err)
	}
	want := "action: created\npath: new-file.txt\nsha: abc123\nbranch: main\n"
	if buf.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestPrintMutationResult_TextNoBranch(t *testing.T) {
	var buf bytes.Buffer
	r := MutationResultData{
		Action: "deleted",
		Path:   "old-file.txt",
		SHA:    "def456",
	}
	if err := PrintMutationResult(&buf, r, false); err != nil {
		t.Fatal(err)
	}
	want := "action: deleted\npath: old-file.txt\nsha: def456\n"
	if buf.String() != want {
		t.Errorf("got:\n%s\nwant:\n%s", buf.String(), want)
	}
}

func TestPrintMutationResult_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := MutationResultData{
		Action: "updated",
		Path:   "file.txt",
		SHA:    "sha789",
		Branch: "dev",
	}
	if err := PrintMutationResult(&buf, r, true); err != nil {
		t.Fatal(err)
	}
	var decoded MutationResultData
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.Action != "updated" || decoded.Path != "file.txt" || decoded.SHA != "sha789" || decoded.Branch != "dev" {
		t.Errorf("unexpected decoded result: %+v", decoded)
	}
}
