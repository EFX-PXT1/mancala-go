package game

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v2"
	bh "github.com/timshannon/badgerhold/v2"
)

type Item struct {
	ID       int
	Category string `badgerholdIndex:"Category"`
	Created  time.Time
}

func TestBadgerHold(t *testing.T) {
	//	assert := assert.New(t)

	data := []Item{
		{
			ID:       0,
			Category: "blue",
			Created:  time.Now().Add(-4 * time.Hour),
		},
		{
			ID:       1,
			Category: "red",
			Created:  time.Now().Add(-3 * time.Hour),
		},
		{
			ID:       2,
			Category: "blue",
			Created:  time.Now().Add(-2 * time.Hour),
		},
		{
			ID:       3,
			Category: "blue",
			Created:  time.Now().Add(-20 * time.Minute),
		},
	}

	dir := tempdir()
	defer os.RemoveAll(dir)

	options := bh.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir
	options.Logger = nil
	store, err := bh.Open(options)
	defer store.Close()

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// insert the data in one transaction

	err = store.Badger().Update(func(tx *badger.Txn) error {
		for i := range data {
			err := store.TxInsert(tx, data[i].ID, data[i])
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// Find all items in the blue category that have been created in the past hour
	var result []Item

	err = store.Find(&result, bh.Where("Category").Eq("blue").And("Created").Ge(time.Now().Add(-1*time.Hour)))

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	fmt.Println(result[0].ID)
	// Output: 3

}

// tempdir returns a temporary dir path.
func tempdir() string {
	name, err := ioutil.TempDir("", "badgerhold-")
	if err != nil {
		panic(err)
	}
	return name
}
