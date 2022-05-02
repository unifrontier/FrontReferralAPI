package repository

import (
	"FrontReferralAPI/entity"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DeviceRepository struct {
	Save    func(record *entity.Device) error
	Find    func(device_id string) (*entity.Device, error)
	FindAll func() ([]entity.Device, error)
}

// NewRepository returns a new repository
func NewRepository() DeviceRepository {
	return DeviceRepository{
		Save:    Save,
		Find:    Find,
		FindAll: FindAll,
	}
}

const (
	projectID      = "gofrontierreferrals"
	collectionName = "device_info"
)

func Save(record *entity.Device) error {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	_, err = client.Collection(collectionName).Doc(record.UniqueID).Set(ctx, record)
	if err != nil {
		log.Fatalf("Failed to set data: %v", err)
	}
	return nil
}

func Find(device_id string) (*entity.Device, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	doc := client.Collection(collectionName).Doc(device_id)
	docSnap, err := doc.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	var record entity.Device
	docSnap.DataTo(&record)
	return &record, nil
}

func FindAll() ([]entity.Device, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	iter := client.Collection(collectionName).Documents(ctx)
	var records []entity.Device
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		var record entity.Device
		doc.DataTo(&record)
		records = append(records, record)
	}
	return records, nil
}
