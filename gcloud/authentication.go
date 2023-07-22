package gcloud

import (
	"firebase.google.com/go/auth"
	"google.golang.org/api/iterator"
)

type AuthClient struct {
	app    *FirebaseApp
	Client *auth.Client
}

type UserRecord = *auth.UserRecord
type UserRecords = *[]UserRecord

func (a *AuthClient) GetUser(uid string) (*auth.UserRecord, error) {
	u, err := a.Client.GetUser(app.ctx, uid)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a *AuthClient) GetUsers() (UserRecords, error) {
	userList := []*auth.UserRecord{}

	iter := a.Client.Users(a.app.ctx, "")

	for {
		user, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, errorf("error listing users")
		}

		userList = append(userList, user.UserRecord)
	}

	return &userList, nil
}

func NewFirebaseAuth() (*AuthClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, wrap(err)
	}

	client, err := app.app.Auth(app.ctx)

	if err != nil {
		return nil, wrap(err)
	}

	aClient := AuthClient{
		app:    app,
		Client: client,
	}

	return &aClient, nil

}
