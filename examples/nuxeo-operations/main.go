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

	// initialize operation manager
	operationManager := nuxeoClient.OperationManager()

	////////////////
	// Operation 1 : fetch a document
	////////////////
	operation := nuxeo.NewOperation("Document.Fetch")
	operation.SetParam("value", "797d8306-d3b7-410e-b30e-5972c2cb6eb7")

	opRes, err := operationManager.Execute(ctx, *operation, nil)
	if err != nil {
		panic(err)
	}

	document, err := opRes.AsDocument()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fetched document: %s (type: %s)\n", document.Title, document.Type)

	////////////////
	// Operation 2 : fetch a list of documents
	////////////////
	operation2 := nuxeo.NewOperation("Document.Query")
	operation2.SetParam("query", "SELECT * FROM Document")

	opRes2, err := operationManager.Execute(ctx, *operation2, nil)
	if err != nil {
		panic(err)
	}

	documents, err := opRes2.AsDocumentList()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fetched %d documents:\n", len(documents.Entries))
	for _, doc := range documents.Entries {
		fmt.Printf("- %s (type: %s)\n", doc.Title, doc.Type)
	}

	////////////////
	// Operation 3 : fetch a blob (file content)
	////////////////
	operation3 := nuxeo.NewOperation("Document.GetBlob")
	operation3.SetInputDocumentId("d6a4663d-d6d0-47c7-9cde-3c8a71e78f3f")
	operation3.SetParam("xpath", nuxeo.DocumentPropertyFileContent)

	opRes3, err := operationManager.Execute(ctx, *operation3, nil)
	if err != nil {
		panic(err)
	}

	blob, err := opRes3.AsBlob()
	if err != nil {
		panic(err)
	}
	defer blob.Close()

	println("Fetched blob: " + blob.Filename)

	////////////////
	// Operation 4 : fetch a list of blobs
	////////////////
	operation4 := nuxeo.NewOperation("Document.GetBlobs")
	operation4.SetInputDocumentId("d6a4663d-d6d0-47c7-9cde-3c8a71e78f3f")

	opRes4, err := operationManager.Execute(ctx, *operation4, nil)
	if err != nil {
		panic(err)
	}

	blobs, err := opRes4.AsBlobList()
	if err != nil {
		panic(err)
	}

	println("Fetched blobs:")
	for blob := range blobs { // iterate over blobs
		defer blob.Close()

		println("- " + blob.Filename)

		// write to local file
		file, err := os.Create(blob.Filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = io.Copy(file, blob)
		if err != nil {
			panic(err)
		}
	}
}
