package gcloud

import "firebase.google.com/go/v4/messaging"

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
		return wrapMessaging(err)
	}
	return nil
}

func (m *MessagingClient) SendDataMessage(data map[string]string, tokens []string) error {
	_, err := m.Client.SendMulticast(app.ctx, &messaging.MulticastMessage{
		Tokens: tokens,
		Data:   data,
	})

	if err != nil {
		return wrapMessaging(err)
	}
	return nil
}

func NewMessaging() (*MessagingClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, wrapMessaging(err)
	}

	client, err := app.app.Messaging(app.ctx)

	if err != nil {
		return nil, wrapMessaging(err)
	}

	mClient := MessagingClient{
		app:    app,
		Client: client,
	}
	return &mClient, nil
}
