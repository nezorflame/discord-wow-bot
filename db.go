package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db    *bolt.DB
	today string
)

func logOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func panicOnError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func Init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	today = getStringDateFromToday(0)
	yesterday := getStringDateFromToday(-1)
	log.Printf("Initiating bolt connection... Current bucket name is %s\n", today)
	db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	panicOnError(err)
	err = createBucket("Items")
	panicOnError(err)
	err = createBucket(today)
	panicOnError(err)
	err = deleteBucket(yesterday)
	logOnError(err)
}

func Close() {
	db.Close()
}

func Watcher() {
	for {
		yesterday := getStringDateFromToday(-1)
		if today == yesterday {
			Close()
			Init()
		}
		time.Sleep(5 * time.Minute)
	}
}

func Get(bucketName, key string) (value []byte) {
	bName := bucketName
	if bucketName == "Main" {
		bName = today
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bName))
		value = b.Get([]byte(key))
		return nil
	})
	return
}

func Put(bucketName, key string, value []byte) (err error) {
	bName := bucketName
	if bucketName == "Main" {
		bName = today
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bName))
		err = b.Put([]byte(key), value)
		return err
	})
	return
}

func createBucket(name string) (err error) {
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
	return
}

func deleteBucket(name string) (err error) {
	db.Update(func(tx *bolt.Tx) error {
		err = tx.DeleteBucket([]byte(name))
		return err
	})
	return
}

func getStringDateFromToday(diff int) string {
	dest := time.Now().AddDate(0, 0, diff)
	year, month, day := dest.Date()
	return fmt.Sprintf("%d/%d/%d", year, month, day)
}
