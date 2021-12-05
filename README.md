# read-docx
A simple library to read text from a docx file
[![Go Report Card](https://goreportcard.com/badge/github.com/opencontrol/doc-template)](https://goreportcard.com/report/github.com/opencontrol/doc-template)

```go
url := "https://example/com/document.docx"
doc1 := word.NewDocInUrl(url)
err1 := doc1.Read()
if err1 != nil {
	log.Fatal(err1)
}
fmt.Println(doc1.GetContent())

filePath := "test.docx"
doc2 := word.NewDocInFile(filePath)
err2 := doc2.Read()
if err2 != nil {
	log.Fatal(err2)
}
fmt.Println(doc2.GetContentText())
```
