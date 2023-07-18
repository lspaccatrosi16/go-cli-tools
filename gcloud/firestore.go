package gcloud

import (
	"fmt"
	"log"
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

func (f *FirestoreClient) SetDoc(path string, data map[string]interface{}) {
	split := strings.Split(path, "/")

	if len(split)%2 != 0 {
		errMsg := fmt.Sprintf("incomplete path: must be col/doc/col/doc..../doc, not %s", path)
		panic(errMsg)
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
		panic(err)
	}
}

func (f *FirestoreClient) GetDoc(path string) map[string]interface{} {
	split := strings.Split(path, "/")

	if len(split)%2 != 0 {
		errMsg := fmt.Sprintf("incomplete path: must be col/doc/col/doc..../doc, not %s", path)
		panic(errMsg)
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
		panic(err)
	}

	if err != nil {
		if status.Code(err) == codes.NotFound {
			errMsg := fmt.Sprintf("document does not exist %s", path)
			panic(errMsg)
		} else {
			log.Fatalln(err)
		}
	}

	if !doc.Exists() {
		errMsg := fmt.Sprintf("document does not exist %s", path)
		panic(errMsg)
	}

	data := doc.Data()
	return data
}

func (f *FirestoreClient) GetManyDocs(path string) []*firestore.DocumentSnapshot {
	split := strings.Split(path, "/")

	if len(split)%2 != 1 {
		errMsg := fmt.Sprintf("incomplete path: must be col/doc/col/doc..../col, not %s", path)
		panic(errMsg)
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
			panic(err)
		}

		datas = append(datas, doc)

	}

	return datas

}

func NewFirestore() *FirestoreClient {
	app := getFirebase()

	client, err := app.app.Firestore(app.ctx)

	if err != nil {
		panic(err)
	}

	fsClient := FirestoreClient{
		app:    app,
		Client: client,
	}

	return &fsClient
}
