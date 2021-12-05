package docx

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type docInUrl struct {
	Url         string
	content     string
	contentText string
}
type docInFile struct {
	Path        string
	content     string
	contentText string
}

func NewDocInUrl(url string) *docInUrl {
	var f docInUrl
	f.Url = url
	return &f
}

func NewDocInFile(filePath string) *docInFile {
	var f docInFile
	f.Path = filePath
	return &f
}

func (doc *docInUrl) Read() error {
	// Connect to url
	resp, err := http.Get(doc.Url)
	if err != nil {
		return fmt.Errorf("not connected, %s", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connected, bad status: %s", resp.Status)
	}

	// Read file
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("file read error: %s", err)
	}
	reader := bytes.NewReader(bodyBytes)

	// Open a zip archive for reading.
	zipReader, err := zip.NewReader(reader, int64(len(bodyBytes)))
	if err != nil {
		return fmt.Errorf("error reading archive : %s", err)
	}

	cont, err := readText(zipReader.File)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	doc.content = cont
	doc.contentText = exactTextDoc(&cont)
	return nil
}

func (doc *docInUrl) GetContent() string {
	return doc.content
}

func (doc *docInUrl) GetContentText() string {
	return doc.contentText
}

func (doc *docInFile) Read() error {
	read, err := zip.OpenReader(doc.Path)
	if err != nil {
		return fmt.Errorf("cannot open file, %s", err)
	}
	content, err := readText(read.File)
	if err != nil {
		return fmt.Errorf("cannot read file, %s", err)
	}
	doc.content = content
	doc.contentText = exactTextDoc(&content)
	return nil
}

func (doc *docInFile) GetContent() string {
	return doc.content
}

func (doc *docInFile) GetContentText() string {
	return doc.contentText
}

func readText(files []*zip.File) (text string, err error) {
	var documentFile *zip.File
	documentFile, err = receiveWordDoc(files)
	if err != nil {
		return text, err
	}
	var documentReader io.ReadCloser
	documentReader, err = documentFile.Open()
	if err != nil {
		return text, err
	}
	text, err = wordDocToString(documentReader)
	if err != nil {
		return text, err
	}
	return
}

// Чтение контента из документа
func wordDocToString(reader io.Reader) (string, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Извлеч весть контент из документа
func receiveWordDoc(files []*zip.File) (file *zip.File, err error) {
	for _, f := range files {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		err = errors.New("document.xml file not found")
	}
	return
}

// Извлечь текст из документа
func exactTextDoc(content *string) string {
	var result []byte
	var data []byte = []byte(*content)
	var writeToResult bool = false
	var isNotText bool = false

	for i := 0; i < len(data); i++ {
		if data[i] == '>' {
			isNotText = false
			continue
		}
		if isNotText {
			continue
		}
		if data[i] == '<' {
			if writeToResult {
				writeToResult = false
				result = append(result, ' ')
			}
			isNotText = true
			continue
		}
		if !writeToResult {
			writeToResult = true
			result = append(result, data[i])
		} else {
			result = append(result, data[i])
		}
	}

	return strings.Trim(string(result), "\n ")
}
