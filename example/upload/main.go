/*
	author: Michael Putong @2024
	This is an example on how to implement go-gcs-wrapper
	to upload a file to Google Cloud Storage
	visit the code repository in github.com/michaelwp/go-gcs-wrapper
*/

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
