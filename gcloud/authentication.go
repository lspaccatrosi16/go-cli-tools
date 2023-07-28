package gcloud

import (
	"crypto/rand"
	"math"

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

func (a *AuthClient) GenerateTemporaryPassword(pLen int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	bytes := make([]byte, pLen)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	randomString := make([]byte, pLen)
	letterLen := float64(len(letters))

	for i, b := range bytes {
		ratio := float64(b / math.MaxUint8)
		letterIdx := int(math.Floor(ratio * letterLen))
		randomString[i] = letters[letterIdx-1]
	}

	return string(randomString), nil
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
