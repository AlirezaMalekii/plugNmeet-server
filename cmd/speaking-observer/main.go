// speaking-observer records server-observed speech transitions independently of browser telemetry.
package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type journal struct {
	mu      sync.Mutex
	file    *os.File
	room    string
	session string
}

func (j *journal) write(kind, identity, name string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	event := map[string]any{"event": kind, "room_id": j.room, "room_sid": j.session, "participant_key": identity, "name": name, "timestamp_ms": time.Now().UnixMilli()}
	if err := json.NewEncoder(j.file).Encode(event); err != nil {
		log.Printf("journal encode: %v", err)
		return
	}
	if err := j.file.Sync(); err != nil {
		log.Printf("journal sync: %v", err)
	}
}
func observe(ctx context.Context, info *livekit.Room, output, host, key, secret string) error {
	// A SID-scoped filename separates repeated sessions using the same room name.
	f, err := os.OpenFile(filepath.Join(output, info.Sid+"-"+uuid.NewString()+".jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	j := journal{file: f, room: info.Name, session: info.Sid}
	done := make(chan struct{})
	var once sync.Once
	room := lksdk.NewRoom(&lksdk.RoomCallback{
		ParticipantCallback: lksdk.ParticipantCallback{
			OnIsSpeakingChanged: func(p lksdk.Participant) {
				kind := "speech_ended"
				if p.IsSpeaking() {
					kind = "speech_started"
				}
				j.write(kind, p.Identity(), p.Name())
			},
		},
		OnParticipantDisconnected: func(p *lksdk.RemoteParticipant) { j.write("participant_left", p.Identity(), p.Name()) },
		OnDisconnected:            func() { once.Do(func() { close(done) }) },
	})
	no := false
	token, err := auth.NewAccessToken(key, secret).SetIdentity("ineshast-observer-" + info.Sid).
		SetVideoGrant(&auth.VideoGrant{RoomJoin: true, Room: info.Name, Hidden: true, CanPublish: &no, CanSubscribe: &no}).
		SetValidFor(10 * time.Minute).ToJWT()
	if err != nil {
		return err
	}
	if err = room.JoinWithToken(host, token, lksdk.WithAutoSubscribe(false)); err != nil {
		return err
	}
	j.write("observer_connected", "", "")
	select {
	case <-ctx.Done():
	case <-done:
	}
	room.Disconnect()
	j.write("observer_disconnected", "", "")
	if err := os.WriteFile(f.Name()+".ready", []byte("ready"), 0600); err != nil {
		return err
	}
	return nil
}
func main() {
	host, key, secret := os.Getenv("LIVEKIT_URL"), os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET")
	if host == "" || key == "" || secret == "" {
		log.Fatal("LIVEKIT_URL, LIVEKIT_API_KEY and LIVEKIT_API_SECRET required")
	}
	output := os.Getenv("SPEAKING_OUTPUT")
	if output == "" {
		output = "/data/speaking"
	}
	if err := os.MkdirAll(output, 0700); err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := recoverJournals(output); err != nil {
		log.Fatal(err)
	}
	go deliver(ctx, output)
	client := lksdk.NewRoomServiceClient(host, key, secret)
	var mu sync.Mutex
	active := map[string]bool{}
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		rooms, err := client.ListRooms(ctx, &livekit.ListRoomsRequest{})
		if err != nil {
			log.Printf("list rooms: %v", err)
			continue
		}
		for _, info := range rooms.Rooms {
			mu.Lock()
			running := active[info.Sid]
			if !running {
				active[info.Sid] = true
			}
			mu.Unlock()
			if running {
				continue
			}
			go func(info *livekit.Room) {
				defer func() { mu.Lock(); delete(active, info.Sid); mu.Unlock() }()
				if err := observe(ctx, info, output, host, key, secret); err != nil {
					log.Print(fmt.Errorf("observe %s: %w", info.Sid, err))
				}
			}(info)
		}
	}
}

func deliver(ctx context.Context, output string) {
	url, key, secret := os.Getenv("SPEAKING_WEBHOOK_URL"), os.Getenv("PLUGNMEET_API_KEY"), os.Getenv("PLUGNMEET_API_SECRET")
	if url == "" || key == "" || secret == "" {
		log.Print("speaking webhook delivery disabled: credentials or URL missing")
		return
	}
	client := &http.Client{Timeout: 20 * time.Second}
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		paths, err := filepath.Glob(filepath.Join(output, "*.jsonl.ready"))
		if err != nil {
			log.Print(err)
			continue
		}
		for _, path := range paths {
			if err := sendJournal(ctx, client, strings.TrimSuffix(path, ".ready"), url, key, secret); err != nil {
				log.Printf("speaking delivery retry: %v", err)
				continue
			}
			if err := os.Rename(path, strings.TrimSuffix(path, ".ready")+".sent"); err != nil {
				log.Print(err)
			}
		}
	}
}
func sendJournal(ctx context.Context, client *http.Client, path, url, key, secret string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	rows := []map[string]any{}
	scan := bufio.NewScanner(f)
	scan.Buffer(make([]byte, 4096), 1024*1024)
	for scan.Scan() {
		var row map[string]any
		if err := json.Unmarshal(scan.Bytes(), &row); err != nil {
			return err
		}
		rows = append(rows, row)
	}
	if err := scan.Err(); err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("empty speaking journal")
	}
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(filepath.Base(path))).String()
	body, err := json.Marshal(map[string]any{"id": id, "event": "ineshast_speaking_events", "createdAt": time.Now().Unix(), "room": map[string]any{"room_id": rows[0]["room_id"], "sid": rows[0]["room_sid"]}, "speaking_events": rows})
	if err != nil {
		return err
	}
	sum := sha256.Sum256(body)
	enc := base64.RawURLEncoding.EncodeToString
	header := enc([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims, _ := json.Marshal(map[string]any{"iss": key, "exp": time.Now().Add(5 * time.Minute).Unix(), "sha256": base64.StdEncoding.EncodeToString(sum[:])})
	token := header + "." + enc(claims)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(token))
	token += "." + enc(mac.Sum(nil))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook HTTP %d", resp.StatusCode)
	}
	return nil
}

// recoverJournals makes interrupted journals deliverable without extending speech
// across an unobserved outage. The gap remains explicit in the payload.
func recoverJournals(output string) error {
	paths, err := filepath.Glob(filepath.Join(output, "*.jsonl"))
	if err != nil {
		return err
	}
	for _, path := range paths {
		if _, err := os.Stat(path + ".sent"); err == nil {
			continue
		}
		if _, err := os.Stat(path + ".ready"); err == nil {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		var last map[string]any
		scan := bufio.NewScanner(f)
		for scan.Scan() {
			var row map[string]any
			if err := json.Unmarshal(scan.Bytes(), &row); err != nil {
				f.Close()
				return fmt.Errorf("invalid journal %s: %w", filepath.Base(path), err)
			}
			last = row
		}
		scanErr := scan.Err()
		f.Close()
		if scanErr != nil {
			return scanErr
		}
		if last == nil {
			continue
		}
		last["event"] = "observer_disconnected"
		last["participant_key"] = ""
		last["name"] = ""
		last["observation_gap"] = true
		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		err = json.NewEncoder(f).Encode(last)
		if err == nil {
			err = f.Sync()
		}
		f.Close()
		if err != nil {
			return err
		}
		if err = os.WriteFile(path+".ready", []byte("recovered"), 0600); err != nil {
			return err
		}
	}
	return nil
}
