package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/anselm94/nuxeo-go-client"
	nuxeoauth "github.com/anselm94/nuxeo-go-client/auth"
)

func main() {
	ctx := context.Background()

	// client options
	nuxeoClientOptions := nuxeo.DefaultNuxeoClientOptions()

	// basic authenticator
	nuxeoClientOptions.Authenticator = nuxeoauth.NewBasicAuthenticator("Administrator", "Administrator")

	// initialize client
	nuxeoClient := nuxeo.NewClient("https://demo.nuxeo.com/nuxeo", &nuxeoClientOptions)

	/////////

	repo := nuxeoClient.Repository()
	uploadManager := nuxeoClient.BatchUploadManager()

	// create a batch
	batch, err := uploadManager.CreateBatch(ctx, nil)
	if err != nil {
		panic(err)
	}

	// read the file along with its size
	file, err := os.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileLength, _ := file.Stat()

	// upload the file to the batch
	uploadOpts := nuxeo.NewUploadOptions("example.pdf", fileLength.Size(), "application/pdf")
	batchUploadInfo, err := uploadManager.Upload(ctx, batch.BatchId, "0", uploadOpts, file, nil)
	if err != nil {
		panic(err)
	}

	// create a document with the uploaded blob
	newDoc := nuxeo.NewDocument("File", "example.pdf")
	newDoc.SetProperty(nuxeo.DocumentPropertyDCDescription, nuxeo.NewStringField("An example PDF file"))
	newDoc.SetUploadInfoProperty(nuxeo.DocumentPropertyFileContent, nuxeo.UploadInfo{
		Batch:  batchUploadInfo.BatchId,
		FileId: batchUploadInfo.FileIdx,
	})
	createdDoc, err := repo.CreateDocumentById(ctx, "7ee74b3c-ab1f-4213-9467-6b68f64a4f88", *newDoc, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created document:", createdDoc.Title)

	// download the file content
	blob, err := repo.StreamBlobById(ctx, createdDoc.ID, nuxeo.DocumentPropertyFileContent, nil)
	if err != nil {
		panic(err)
	}
	defer blob.Close()

	// create a local file
	outFile, err := os.Create("downloaded_example.pdf")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// copy the content from the blob to the local file
	_, err = io.Copy(outFile, blob)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded file content to: downloaded_example.pdf")

	// // cleanup: delete the created document
	// err = repo.DeleteDocument(ctx, createdDoc.ID)
	// if err != nil {
	// 	panic(err)
	// }
}
