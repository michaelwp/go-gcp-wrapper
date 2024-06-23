/*
	author: Michael Putong @2024
	This is an example on how to implement go-gcs-wrapper
	to generate a Google Cloud Storage's signedUrl
	visit the code repository in github.com/michaelwp/go-gcs-wrapper
*/
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
