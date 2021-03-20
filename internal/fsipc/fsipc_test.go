package fsipc

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestFsipc(t *testing.T) {
	dir := t.TempDir()

	ipc, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("PingRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "bda6a8a9d7924c149697e13b93aa68bf.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "ping"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		pingRequest, ok := request.(*PingRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := PingRequest{
			Id: "bda6a8a9d7924c149697e13b93aa68bf",
		}
		if !reflect.DeepEqual(*pingRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GsNewSessionRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "b787ae38ea0c465e8d853015db940915.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "groovestats/new-session"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		newSessionRequest, ok := request.(*GsNewSessionRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := GsNewSessionRequest{
			Id: "b787ae38ea0c465e8d853015db940915",
		}
		if !reflect.DeepEqual(*newSessionRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("GsPlayerScoresRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "8eef8847d10041e2b519ac28165e9a24.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "groovestats/player-scores",
				"chart": "H",
				"api-key-player-2": "K"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		playerScoresRequest, ok := request.(*GsPlayerScoresRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		key := "K"
		expected := GsPlayerScoresRequest{
			Id:            "8eef8847d10041e2b519ac28165e9a24",
			Chart:         "H",
			ApiKeyPlayer1: nil,
			ApiKeyPlayer2: &key,
		}
		if !reflect.DeepEqual(*playerScoresRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("SubmitScoreRequest", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "25a1506cdeff4d01b50f8207313f5db1.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "groovestats/submit-score",
				"api-key": "K",
				"profile-name": "N",
				"hash": "H",
				"score": 10000,
				"rate": 100
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		submitScoreRequest, ok := request.(*SubmitScoreRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := SubmitScoreRequest{
			Id:          "25a1506cdeff4d01b50f8207313f5db1",
			ApiKey:      "K",
			ProfileName: "N",
			Hash:        "H",
			Score:       10000,
			Rate:        100,
		}
		if !reflect.DeepEqual(*submitScoreRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("PartialWrite", func(t *testing.T) {
		go func() {
			filename := filepath.Join(dir, "requests", "cb95f27932174600bafab86e2e5204c7.json")
			err := os.WriteFile(filename, []byte(`{
				"action": "ping"
			}`), 0700)
			if err != nil {
				t.Fatal(err)
			}

			f, err := os.Create(filename)
			if err != nil {
				t.Fatal(err)
			}

			_, err = f.Write([]byte(`{"action": `))
			if err != nil {
				t.Fatal(err)
			}

			err = f.Sync()
			if err != nil {
				t.Fatal(err)
			}

			// give fsipc some time to process the write event
			// before issueing the next one
			time.Sleep(time.Second)

			_, err = f.Write([]byte(`"ping"}`))
			if err != nil {
				t.Fatal(err)
			}

			err = f.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		request, ok := <-ipc.Requests
		if !ok {
			t.Fatal("requests channel closed prematurely")
		}

		pingRequest, ok := request.(*PingRequest)
		if !ok {
			t.Fatal("incorrect request type")
		}

		expected := PingRequest{
			Id: "cb95f27932174600bafab86e2e5204c7",
		}
		if !reflect.DeepEqual(*pingRequest, expected) {
			t.Fatal("unexpected request")
		}
	})

	t.Run("WriteResponse", func(t *testing.T) {
		var data struct {
			Payload string `json:"payload"`
		}
		data.Payload = "foobar"

		err := ipc.WriteResponse("848431421ad846ebbb26269c749a5a43", data)
		if err != nil {
			t.Fatal(err)
		}

		filename := filepath.Join(dir, "responses", "848431421ad846ebbb26269c749a5a43.json")
		serialized, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		expected := []byte(`{"payload":"foobar"}`)
		if !bytes.Equal(serialized, expected) {
			t.Fatal("unexpected data")
		}
	})

	t.Run("Close", func(t *testing.T) {
		err := ipc.Close()
		if err != nil {
			t.Fatal(err)
		}

		_, ok := <-ipc.Requests
		if ok {
			t.Fatal("requests channel not closed")
		}
	})
}
