package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func mediaUpload(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		renderError(w, "INVALID_FILE", http.StatusUnprocessableEntity)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// get current path
	dir, _ := os.Getwd()

	fileType := http.DetectContentType(fileBytes)
	fmt.Println(fileType)
	if fileType != "audio/mp3" && fileType != "audio/mpeg" && fileType != "application/octet-stream" {
		renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return
	}

	_, typeErr := mime.ExtensionsByType(fileType)
	if typeErr != nil {
		renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return
	}

	fileName := randToken()
	newPath := filepath.Join(fmt.Sprintf("%s/%s", dir, "media/uploads"), fmt.Sprintf("%s.mp3", fileName))
	fmt.Println(newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}

	if _, err := newFile.Write(fileBytes); err != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}

	newFile.Close()

	outputPath := "/media/output"
	splitErr := split(newPath, fmt.Sprintf("%s/%s", dir, outputPath))
	if splitErr != nil {
		fmt.Println(splitErr)
		renderError(w, "PROCESSING_ERROR", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	content := map[string]string{
		"vocals":        fmt.Sprintf("%s/%s/%s", outputPath, fileName, "vocals.wav"),
		"accompaniment": fmt.Sprintf("%s/%s/%s", outputPath, fileName, "accompaniment.wav"),
	}
	jsonBody, _ := json.Marshal(content)
	fmt.Fprint(w, string(jsonBody))
}

func randToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func renderError(w http.ResponseWriter, message string, status int) {
	fmt.Println(message)
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, message)
	}
}
