package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	clerrors "githubRAGCli/internal/exitcode"
)

// newTestService creates a RepoService pointing at a test server.
func newTestService(srvURL string) *RepoService {
	return NewRepoService(srvURL, "test-token", 5*time.Second, "owner", "repo")
}

// --- Stat tests ---

func TestStat_File(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/contents/README.md" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"type":         "file",
			"path":         "README.md",
			"sha":          "abc123",
			"size":         42,
			"download_url": "https://raw.githubusercontent.com/owner/repo/main/README.md",
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	entry, err := svc.Stat("", "README.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Type != "file" {
		t.Errorf("type: got %q, want file", entry.Type)
	}
	if entry.Path != "README.md" {
		t.Errorf("path: got %q, want README.md", entry.Path)
	}
	if entry.SHA != "abc123" {
		t.Errorf("sha: got %q, want abc123", entry.SHA)
	}
	if entry.Size != 42 {
		t.Errorf("size: got %d, want 42", entry.Size)
	}
}

func TestStat_Dir(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Contents API returns array for directories.
		json.NewEncoder(w).Encode([]map[string]any{
			{"type": "file", "path": "docs/a.md", "sha": "aaa", "size": 10},
			{"type": "file", "path": "docs/b.md", "sha": "bbb", "size": 20},
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	entry, err := svc.Stat("", "docs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Type != "dir" {
		t.Errorf("type: got %q, want dir", entry.Type)
	}
}

func TestStat_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	_, err := svc.Stat("", "nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitNotFound {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitNotFound)
	}
}

// --- List tests ---

func TestList_Flat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"type": "file", "path": "docs/a.md", "sha": "aaa", "size": 10},
			{"type": "dir", "path": "docs/sub", "sha": "ddd", "size": 0},
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	entries, err := svc.List("", "docs", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Type != "file" || entries[0].Path != "docs/a.md" {
		t.Errorf("entry[0]: got %+v", entries[0])
	}
	if entries[1].Type != "dir" || entries[1].Path != "docs/sub" {
		t.Errorf("entry[1]: got %+v", entries[1])
	}
}

func TestList_Recursive(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch {
		case r.URL.Path == "/repos/owner/repo/contents/docs":
			// First call: return directory listing (to check it's a dir).
			json.NewEncoder(w).Encode([]map[string]any{
				{"type": "file", "path": "docs/a.md", "sha": "aaa", "size": 10},
			})
		case r.URL.Path == "/repos/owner/repo/contents/":
			// getDirSHA for parent of "docs"
			json.NewEncoder(w).Encode([]map[string]any{
				{"type": "dir", "path": "docs", "sha": "treeSHA", "size": 0},
			})
		case r.URL.Path == "/repos/owner/repo/git/trees/treeSHA":
			json.NewEncoder(w).Encode(map[string]any{
				"sha": "treeSHA",
				"tree": []map[string]any{
					{"path": "a.md", "type": "blob", "sha": "aaa", "size": 10},
					{"path": "sub", "type": "tree", "sha": "sss", "size": 0},
					{"path": "sub/b.md", "type": "blob", "sha": "bbb", "size": 20},
				},
				"truncated": false,
			})
		default:
			t.Errorf("unexpected request: %s", r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	entries, err := svc.List("", "docs", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	// Check paths include the base prefix.
	if entries[0].Path != "docs/a.md" {
		t.Errorf("entry[0].Path: got %q, want docs/a.md", entries[0].Path)
	}
}

func TestList_FilePathFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"type": "file",
			"path": "README.md",
			"sha":  "abc",
			"size": 42,
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	_, err := svc.List("", "README.md", false)
	if err == nil {
		t.Fatal("expected error for file path")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitBadArgs {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitBadArgs)
	}
}

// --- ReadFile tests ---

func TestReadFile_Success(t *testing.T) {
	content := "Hello, World!"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"type":     "file",
			"path":     "README.md",
			"sha":      "abc123",
			"size":     len(content),
			"content":  encoded,
			"encoding": "base64",
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	data, err := svc.ReadFile("", "README.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != content {
		t.Errorf("content: got %q, want %q", string(data), content)
	}
}

func TestReadFile_DirFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"type": "file", "path": "docs/a.md", "sha": "aaa", "size": 10},
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	_, err := svc.ReadFile("", "docs")
	if err == nil {
		t.Fatal("expected error for directory")
	}
}

