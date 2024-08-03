package server

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lostgrace/utility"
	"net/http"
	"os"
	"path"
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
	defer utility.LogIfError(req.Body.Close())
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	bodyBuffer := make([]byte, 0)
	b := bytes.NewBuffer(bodyBuffer)
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return err
	}
	defer utility.LogIfError(resp.Body.Close())
	fmt.Printf("Response: %s\n", b.String())
	return err
}

func DownloadSave(user string, key string) ([]byte, error) {
	sm := &actionMessage{Action: "download", User: user, Key: key}
	resp, err := client.Post(url, "application/json", bytes.NewBufferString(sm.String()))
	var b = make([]byte, 0)
	if err != nil {
		return nil, err
	}
	defer utility.LogIfError(resp.Body.Close())
	buffer := bytes.NewBuffer(b)
	w, err := io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, err
	}

	if w < resp.ContentLength {
		fmt.Printf("Response may not be parsed. Only %d of %d bytes read.\n", w, resp.ContentLength)
	}

	ar := &actionResponse{}
	err = json.Unmarshal(buffer.Bytes(), ar)

	if err != nil {
		str := buffer.String()
		fmt.Println(str)
		return nil, err
	}

	if !ar.Okay {
		return nil, errors.New("could not download file")
	}
	data, err := base64.StdEncoding.DecodeString(ar.Data)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func InstallCoop(outPath string) error {
	//TODO Clean this up a bit. Would probably look better if the individual steps were broken into functions.

	//This request will redirect to the latest release
	req, err := http.NewRequest("GET", coopUrl, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// get the location we were redirected to
	newLocation := resp.Request.URL.Path

	// last part of url should contain the release version
	segments := strings.Split(newLocation, "/")
	version := segments[len(segments)-1]

	// replace the <v> tag in the download string with the actual version
	fullDownload := strings.Replace(download, "<v>", version, 1)

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	// write the download into a buffer
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

	// write the buffer to disk
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

	// copy the files contained in the zip
	for _, zf := range r.File {

		if zf.FileInfo().IsDir() {
			err = os.MkdirAll(path.Join(outPath, zf.Name), 0755)
			if err != nil {
				return err
			}
		}
		nf, err := zf.Open()
		if err != nil {
			return err
		}
		buff := bytes.Buffer{}

		// do not need the read byte count
		_, err = io.Copy(&buff, nf)
		if err != nil {
			return err
		}

		filePath := path.Join(outPath, zf.Name)
		outputPath := path.Dir(filePath)
		err = os.MkdirAll(outputPath, 0755)
		if err != nil {
			return nil
		}
		err = os.WriteFile(filePath, buff.Bytes(), 0755)
		if err != nil {
			return err
		}

	}
	// Close the file so we can clean it up
	err = f.Close()
	if err != nil {
		return err
	}

	//clean up
	err = os.RemoveAll("./release.zip")
	if err != nil {
		return err
	}
	return nil
}
