package gcloud

import (
	"firebase.google.com/go/v4/db"
)

type RTDBClient struct {
	app    *FirebaseApp
	Client *db.Client
}

func (r *RTDBClient) Set(path string, data interface{}) error {
	ref := r.Client.NewRef(path)
	err := ref.Set(r.app.ctx, data)
	if err != nil {
		return wrapRTDB(err)
	}

	return nil
}

func (r *RTDBClient) Read(path string) (interface{}, error) {
	ref := r.Client.NewRef(path)
	var data interface{}

	err := ref.Get(r.app.ctx, &data)
	if err != nil {
		return nil, wrapRTDB(err)
	}

	return data, nil
}

func NewDefaultRTDB() (*RTDBClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, wrapRTDB(err)
	}
	client, err := app.app.Database(app.ctx)
	if err != nil {
		return nil, wrapRTDB(err)
	}
	mClient := RTDBClient{
		app:    app,
		Client: client,
	}
	return &mClient, nil
}

func NewRTDB(url string) (*RTDBClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, wrapRTDB(err)
	}
	client, err := app.app.DatabaseWithURL(app.ctx, url)
	if err != nil {
		return nil, wrapRTDB(err)
	}
	mClient := RTDBClient{
		app:    app,
		Client: client,
	}
	return &mClient, nil
}
