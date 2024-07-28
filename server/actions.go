package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	path2 "path"
	"strings"
)

var client = &http.Client{}

const (
	url      = "https://skulle.xyz/actions.php"
	coopUrl  = "https://github.com/LukeYui/EldenRingSeamlessCoopRelease/releases/latest"
	download = "https://github.com/LukeYui/EldenRingSeamlessCoopRelease/releases/download/<v>/ersc.zip"
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

func InstallCoop(outPath string) error {

	req, err := http.NewRequest("GET", coopUrl, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	newLocation := resp.Request.URL.Path
	tags := strings.Split(newLocation, "/")
	version := tags[len(tags)-1]
	fullDownload := strings.Replace(download, "<v>", version, 1)
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	req, err = http.NewRequest("GET", fullDownload, nil)
	if err != nil {
		return nil
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil
	}
	fileBuffer := bytes.Buffer{}
	_, err = io.Copy(&fileBuffer, resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile("./release.zip", fileBuffer.Bytes(), 0664)
	if err != nil {
		return err
	}
	f, err := os.Open("./release.zip")
	if err != nil {
		return err
	}
	r, err := zip.NewReader(f, int64(fileBuffer.Len()))

	if err != nil {
		return err
	}

	for _, zf := range r.File {
		if zf.FileInfo().IsDir() {
			err = os.MkdirAll(path2.Join(outPath, zf.Name), 0755)
			if err != nil {
				return err
			}
		}
		nf, err := zf.Open()
		if err != nil {
			return err
		}
		buff := bytes.Buffer{}
		_, err = io.Copy(&buff, nf)
		if err != nil {
			return err
		}

		path := path2.Dir(path2.Join(outPath, zf.Name))
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return nil
		}
		err = os.WriteFile(path2.Join(outPath, zf.Name), buff.Bytes(), 764)
		if err != nil {
			return err
		}

	}
	err = os.RemoveAll("./release.zip")
	if err != nil {
		return err
	}
	return nil
}