// --- Download tests ---

func TestDownload_SingleFile(t *testing.T) {
	content := "file content"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"type":     "file",
			"path":     "README.md",
			"sha":      "abc123",
			"size":     len(content),
			"content":  encoded,
			"encoding": "base64",
		})
	}))
	defer srv.Close()

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "README.md")

	svc := newTestService(srv.URL)
	if err := svc.Download("", "README.md", outPath, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if string(data) != content {
		t.Errorf("content: got %q, want %q", string(data), content)
	}
}

func TestDownload_OverwriteBlocked(t *testing.T) {
	content := "file content"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"type":     "file",
			"path":     "README.md",
			"sha":      "abc123",
			"size":     len(content),
			"content":  encoded,
			"encoding": "base64",
		})
	}))
	defer srv.Close()

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "README.md")

	// Create existing file.
	if err := os.WriteFile(outPath, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := newTestService(srv.URL)
	err := svc.Download("", "README.md", outPath, false)
	if err == nil {
		t.Fatal("expected error when file exists without --overwrite")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitLocalWriteErr {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitLocalWriteErr)
	}
}

func TestDownload_OverwriteAllowed(t *testing.T) {
	content := "new content"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"type":     "file",
			"path":     "README.md",
			"sha":      "abc123",
			"size":     len(content),
			"content":  encoded,
			"encoding": "base64",
		})
	}))
	defer srv.Close()

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "README.md")

	// Create existing file.
	if err := os.WriteFile(outPath, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := newTestService(srv.URL)
	if err := svc.Download("", "README.md", outPath, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if string(data) != content {
		t.Errorf("content: got %q, want %q", string(data), content)
	}
}

func TestDownload_Directory(t *testing.T) {
	fileContent := "hello"
	encoded := base64.StdEncoding.EncodeToString([]byte(fileContent))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/owner/repo/contents/docs":
			// First call: Download checks if file or dir.
			// Second call: List flat (to verify it's a dir).
			// Third call: List recursive, first checks Contents API.
			json.NewEncoder(w).Encode([]map[string]any{
				{"type": "file", "path": "docs/a.md", "sha": "aaa", "size": 5},
			})
		case "/repos/owner/repo/contents/":
			// getDirSHA: parent listing to find docs SHA.
			json.NewEncoder(w).Encode([]map[string]any{
				{"type": "dir", "path": "docs", "sha": "treeSHA", "size": 0},
			})
		case "/repos/owner/repo/git/trees/treeSHA":
			json.NewEncoder(w).Encode(map[string]any{
				"sha": "treeSHA",
				"tree": []map[string]any{
					{"path": "a.md", "type": "blob", "sha": "aaa", "size": 5},
				},
				"truncated": false,
			})
		case "/repos/owner/repo/contents/docs/a.md":
			json.NewEncoder(w).Encode(map[string]any{
				"type":     "file",
				"path":     "docs/a.md",
				"sha":      "aaa",
				"size":     5,
				"content":  encoded,
				"encoding": "base64",
			})
		default:
			fmt.Printf("unexpected path: %s\n", r.URL.Path)
			w.WriteHeader(404)
			w.Write([]byte(`{"message":"Not Found"}`))
		}
	}))
	defer srv.Close()

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "local-docs")

	svc := newTestService(srv.URL)
	if err := svc.Download("", "docs", outPath, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outPath, "a.md"))
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if string(data) != fileContent {
		t.Errorf("content: got %q, want %q", string(data), fileContent)
	}
}

// --- Ref flag tests ---

