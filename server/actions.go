package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var client = &http.Client{}

const (
	url = "https://skulle.xyz/actions.php"
)

type actionMessage struct {
	Action   string
	User     string
	Key      string
	FileName string
	File     []byte
}
type actionResponse struct {
	Okay bool   `json:"okay"`
	Data string `json:"data,omitempty"`
}

func (sm *actionMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Action   string `json:"action,omitempty"`
		User     string `json:"user,omitempty"`
		Key      string `json:"key,omitempty"`
		FileName string `json:"file_name,omitempty"`
		File     string `json:"file,omitempty"`
	}{
		Action:   sm.Action,
		Key:      sm.Key,
		User:     sm.User,
		FileName: sm.FileName,
		File:     base64.StdEncoding.EncodeToString(sm.File),
	})
}

func (sm *actionMessage) String() string {
	j, err := json.Marshal(sm)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("x=%s", j)
}

func UploadSave(user string, key string, filename string, save []byte) error {
	sm := &actionMessage{Action: "upload", User: user, Key: key, FileName: filename, File: save}
	j := sm.String()
	resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(j)))
	if err != nil {
		return err
	}
	bodyBuffer := make([]byte, resp.ContentLength)
	defer resp.Body.Close()
	_, _ = resp.Body.Read(bodyBuffer)
	fmt.Printf("Response: %s\n", string(bodyBuffer))
	return err
}

func DownloadSave(user string, key string) (string, error) {
	sm := &actionMessage{Action: "download", User: user, Key: key}
	j := sm.String()
	resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(j)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBuffer := make([]byte, resp.ContentLength)
	var c int64 = 0
	for {
		i, err := resp.Body.Read(bodyBuffer[c:])
		c += int64(i)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		if c >= resp.ContentLength {
			break
		}

	}

	ar := &actionResponse{}
	err = json.Unmarshal(bodyBuffer, ar)
	if err != nil {
		return "", err
	}
	if !ar.Okay {
		return "", errors.New("could not download file")
	}
	return ar.Data, nil
}
