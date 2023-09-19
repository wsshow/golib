package bbolt

import (
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"time"

	bolt "go.etcd.io/bbolt"
)

type BBolt struct {
	db         *bolt.DB
	buckets    map[string]string
	bucketName string
	structObj  any
	err        error
	debug      bool
}

type Options struct {
	Timeout time.Duration
	Debug   bool
}

func NewBBolt(path string, mode fs.FileMode, options *Options) *BBolt {
	boltOptions := &bolt.Options{Timeout: options.Timeout}
	bboltDB, err := bolt.Open(path, mode, boltOptions)
	if err != nil {
		panic(err)
	}
	bbolt := &BBolt{
		db:      bboltDB,
		buckets: make(map[string]string),
		err:     nil,
		debug:   options.Debug,
	}
	return bbolt
}

func (bb *BBolt) Close() *BBolt {
	bb.err = bb.db.Close()
	return bb
}

func (bb *BBolt) CreateBuckets(structObjs ...any) *BBolt {
	for _, structObj := range structObjs {
		structType := reflect.TypeOf(structObj)
		if structType.Kind() != reflect.Struct {
			panic(fmt.Sprintf("require type: struct, current type: %s", structType.Kind()))
		}
		bb.createBucket(structType.Name())
		if bb.debug {
			fmt.Println("create bucket:", structType.Name())
		}
	}
	return bb
}

func (bb *BBolt) createBucket(name string) *BBolt {
	if bb.err != nil {
		return bb
	}
	bb.err = bb.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		bb.buckets[name] = name
		bb.bucketName = name
		return nil
	})
	return bb
}

func (bb *BBolt) Bucket(structObj any) *BBolt {
	if bb.err != nil {
		return bb
	}
	bb.structObj = structObj
	structType := reflect.TypeOf(structObj)
	if structType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("require type: struct, current type: %s", structType.Kind()))
	}
	name := structType.Name()
	if bucket, ok := bb.buckets[name]; !ok {
		bb.err = errors.New("unknown bucket name: " + name)
	} else {
		bb.bucketName = bucket
	}
	return bb
}

func (bb *BBolt) Set(key string, bs []byte) *BBolt {
	if bb.err != nil {
		return bb
	}
	bb.err = bb.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bb.bucketName))
		err := bucket.Put([]byte(key), bs)
		return err
	})
	return bb
}

func (bb *BBolt) Del(key string) *BBolt {
	if bb.err != nil {
		return bb
	}
	bb.err = bb.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bb.bucketName))
		err := bucket.Delete([]byte(key))
		return err
	})
	return bb
}

func (bb *BBolt) Get(key string) (bs []byte) {
	if bb.err != nil {
		return
	}
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bb.bucketName))
		bs = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		bb.err = err
		return nil
	}
	return
}

func (bb *BBolt) GetAll() (bss [][]byte) {
	if bb.err != nil {
		return
	}
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bb.bucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			bss = append(bss, v)
		}
		return nil
	})
	if err != nil {
		bb.err = err
		return nil
	}
	return
}

func (bb *BBolt) Error() error {
	return bb.err
}
