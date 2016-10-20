package db

import (
	"fmt"
    "log"
	"time"

	"github.com/boltdb/bolt"
)

var (
    db      *bolt.DB
    today   string
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
    now := time.Now()
    before := now.AddDate(0, 0, -1)
    nYear, nMonth, nDay := now.Date()
    yYear, yMonth, yDay := before.Date()
    today = fmt.Sprintf("%d/%d/%d", nYear, nMonth, nDay)
    yesterday := fmt.Sprintf("%d/%d/%d", yYear, yMonth, yDay)
    log.Printf("Initiating bolt connection... Current bucket name is %s\n", today)
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	panicOnError(err)
    log.Printf("DB opened succesfully")
    err = createBucket("Items")
    panicOnError(err)
    log.Printf("Items bucket created")
    err = createBucket(today)
    panicOnError(err)
    log.Printf("Today's bucket created")
    err = deleteBucket(yesterday)
    logOnError(err)
    log.Printf("Yesterday's bucket deleted")
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