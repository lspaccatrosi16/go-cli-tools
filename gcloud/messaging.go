package gcloud

import "firebase.google.com/go/messaging"

type MessagingClient struct {
	app    *FirebaseApp
	Client *messaging.Client
}

func (m *MessagingClient) SendMessage(title string, body string, data map[string]string, image string, tokens []string) error {
	_, err := m.Client.SendMulticast(app.ctx, &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: image,
		},
		Data: data,
	})

	if err != nil {
		return wrap(err)
	}
	return nil
}

func NewMessaging() (*MessagingClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, wrap(err)
	}

	client, err := app.app.Messaging(app.ctx)

	if err != nil {
		return nil, wrap(err)
	}

	mClient := MessagingClient{
		app:    app,
		Client: client,
	}
	return &mClient, nil
}
