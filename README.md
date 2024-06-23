# go-gcs-wrapper
go-gcs-wrapper is a library that makes it easy for Go applications to interact with Google Cloud Storage. 
It provides simple functions for tasks like uploading files and generating signed URLs for secure access. 
This wrapper handles the complex details of the Google Cloud Storage API, allowing developers to use storage features 
with less code and effort.

### installation
```shell
go get -d github.com/michaelwp/go-gcs-wrapper
```

### basic of use
- Upload file
```go
package main

import (
	"context"
	"github.com/joho/godotenv"
	gogcswrapper "github.com/michaelwp/go-gcs-wrapper"
	"log"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	projectId := os.Getenv("GOOGLE_APPLICATION_PROJECT_ID")
	bucket := os.Getenv("GOOGLE_APPLICATION_BUCKET")
	object := "go_gcs.png"
	uploadObjPath := "upload_tes"
	localObjPath := "."

	ctx := context.Background()
	gcs := gogcswrapper.NewGCS(ctx, projectId)

	params := &gogcswrapper.UploadParams{
		BucketAndObject: &gogcswrapper.BucketAndObject{
			Bucket: bucket,
			Object: object,
		},
		LocalObjPath:  localObjPath,
		UploadObjPath: uploadObjPath,
	}

	err := gcs.Upload(ctx, params)
	if err != nil {
		log.Fatal("error uploading file", err)
	}
}
```

- Generate signed Url
```go
package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	gogcswrapper "github.com/michaelwp/go-gcs-wrapper"
	"log"
	"os"
	"time"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	projectId := os.Getenv("GOOGLE_APPLICATION_PROJECT_ID")
	bucket := os.Getenv("GOOGLE_APPLICATION_BUCKET")
	object := "go_gcs.png"
	uploadObjPath := "upload_tes"

	ctx := context.Background()
	gcs := gogcswrapper.NewGCS(ctx, projectId)

	params := &gogcswrapper.GenerateSignedURLParams{
		BucketAndObject: &gogcswrapper.BucketAndObject{
			Bucket: bucket,
			Object: object,
		},
		UploadObjPath:  uploadObjPath,
		ExpirationTime: time.Now().Add(time.Minute * 10),
	}

	signedUrl, err := gcs.GenerateSignedURL(params)
	if err != nil {
		log.Fatal("error generating sign url", err)
	}

	fmt.Printf("Signed URL: %s\n", signedUrl)
}
```
