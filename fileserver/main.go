package main

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const MAX_FILE_SIZE = 25 * 1024 * 1024
const UPLOAD_DIR = "./uploads"

type FileUploadEntry struct {
	FileName string
	Size     int64
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	htmlFile, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		fmt.Fprintln(w, "Failed to load templates file ", err)
		return
	}
	htmlFile.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Will handle FileUpload")

	if err := r.ParseMultipartForm(MAX_FILE_SIZE); err != nil {
		fmt.Fprintln(w, "Unable to parse Multipart Form data", err)
		return
	}

	file, header, err := r.FormFile("aFile")

	if err != nil {
		fmt.Fprintln(w, "Failed to read File from Formdata ", err)
		return
	}

	if header.Size > MAX_FILE_SIZE {
		fmt.Fprintf(w, "File Upload size exceeds %f MB \n", MAX_FILE_SIZE/(math.Pow(1024, 2)))
		return
	}

	defer file.Close()

	defer r.Body.Close()

	fmt.Printf("File Name Uploaded is %s of size %d \n", header.Filename, header.Size)

	dst, err := os.Create(fmt.Sprintf("%s/%s", UPLOAD_DIR, header.Filename))

	if err != nil {
		fmt.Fprintln(w, "Failed to Create File in Destination ", err)
		return
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)

	if err != nil {
		fmt.Fprintln(w, "Failed to Copy File ", err)
		return
	}

	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileEntries, err := os.ReadDir(UPLOAD_DIR)

	FileUploadEntries := make([]FileUploadEntry, len(fileEntries))

	if err != nil {
		fmt.Fprintln(w, "Failed to read file contents ", err)
		return
	}

	for i, fileInfoEntry := range fileEntries {
		fileInfo, _ := fileInfoEntry.Info()
		FileUploadEntries[i] = FileUploadEntry{
			FileName: fileInfoEntry.Name(),
			Size:     int64(fileInfo.Size()),
		}
	}

	htmlFile, err := template.ParseFiles("./templates/dirContent.html")
	if err != nil {
		fmt.Fprintln(w, "Failed to load templates file ", err)
		return
	}
	htmlFile.Execute(w, FileUploadEntries)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {

	fileName, ok := strings.CutPrefix(r.RequestURI, "/downloads/")

	fmt.Println("Request to Download ", fileName)

	if !ok {
		fmt.Fprintln(w, "Invalid FileName")
		return
	}

	fileName, _ = url.QueryUnescape(fileName)

	file, err := os.Open(fmt.Sprintf("%s/%s", UPLOAD_DIR, fileName))
	if err != nil {
		fmt.Fprintln(w, "Failed to open file ", err)
		return
	}

	defer file.Close()

	extn := mime.TypeByExtension(filepath.Ext(fileName))
	if err != nil {
		fmt.Fprintln(w, "Failed to deduce MimeType of file ", err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", extn)

	io.Copy(w, file)

}

func main() {
	fmt.Println("Initializing Go FileServer")
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/downloads/", downloadHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8200", nil)
}
