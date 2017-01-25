package main

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

var (
	db *bolt.DB

	today, yesterday string
)

// InitDB initializes DB
func InitDB() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	today = getStringDateFromToday(0)
	yesterday = getStringDateFromToday(-1)
	glog.Infof("Initiating bolt connection... Current bucket name is %s\n", today)
	if db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
		glog.Fatalf("Unable to open db: %s", err)
	}
	if err = createBucket("Items"); err != nil {
		glog.Fatalf("Unable to create bucket Items: %s", err)
	}
	if err = createBucket(today); err != nil {
		glog.Fatalf("Unable to create bucket %s: %s", today, err)
	}
	if err = deleteBucket(yesterday); err != nil && err.Error() != "bucket not found" {
		glog.Fatalf("Unable to delete bucket %s: %s", yesterday, err)
	}
}

// CloseDB closes the connection
func CloseDB() {
	db.Close()
	glog.Flush()
}

// DBWatcher runs on schedule
func DBWatcher() {
	for {
		yesterday := getStringDateFromToday(-1)
		if today == yesterday {
			CloseDB()
			InitDB()
		}
		time.Sleep(5 * time.Minute)
	}
}

// Get returns the value for key
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

// Put writes the vakue to DB
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
