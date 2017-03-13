package util

import (
	"os"
	"bytes"
	"io"
	"strings"
)

func Parse(s string) ( result string) {
	f, _ := os.Open(s)
	defer f.Close()
	result = getContent(f)
	return
}

func getContent(file *os.File) (finalContent string) {
	var result *bytes.Buffer = new(bytes.Buffer)
	content := make([]byte, 10)
	for {
		l,err := file.Read(content)
		if l == 0 || err == io.EOF {
			break
		}
		if result == nil {
			result = bytes.NewBuffer(content[:l])
		} else {
			result.Write(content[:l])
		}
	}
	
	initialContent := result.String()
	splittedContent := strings.Split(initialContent,"\n")
	tempContent := make([]string,1)
	
	for i,values := range splittedContent {
		if values != "" {
			if i==0 {
				tempContent[i] = values
			} else {
				test := append(tempContent,values)
				tempContent = test
			}
		}
	}
	
	finalContent = strings.Join(tempContent,"\n")
	
	return
}