package repository

import (
	"FrontReferralAPI/entity"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DeviceRepository struct {
	Save             func(record *entity.Device) error
	FindByReferrer   func(referrer_id string) (*entity.Device, error)
	FindAll          func() ([]entity.Device, error)
	Update           func(referrer_id string, device_id string)
	FindDevice       func(device_id string) (*entity.Device, error)
	CountReferredIDS func(device_id string) int
}

// NewRepository returns a new repository
func NewRepository() DeviceRepository {
	return DeviceRepository{
		Save:             Save,
		FindByReferrer:   FindByReferrer,
		FindAll:          FindAll,
		Update:           Update,
		FindDevice:       FindDevice,
		CountReferredIDS: CountReferredIDS,
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

func FindByReferrer(referrer_id string) (*entity.Device, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	fmt.Println("referrer_id: ", referrer_id)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	doc := client.Collection(collectionName).Doc(referrer_id)
	docSnap, err := doc.Get(ctx)
	if err != nil {
		fmt.Println("Referrer not Found")
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

func Update(referrer_id string, device_id string) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	doc := client.Collection(collectionName).Doc(referrer_id)
	docSnap, err := doc.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	var record entity.Device
	docSnap.DataTo(&record)
	record.ReferredIDS = append(record.ReferredIDS, device_id)
	_, err = doc.Set(ctx, record)
	if err != nil {
		log.Fatalf("Failed to set data: %v", err)
	}
}

func FindDevice(device_id string) (*entity.Device, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	iter := client.Collection(collectionName).Where("DeviceID", "==", device_id).Documents(ctx)
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
		return &record, nil
	}
	return nil, nil
}

//count the number of referredids in the device
func CountReferredIDS(device_id string) int {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	iter := client.Collection(collectionName).Where("DeviceID", "==", device_id).Documents(ctx)
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
		return len(record.ReferredIDS)
	}
	return 0
}
