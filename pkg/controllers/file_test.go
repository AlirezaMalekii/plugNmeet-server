package controllers

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/mynaparrot/plugnmeet-protocol/plugnmeet"
	"github.com/mynaparrot/plugnmeet-server/pkg/config"
	"github.com/mynaparrot/plugnmeet-server/pkg/helpers"
	"github.com/sirupsen/logrus"
)

// TestAuthApiFileMessagesArePlainEnglish guards the integrator contract: the Auth API
// file endpoint answers with curated English text, never with the client i18n keys
// (room.errors.*, file.errors.*) that browsers resolve via t().
func TestAuthApiFileMessagesArePlainEnglish(t *testing.T) {
	messages := map[string]string{
		"room id required":      config.MsgRoomIdRequired,
		"room not active":       config.MsgFileRoomNotActive,
		"invalid request":       config.MsgFileInvalidRequest,
		"upload in progress":    config.MsgFileUploadInProgress,
		"invalid document link": config.MsgFileInvalidDocumentLink,
		"too large":             config.MsgFileTooLarge,
		"type not allowed":      config.MsgFileTypeNotAllowed,
		"upload failed":         config.MsgFileUploadFailed,
		"conversion failed":     config.MsgFileConversionFailed,
	}

	for name, msg := range messages {
		if !strings.Contains(msg, " ") {
			t.Errorf("%s: %q is not plain English", name, msg)
		}
		if strings.Contains(msg, "errors.") || strings.HasPrefix(msg, "notifications.") {
			t.Errorf("%s: %q looks like an i18n key", name, msg)
		}
	}

	// the room model's messages are client i18n keys and must never reach integrators
	for _, key := range []error{config.ErrRoomNotActive, config.ErrRoomNotFound, config.ErrRoomUnexpectedError} {
		if config.MsgFileRoomNotActive == key.Error() {
			t.Errorf("room i18n key leaked into the Auth API file message: %q", config.MsgFileRoomNotActive)
		}
	}
}

func TestChatAttachmentPermission(t *testing.T) {
	room := &plugnmeet.RoomMetadata{
		RoomFeatures: &plugnmeet.RoomCreateFeatures{
			ChatFeatures: &plugnmeet.ChatFeatures{IsAllow: true, IsAllowFileUpload: true},
		},
		DefaultLockSettings: &plugnmeet.LockSettings{LockChatFileShare: new(true)},
	}
	attendee := &plugnmeet.NatsKvUserInfo{}
	metadata := &plugnmeet.UserMetadata{}
	if chatFileSharingAllowed(room, attendee, metadata) {
		t.Fatal("a participant must inherit the room's file lock")
	}
	if !chatFileSharingAllowed(room, &plugnmeet.NatsKvUserInfo{IsAdmin: true}, metadata) {
		t.Fatal("a moderator must be able to share when attendees are locked")
	}
	metadata.LockSettings = &plugnmeet.LockSettings{LockChatFileShare: new(false)}
	if !chatFileSharingAllowed(room, attendee, metadata) {
		t.Fatal("an explicit grant must override the room's file lock")
	}
	metadata.LockSettings.LockChatSendMessage = new(true)
	if chatFileSharingAllowed(room, attendee, metadata) {
		t.Fatal("a file grant must not bypass the chat message lock")
	}
}

func TestChatAttachmentDownload(t *testing.T) {
	directory := t.TempDir()
	if err := os.Mkdir(filepath.Join(directory, "session-one"), 0700); err != nil {
		t.Fatal(err)
	}
	name := helpers.IneshastChatFilePrefix + "example.txt"
	if err := os.WriteFile(filepath.Join(directory, "session-one", name), []byte("meeting file"), 0600); err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	controller := &FileController{
		AppConfig: &config.AppConfig{UploadFileSettings: config.UploadFileSettings{Path: directory}},
		logger:    logrus.NewEntry(logrus.New()),
	}
	app.Get("/download/*", controller.HandleDownloadUploadedFile)
	app.Get("/api/download/*", func(c fiber.Ctx) error {
		// The API middleware supplies this value after authorizing the room token.
		c.Locals("roomSid", c.Get("Test-Room"))
		return controller.HandleDownloadUploadedFile(c)
	})

	for _, test := range []struct {
		name   string
		target string
		sid    string
		want   int
	}{
		{"public route", "/download", "", 404},
		{"another room", "/api/download", "session-two", 404},
		{"the same room", "/api/download", "session-one", 200},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", test.target+"/session-one/"+name, nil)
			request.Header.Set("Test-Room", test.sid)
			response, err := app.Test(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != test.want {
				t.Fatalf("status = %d, want %d", response.StatusCode, test.want)
			}
		})
	}
}
