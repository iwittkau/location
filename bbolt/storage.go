package bbolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/iwittkau/location"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketCheckIns = "CheckIns"
	checkinPrefix  = "checkin"
	tsPrefix       = "ts"
)

var _ location.Storage = &Storage{}

type Storage struct {
	db *bolt.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Open(path string) error {

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return err
	}
	s.db = db

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketCheckIns))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return err
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) CreateCheckIn(c location.CheckIn) error {

	if ok, err := s.CheckInExists(c); err != nil {
		return err
	} else if ok {
		return errors.New("checkin exists")
	}

	if c.Time.IsZero() {
		c.Time = time.Now()
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketCheckIns))
		if err := b.Put([]byte(checkInTS(c)), []byte(data)); err != nil {
			return err
		}
		if err := b.Put([]byte(checkInID(c.ID)), []byte(checkInTS(c))); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) DeleteCheckIn(id string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketCheckIns))
		if tsData := b.Get([]byte(checkInID(id))); tsData != nil {
			if err := b.Delete(tsData); err != nil {
				return err
			}
		}
		if err := b.Delete([]byte(checkInID(id))); err != nil {
			return err
		}

		return nil
	})
	return err
}

func (s *Storage) CheckIn(id string) (result location.CheckIn, err error) {
	var data []byte
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketCheckIns))
		tsData := b.Get([]byte(checkInID(id)))
		data = b.Get(tsData)
		return nil
	})
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return
}

func (s *Storage) CheckInExists(c location.CheckIn) (bool, error) {
	var data []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketCheckIns))
		data = b.Get([]byte(checkInTS(c)))
		return nil
	})
	return data != nil, err
}

func (s *Storage) ListCheckIns(size int, offset int) (result []location.CheckIn, err error) {
	data := [][]byte{}
	err = s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(bucketCheckIns)).Cursor()

		prefix := []byte(tsPrefix)
		i := 0
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			i++
			if i <= offset {
				continue
			}
			data = append(data, v)
			if len(data) == size {
				break
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	result = make([]location.CheckIn, len(data))
	c := location.CheckIn{}
	for i := range data {
		if err := json.Unmarshal(data[i], &c); err != nil {
			return nil, err
		}
		result[i] = c
	}
	return
}

func checkInID(id string) string {
	return fmt.Sprintf("%s-%s", checkinPrefix, id)
}

func checkInTS(c location.CheckIn) string {
	return fmt.Sprintf("%s-%d", tsPrefix, c.Time.UnixNano())
}
