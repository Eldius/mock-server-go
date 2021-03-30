package request

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/Eldius/mock-server-go/logger"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	db                   *bolt.DB
	log                  = logger.Log()
	requestsDbBucketName = "requests"
)

func init() {
	var err error
	db, err = bolt.Open("mocky.db", 0666, nil)
	if err != nil {
		fmt.Println("Failed to open db file")
		panic(err)
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(requestsDbBucketName))
		if err != nil {
			fmt.Println("Failed to create requests bucket")
			panic(err)
		}
		return nil
	}); err != nil {
		fmt.Println("Failed to open transaction to initialize bucket")
		panic(err)
	}
}

func Persist(r *Record) {
	if db == nil {
		panic("DB is nil")
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(requestsDbBucketName))
		id, err := b.NextSequence()
		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{
					"record": r,
				}).
				Error("Failed to get next sequential ID")
			return err
		}
		r.ID = int(id)
		bin, err := json.Marshal(r)
		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{
					"record": r,
				}).
				Error("Failed to marshal request json")
			return err
		}
		err = b.Put(itob(r.ID), bin)
		if err != nil {
			log.WithError(err).
				WithFields(logrus.Fields{
					"record": r,
				}).
				Error("Failed to marshal request json")
			return err
		}
		return err
	}); err != nil {
		log.WithError(err).
			WithFields(logrus.Fields{
				"record": r,
			}).
			Error("Failed to open transaction")
	}
}

func GetRequests() []Record {
	records := make([]Record, 0)
	if err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(requestsDbBucketName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			var r Record
			err := json.Unmarshal(v, &r)
			if err != nil {
				log.WithError(err).
					WithFields(logrus.Fields{
						"value": string(v),
						"key":   string(k),
					}).
					Error("Failed to marshal request json")
			}
			records = append(records, r)
		}

		return nil
	}); err != nil {
		log.WithError(err).
			Error("Failed to open View Transaction")
	}
	return records
}

func debug(obj interface{}) {
	log.Debug(obj)
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
