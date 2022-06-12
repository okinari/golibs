package golibs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

type GoogleCloudStorage struct {
	projectID string
	client    *storage.Client
}

func NewGoogleCloudStorage(ctx context.Context, ProjectID string) (*GoogleCloudStorage, error) {
	gcs := new(GoogleCloudStorage)
	gcs.projectID = ProjectID

	// ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	gcs.client = client

	return gcs, nil
}

func (gcs *GoogleCloudStorage) Close() error {
	return gcs.client.Close()
}

func (gcs *GoogleCloudStorage) Write(ctx context.Context, bucketName string, objectPath string, text string) error {

	wc := gcs.client.Bucket(bucketName).Object(objectPath).NewWriter(ctx)

	_, err := fmt.Fprintf(wc, text)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("fmt.Fprintf: %v", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (gcs *GoogleCloudStorage) Copy(ctx context.Context, bucketName string, objectPath string, localFilePath string) error {

	// Open local file.
	f, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	// Copy GoogleCloudStorage
	wc := gcs.client.Bucket(bucketName).Object(objectPath).NewWriter(ctx)

	_, err = io.Copy(wc, f)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("io.Copy: %v", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("Copy.Close: %v", err)
	}

	return nil
}

// func sample() {

// 	// クライアントを作成する
// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// GCSオブジェクトを書き込むファイルの作成
// 	f, err := os.Create("sample.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// オブジェクトのReaderを作成
// 	bucketName := "xxx-bucket"
// 	objectPath := "yyy-obj"
// 	obj := client.Bucket(bucketName).Object(objectPath)
// 	reader, err := obj.NewReader(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer reader.Close()

// 	// 書き込み
// 	tee := io.TeeReader(reader, f)
// 	s := bufio.NewScanner(tee)
// 	for s.Scan() {
// 	}
// 	if err := s.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("done")
// }

// // uploadFile uploads an object.
// func UploadFile(w io.Writer, bucket, object string) error {
// 	bucket := "bucket-name"
// 	object := "object-name"
// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		return fmt.Errorf("storage.NewClient: %v", err)
// 	}
// 	defer client.Close()

// 	// Open local file.
// 	f, err := os.Open("notes.txt")
// 	if err != nil {
// 		return fmt.Errorf("os.Open: %v", err)
// 	}
// 	defer f.Close()

// 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
// 	defer cancel()

// 	// Upload an object with storage.Writer.
// 	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
// 	if _, err = io.Copy(wc, f); err != nil {
// 		return fmt.Errorf("io.Copy: %v", err)
// 	}
// 	if err := wc.Close(); err != nil {
// 		return fmt.Errorf("Writer.Close: %v", err)
// 	}
// 	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
// 	return nil
// }
