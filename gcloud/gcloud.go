package gcloud

import (
	"context"
	_ "embed"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var credJSON *[]byte

type FirebaseApp struct {
	app *firebase.App
	ctx context.Context
}

var app *FirebaseApp

func RegisterServiceAccount(json []byte) {
	credJSON = &json
}

func getFirebase() (*FirebaseApp, error) {
	if credJSON == nil {
		err := fmt.Errorf("credentials json must be registered first")
		return nil, err
	}

	if app == nil {
		ctx := context.Background()
		opt := option.WithCredentialsJSON(*credJSON)

		i_app, err := firebase.NewApp(ctx, nil, opt)

		if err != nil {
			return nil, err
		}

		app = &FirebaseApp{app: i_app, ctx: ctx}
	}

	return app, nil
}
