package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	clerrors "githubRAGCli/internal/exitcode"
	"githubRAGCli/internal/githubapi"
)

// Entry is the unified model for repository content metadata.
type Entry struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"download_url,omitempty"`
}

// contentsItem maps the JSON returned by the GitHub Contents API.
type contentsItem struct {
	Type        string `json:"type"` // "file", "dir", "symlink", "submodule"
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"download_url"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

// RepoService provides business logic for repository content operations.
type RepoService struct {
	Client *githubapi.Client
	Owner  string
	Repo   string
}

// NewRepoService creates a RepoService for the given owner/repo.
func NewRepoService(apiBase, token string, timeout time.Duration, owner, repo string) *RepoService {
	return &RepoService{
		Client: githubapi.NewClient(apiBase, token, timeout),
		Owner:  owner,
		Repo:   repo,
	}
}

// Stat returns metadata for a single path.
func (s *RepoService) Stat(ref, path string) (*Entry, error) {
	raw, err := s.Client.GetContents(s.Owner, s.Repo, path, ref)
	if err != nil {
		return nil, err
	}

	// Contents API returns an object for a file and an array for a directory.
	// For stat, we need the single object. If it's an array, the path is a directory
	// and we use the first-level info from the API.
	var item contentsItem
	if err := json.Unmarshal(raw, &item); err == nil && item.SHA != "" {
		return itemToEntry(&item), nil
	}

	// It might be a directory — API returns an array.
	// For stat on a directory, we need to call the API with trailing slash behavior.
	// The Contents API on a dir path returns the children; we derive dir metadata.
	var items []contentsItem
	if err := json.Unmarshal(raw, &items); err == nil {
		return &Entry{
			Type: "dir",
			Path: path,
		}, nil
	}

	return nil, clerrors.NewTransport("unexpected contents response format", nil)
}

// List returns directory entries. Non-recursive uses the Contents API;
// recursive uses the Trees API.
func (s *RepoService) List(ref, path string, recursive bool) ([]Entry, error) {
	if !recursive {
		return s.listFlat(ref, path)
	}
	return s.listRecursive(ref, path)
}

func (s *RepoService) listFlat(ref, path string) ([]Entry, error) {
	raw, err := s.Client.GetContents(s.Owner, s.Repo, path, ref)
	if err != nil {
		return nil, err
	}

	// If the path is a file, the API returns an object, not an array.
	var single contentsItem
	if err := json.Unmarshal(raw, &single); err == nil && single.SHA != "" && single.Type != "dir" {
		return nil, clerrors.NewBadArgs(fmt.Sprintf("path %q is a file, not a directory", path), nil)
	}

	var items []contentsItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, clerrors.NewBadArgs(fmt.Sprintf("path %q is a file, not a directory", path), nil)
	}

	entries := make([]Entry, 0, len(items))
	for i := range items {
		entries = append(entries, *itemToEntry(&items[i]))
	}
	return entries, nil
}

func (s *RepoService) listRecursive(ref, path string) ([]Entry, error) {
	// First, get the directory SHA via Contents API.
	raw, err := s.Client.GetContents(s.Owner, s.Repo, path, ref)
	if err != nil {
		return nil, err
	}

	// Check if it's a single file (not listable).
	var single contentsItem
	if err := json.Unmarshal(raw, &single); err == nil && single.SHA != "" && single.Type != "dir" {
		return nil, clerrors.NewBadArgs(fmt.Sprintf("path %q is a file, not a directory", path), nil)
	}

	// For a directory, we need its tree SHA. The Contents API for a dir returns children
	// but doesn't directly give us the tree SHA. We need to get it from the parent
	// or use the path itself. Let's get the dir's SHA from stat-level info.
	// We'll call GetContents on the parent to find this dir's SHA, or if path is root, use ref.
	dirSHA, err := s.getDirSHA(ref, path)
	if err != nil {
		return nil, err
	}

	tree, err := s.Client.GetTree(s.Owner, s.Repo, dirSHA, true)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(tree.Tree))
	for _, te := range tree.Tree {
		e := Entry{
			Type: treeTypeToEntryType(te.Type),
			Path: joinPath(path, te.Path),
			SHA:  te.SHA,
			Size: te.Size,
		}
		entries = append(entries, e)
	}

	if tree.Truncated {
		fmt.Fprintf(os.Stderr, "warning: tree listing was truncated by GitHub API\n")
	}

	return entries, nil
}

// getDirSHA resolves the git tree SHA for a directory path.
func (s *RepoService) getDirSHA(ref, path string) (string, error) {
	if path == "" || path == "." || path == "/" {
		// Root tree — use ref directly (branch/tag/sha).
		effectiveRef := ref
		if effectiveRef == "" {
			effectiveRef = "HEAD"
		}
		return effectiveRef, nil
	}

	// Get parent directory listing to find this dir's SHA.
	parent := filepath.Dir(path)
	if parent == "." {
		parent = ""
	}
	base := filepath.Base(path)

	raw, err := s.Client.GetContents(s.Owner, s.Repo, parent, ref)
	if err != nil {
		return "", err
	}

	var items []contentsItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return "", clerrors.NewTransport("unexpected contents response format", err)
	}

	for _, item := range items {
		if filepath.Base(item.Path) == base && item.Type == "dir" {
			return item.SHA, nil
		}
	}

	return "", clerrors.NewNotFound(fmt.Sprintf("directory %q not found", path), nil)
}

// ReadFile returns the decoded content of a file.
func (s *RepoService) ReadFile(ref, path string) ([]byte, error) {
	raw, err := s.Client.GetContents(s.Owner, s.Repo, path, ref)
	if err != nil {
		return nil, err
	}

	var item contentsItem
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, clerrors.NewTransport("unexpected contents response format", err)
	}

	if item.Type != "file" {
		return nil, clerrors.NewBadArgs(fmt.Sprintf("path %q is a %s, not a file", path, item.Type), nil)
	}

	if item.Encoding != "base64" {
		return nil, clerrors.NewTransport(fmt.Sprintf("unsupported encoding %q", item.Encoding), nil)
	}

	// GitHub base64 content may contain newlines; strip them.
	cleaned := strings.ReplaceAll(item.Content, "\n", "")
	data, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return nil, clerrors.NewTransport("failed to decode base64 content", err)
	}

	return data, nil
}

// Download writes repository content to the local filesystem.
// For a file, it writes the decoded content to outPath.
// For a directory, it recursively lists and downloads all files.
func (s *RepoService) Download(ref, remotePath, outPath string, overwrite bool) error {
	// Determine if the path is a file or directory.
	raw, err := s.Client.GetContents(s.Owner, s.Repo, remotePath, ref)
	if err != nil {
		return err
	}

	var item contentsItem
	if err := json.Unmarshal(raw, &item); err == nil && item.SHA != "" && item.Type == "file" {
		return s.downloadFile(ref, remotePath, outPath, overwrite)
	}

	// Directory download.
	return s.downloadDir(ref, remotePath, outPath, overwrite)
}

func (s *RepoService) downloadFile(ref, remotePath, outPath string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(outPath); err == nil {
			return clerrors.NewLocalWriteErr(fmt.Sprintf("file already exists: %s (use --overwrite to replace)", outPath), nil)
		}
	}

	data, err := s.ReadFile(ref, remotePath)
	if err != nil {
		return err
	}

	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return clerrors.NewLocalWriteErr("failed to create directory: "+dir, err)
	}

	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return clerrors.NewLocalWriteErr("failed to write file: "+outPath, err)
	}

	return nil
}

func (s *RepoService) downloadDir(ref, remotePath, outPath string, overwrite bool) error {
	entries, err := s.List(ref, remotePath, true)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Type != "file" {
			continue
		}

		// Compute relative path within the downloaded directory.
		relPath := entry.Path
		if remotePath != "" && remotePath != "." && remotePath != "/" {
			relPath = strings.TrimPrefix(entry.Path, remotePath+"/")
		}

		localPath := filepath.Join(outPath, relPath)
		if err := s.downloadFile(ref, entry.Path, localPath, overwrite); err != nil {
			return err
		}
	}

	return nil
}

// DownloadFileFromURL downloads a file from a raw URL (e.g., download_url).
func (s *RepoService) DownloadFileFromURL(url, outPath string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(outPath); err == nil {
			return clerrors.NewLocalWriteErr(fmt.Sprintf("file already exists: %s (use --overwrite to replace)", outPath), nil)
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return clerrors.NewTransport("failed to download file", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return clerrors.NewTransport(fmt.Sprintf("download failed with HTTP %d", resp.StatusCode), nil)
	}

	dir := filepath.Dir(outPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return clerrors.NewLocalWriteErr("failed to create directory: "+dir, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return clerrors.NewLocalWriteErr("failed to create file: "+outPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return clerrors.NewLocalWriteErr("failed to write file: "+outPath, err)
	}

	return nil
}

func itemToEntry(item *contentsItem) *Entry {
	return &Entry{
		Type:        item.Type,
		Path:        item.Path,
		SHA:         item.SHA,
		Size:        item.Size,
		DownloadURL: item.DownloadURL,
	}
}

func treeTypeToEntryType(t string) string {
	switch t {
	case "blob":
		return "file"
	case "tree":
		return "dir"
	default:
		return t
	}
}

func joinPath(base, rel string) string {
	if base == "" || base == "." || base == "/" {
		return rel
	}
	return base + "/" + rel
}
