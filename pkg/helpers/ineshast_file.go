package helpers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const IneshastChatFilePrefix = "ineshast-chat-"

// ArchiveIneshastFile sends an authenticated copy to the configured private archive.
func ArchiveIneshastFile(ctx context.Context, secret, meeting, sid, id, name, source, path string) error {
	endpoint, err := url.Parse(os.Getenv("INESHAST_SHARED_FILES_URL"))
	if err != nil || endpoint.Scheme != "https" || endpoint.Host == "" || secret == "" {
		return fmt.Errorf("iNeshast shared file archive is not configured")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() < 1 || info.Size() > 20*1024*1024 {
		return fmt.Errorf("iNeshast shared files must be between 1 byte and 20 MB")
	}
	digest := sha256.New()
	if _, err := file.WriteTo(digest); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	metadata, err := json.Marshal(map[string]any{
		"meeting": meeting, "sid": sid, "id": id, "name": name,
		"source": source, "sha256": hex.EncodeToString(digest.Sum(nil)), "time": time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	header := base64.StdEncoding.EncodeToString(metadata)
	signature := hmac.New(sha256.New, []byte(secret))
	signature.Write([]byte(header))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), file)
	if err != nil {
		return err
	}
	req.ContentLength = info.Size()
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Room-File", header)
	req.Header.Set("X-Room-File-Signature", hex.EncodeToString(signature.Sum(nil)))
	client := &http.Client{Timeout: 35 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("iNeshast archive returned HTTP %d", res.StatusCode)
	}

	return nil
}
