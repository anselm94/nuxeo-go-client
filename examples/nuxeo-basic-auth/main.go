package main

import (
	"context"
	"fmt"

	"github.com/anselm94/nuxeo"
	"github.com/anselm94/nuxeo/auth"
)

func main() {
	nuxeoClient := nuxeo.NewClient("https://demo.nuxeo.com/nuxeo",
		nuxeo.WithAuthenticator(auth.NewBasicAuthenticator("Administrator", "Administrator")),
	)

	ctx := context.Background()

	serverVersion, err := nuxeoClient.ServerVersion(ctx)
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("Nuxeo Server Version: %d.%d.%d", serverVersion.Major, serverVersion.Minor, serverVersion.Patch))

	currentUser, err := nuxeoClient.CurrentUser(ctx)
	if err != nil {
		panic(err)
	}
	println("Current User:", currentUser.Id, "Properties:", fmt.Sprintf("%v", currentUser.Properties))
}