func TestStat_WithRef(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("ref"); got != "v1.0" {
			t.Errorf("ref: got %q, want v1.0", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"type": "file",
			"path": "README.md",
			"sha":  "abc123",
			"size": 42,
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	_, err := svc.Stat("v1.0", "README.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- CreateOrUpdateFile tests ---

func TestCreateOrUpdateFile_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/repos/owner/repo/contents/new-file.txt":
			// File does not exist.
			w.WriteHeader(404)
			w.Write([]byte(`{"message":"Not Found"}`))
		case r.Method == "PUT" && r.URL.Path == "/repos/owner/repo/contents/new-file.txt":
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"content": map[string]any{
					"path": "new-file.txt",
					"sha":  "file-sha",
				},
				"commit": map[string]any{
					"sha": "commit-sha-123",
				},
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	result, err := svc.CreateOrUpdateFile("", "new-file.txt", "add file", []byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "created" {
		t.Errorf("action: got %q, want created", result.Action)
	}
	if result.SHA != "commit-sha-123" {
		t.Errorf("sha: got %q, want commit-sha-123", result.SHA)
	}
	if result.Path != "new-file.txt" {
		t.Errorf("path: got %q, want new-file.txt", result.Path)
	}
}

func TestCreateOrUpdateFile_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/repos/owner/repo/contents/existing.txt":
			// File exists.
			json.NewEncoder(w).Encode(map[string]any{
				"type": "file",
				"path": "existing.txt",
				"sha":  "old-sha",
				"size": 5,
			})
		case r.Method == "PUT" && r.URL.Path == "/repos/owner/repo/contents/existing.txt":
			// Verify SHA is included in request body.
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["sha"] != "old-sha" {
				t.Errorf("expected sha old-sha in body, got %v", body["sha"])
			}
			json.NewEncoder(w).Encode(map[string]any{
				"content": map[string]any{
					"path": "existing.txt",
					"sha":  "new-file-sha",
				},
				"commit": map[string]any{
					"sha": "commit-sha-456",
				},
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	result, err := svc.CreateOrUpdateFile("", "existing.txt", "update file", []byte("updated"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "updated" {
		t.Errorf("action: got %q, want updated", result.Action)
	}
	if result.SHA != "commit-sha-456" {
		t.Errorf("sha: got %q, want commit-sha-456", result.SHA)
	}
}

func TestCreateOrUpdateFile_WithBranch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			w.WriteHeader(404)
			w.Write([]byte(`{"message":"Not Found"}`))
		case r.Method == "PUT":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["branch"] != "feature" {
				t.Errorf("expected branch feature, got %v", body["branch"])
			}
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(map[string]any{
				"content": map[string]any{"path": "f.txt", "sha": "s"},
				"commit":  map[string]any{"sha": "c"},
			})
		default:
			w.WriteHeader(500)
		}
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	result, err := svc.CreateOrUpdateFile("feature", "f.txt", "msg", []byte("data"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Branch != "feature" {
		t.Errorf("branch: got %q, want feature", result.Branch)
	}
}

// --- DeleteFile tests ---

func TestDeleteFile_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/repos/owner/repo/contents/old-file.txt":
			json.NewEncoder(w).Encode(map[string]any{
				"type": "file",
				"path": "old-file.txt",
				"sha":  "file-sha",
				"size": 10,
			})
		case r.Method == "DELETE" && r.URL.Path == "/repos/owner/repo/contents/old-file.txt":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			if body["sha"] != "file-sha" {
				t.Errorf("expected sha file-sha, got %v", body["sha"])
			}
			json.NewEncoder(w).Encode(map[string]any{
				"content": nil,
				"commit": map[string]any{
					"sha": "delete-commit-sha",
				},
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	result, err := svc.DeleteFile("", "old-file.txt", "delete file")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Action != "deleted" {
		t.Errorf("action: got %q, want deleted", result.Action)
	}
	if result.SHA != "delete-commit-sha" {
		t.Errorf("sha: got %q, want delete-commit-sha", result.SHA)
	}
}

func TestDeleteFile_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"message":"Not Found"}`))
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	_, err := svc.DeleteFile("", "nonexistent.txt", "delete")
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*clerrors.CLIError)
	if !ok {
		t.Fatalf("expected CLIError, got %T", err)
	}
	if ce.ExitCode() != clerrors.ExitNotFound {
		t.Errorf("exit code: got %d, want %d", ce.ExitCode(), clerrors.ExitNotFound)
	}
}

func TestDeleteFile_Directory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]any{
			{"type": "file", "path": "docs/a.md", "sha": "aaa", "size": 10},
		})
	}))
	defer srv.Close()

	svc := newTestService(srv.URL)
	_, err := svc.DeleteFile("", "docs", "delete dir")
	if err == nil {
		t.Fatal("expected error for directory delete")
	}
}
