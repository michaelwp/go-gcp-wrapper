/*
	author: Michael Putong @2024
	-----------------------------------------------------------------
	This is a library that is wrapping the Google Cloud Storage Library
	to simplify the interaction within Go application and Google Cloud Storage.
	-----------------------------------------------------------------
	This code is free to use, modify and distribute, although
	the author is not responsible for any damage occurred in its use.
	-----------------------------------------------------------------
	visit the code repository in github.com/michaelwp/go-semaphore
*/

package go_gcs_wrapper

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Gcs interface {
	Upload(ctx context.Context, params *UploadParams) error
	GenerateSignedURL(params *GenerateSignedURLParams) (string, error)
}

type gcs struct {
	StorageClient *storage.Client
	ProjectId     string
}

type BucketAndObject struct {
	Bucket string
	Object string
}

type UploadParams struct {
	LocalObjPath  string
	UploadObjPath string
	*BucketAndObject
}

type GenerateSignedURLParams struct {
	*BucketAndObject
	ExpirationTime time.Time
	UploadObjPath  string
}

func NewGCS(ctx context.Context, projectId string) Gcs {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("failed to create GCS client")
	}
	return &gcs{
		StorageClient: client,
		ProjectId:     projectId,
	}
}

func (g gcs) Upload(ctx context.Context, params *UploadParams) error {
	// open the local file that is intended to be uploaded to GCS.
	// ensure the file is closed at the end.
	file, err := os.Open(params.LocalObjPath + "/" + params.Object)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(file)

	// set the timeout to 1 minute
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	object := g.StorageClient.
		Bucket(params.Bucket).
		Object(params.UploadObjPath + "/" + params.Object)

	// set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	object = object.If(storage.Conditions{DoesNotExist: true})

	// Upload an object with storage.Writer, and close it at the end.
	writer := object.NewWriter(ctx)
	defer func(w *storage.Writer) {
		err := w.Close()
		if err != nil {
			log.Printf("failed to close writer: %v", err)
		}
	}(writer)

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	log.Printf("File %s successfully uploaded to %s/%s .\n",
		params.Object, params.Bucket, params.UploadObjPath)

	return nil
}

func (g gcs) GenerateSignedURL(params *GenerateSignedURLParams) (string, error) {
	// Set up the signed URL options.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: params.ExpirationTime,
	}

	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	signedUrl, err := g.StorageClient.
		Bucket(params.Bucket).
		SignedURL(params.UploadObjPath+"/"+params.Object, opts)
	if err != nil {
		log.Fatalf("Failed to generate signed URL: %v", err)
	}

	log.Println("SignedURL generated successfully")
	return signedUrl, nil
}
