package listener

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func processRequest(request *http.Request, response *http.Response) {
	path := path.Clean(request.URL.Path)
	filePath := "./front/build" + path
	info, _ := os.Stat(filePath)

	if info.IsDir() {
		_, err := os.Stat(filePath + "index.html")
		if err == nil {
			file, _ := os.Open(filePath + "index.html")
			response.Body = file
			return
		}
		files, err := readDir(filePath)
		if err != nil {
			response.StatusCode = 500
			response.Body = ioutil.NopCloser(strings.NewReader("Internal server error: " + err.(*os.PathError).Err.Error()))
		}
		filesString := strings.Join(files[:], "\n")
		response.Body = ioutil.NopCloser(strings.NewReader("Index of " + path + ":\n\n" + filesString))
		return
	}
	file, _ := os.Open(filePath)
	response.Body = file
}

func readDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name()+"/")
		} else {
			files = append(files, file.Name())
		}
	}
	return files, nil
}
