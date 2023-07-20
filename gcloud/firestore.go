package gcloud

import (
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreClient struct {
	app    *FirebaseApp
	Client *firestore.Client
}

func (f *FirestoreClient) Close() {
	f.Client.Close()
}

func (f *FirestoreClient) SetDoc(path string, data map[string]interface{}) error {
	split := strings.Split(path, "/")

	if len(split)%2 != 0 {
		err := fmt.Errorf("incomplete path: must be col/doc/col/doc..../doc, not %s", path)

		return err
	}

	var docRef *firestore.DocumentRef

	for i := 0; i < len(split); i += 2 {
		if docRef == nil {
			docRef = f.Client.Collection(split[i]).Doc(split[i+1])
		} else {
			docRef = docRef.Collection(split[i]).Doc(split[i+1])
		}
	}

	_, err := docRef.Set(f.app.ctx, data)

	if err != nil {
		return err
	}

	return nil
}

func (f *FirestoreClient) GetDoc(path string) (*map[string]interface{}, error) {
	split := strings.Split(path, "/")

	if len(split)%2 != 0 {
		err := fmt.Errorf("incomplete path: must be col/doc/col/doc..../doc, not %s", path)
		return nil, err
	}

	var docRef *firestore.DocumentRef

	for i := 0; i < len(split); i += 2 {
		if docRef == nil {
			docRef = f.Client.Collection(split[i]).Doc(split[i+1])
		} else {
			docRef = docRef.Collection(split[i]).Doc(split[i+1])
		}
	}

	doc, err := docRef.Get(f.app.ctx)

	if err != nil {
		return nil, err
	}

	if err != nil {
		if status.Code(err) == codes.NotFound {
			err := fmt.Errorf("document does not exist %s", path)
			return nil, err
		} else {
			return nil, err
		}
	}

	if !doc.Exists() {
		err := fmt.Errorf("document does not exist %s", path)

		return nil, err
	}

	data := doc.Data()
	return &data, nil
}

func (f *FirestoreClient) GetManyDocs(path string) ([]*firestore.DocumentSnapshot, error) {
	split := strings.Split(path, "/")

	if len(split)%2 != 1 {
		err := fmt.Errorf("incomplete path: must be col/doc/col/doc..../col, not %s", path)
		return []*firestore.DocumentSnapshot{}, err
	}

	var colRef *firestore.CollectionRef

	for i := 0; i < len(split); i += 2 {
		if colRef == nil {
			colRef = f.Client.Collection(split[i])
			i--
		} else {
			colRef = colRef.Doc(split[i]).Collection(split[i+1])
		}
	}

	datas := []*firestore.DocumentSnapshot{}
	docs := colRef.Documents(f.app.ctx)

	for {
		doc, err := docs.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return []*firestore.DocumentSnapshot{}, err
		}

		datas = append(datas, doc)

	}

	return datas, nil

}

func NewFirestore() (*FirestoreClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, err
	}

	client, err := app.app.Firestore(app.ctx)

	if err != nil {
		return nil, err
	}

	fsClient := FirestoreClient{
		app:    app,
		Client: client,
	}

	return &fsClient, nil
}
