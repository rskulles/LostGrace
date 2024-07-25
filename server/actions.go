package server

import (
	"bytes"
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
	Action   string `json:"action,omitempty"`
	User     string `json:"user,omitempty"`
	Key      string `json:"key,omitempty"`
	FileName string `json:"file_name,omitempty"`
	File     []byte `json:"file,omitempty"`
}
type actionResponse struct {
	Okay bool   `json:"okay"`
	Data string `json:"data,omitempty"`
}

func (sm *actionMessage) String() string {
	/*
		pr, pw := io.Pipe()
		go func() {

			gw, err := gzip.NewWriterLevel(pw, gzip.BestCompression)

			if err != nil {
				panic(err)
			}
			err = json.NewEncoder(gw).Encode(sm)
			defer gw.Close()
			defer pw.CloseWithError(err)
		}()
		bb := make([]byte, 0)
		b := bytes.NewBuffer(bb)
		c, err := io.Copy(b, pr)
		fmt.Println(c)
		if err != nil {
			panic(err)
		}
				return pr
	*/
	b, _ := json.Marshal(sm)
	return string(b)
}
func UploadSave(user string, key string, filename string, save []byte) error {
	sm := &actionMessage{Action: "upload", User: user, Key: key, FileName: filename, File: save}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(sm.String()))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")
	//	req.Header.Set("Content-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyBuffer := make([]byte, 0)
	b := bytes.NewBuffer(bodyBuffer)
	_, _ = io.Copy(b, resp.Body)
	defer resp.Body.Close()
	fmt.Printf("Response: %s\n", b.String())
	return err
}

func DownloadSave(user string, key string) (string, error) {
	sm := &actionMessage{Action: "download", User: user, Key: key}
	resp, err := client.Post(url, "application/json", bytes.NewBufferString(sm.String()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b := make([]byte, 0)
	buffer := bytes.NewBuffer(b)
	w, err := io.Copy(buffer, resp.Body)
	if err != nil {
		return "", err
	}

	if w < resp.ContentLength {
		fmt.Printf("Response may not be parsed. Only %d of %d bytes read.\n", w, resp.ContentLength)
	}

	ar := &actionResponse{}
	err = json.Unmarshal(buffer.Bytes(), ar)
	if err != nil {
		str := buffer.String()
		fmt.Println(str)
		return "", err
	}
	if !ar.Okay {
		return "", errors.New("could not download file")
	}
	return ar.Data, nil
}
