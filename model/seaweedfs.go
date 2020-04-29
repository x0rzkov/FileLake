package model

import (
	fs "github.com/linxGnu/goseaweedfs"
	"io"
	"time"
	"net/http"
)

var (
	masterURL = "GOSWFS_MASTER_URL"
	filers = []string{"GOSWFS_FILER_URL"}

	sw, _ = fs.NewSeaweed(masterURL, filers, 8096, &http.Client{Timeout: 5 * time.Minute})

)

func SaveSeaweedfs(fileReader io.Reader, fileName string, size int64, collection, ttl string) (fp *fs.FilePart, err error) {
	return sw.Upload(fileReader, fileName , size , collection, ttl)
}