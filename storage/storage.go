package storage

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Storage struct {
	size int
	pool chan *mgo.Session
}

func New(info *mgo.DialInfo) (*Storage, error) {
	if info.PoolLimit == 0 {
		info.PoolLimit = 16
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	pool := make(chan *mgo.Session, info.PoolLimit)
	size := 1
	for size < info.PoolLimit {
		size++
		pool <- session.Copy()
	}
	pool <- session
	return &Storage{size: size, pool: pool}, nil
}

func (s *Storage) Close() {
	for s.size > 0 {
		s.size--
		s.acquire().Close()
	}
}

func (s *Storage) acquire() *mgo.Session {
	return <-s.pool
}

type table struct {
	*Storage
	*mgo.Collection
}

func (t table) release() {
	t.pool <- t.Database.Session
}

func (s *Storage) table(collection string) table {
	return table{s, s.acquire().DB("").C(collection)}
}

const defaultCollection = "ads"

type A struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Date        *time.Time    `json:"date,omitempty" bson:"date"`
	Price       float32       `json:"price" bson:"price"`
	Subject     string        `json:"subject,omitempty" bson:"subject"`
	Description string        `json:"description,omitempty" bson:"description"`
	Photo       []string      `json:"photo,omitempty" bson:"photo"`
}

func (s *Storage) Search(ids []string, fields []string, limit int, offset int, sorts []string) (interface{}, error) {
	t := s.table(defaultCollection)
	defer t.release()
	var e bson.M
	var f bson.M
	if len(ids) > 0 {
		var ins []bson.ObjectId
		for _, v := range ids {
			d, err := hex.DecodeString(v)
			if err != nil {
				return nil, err
			}
			if len(d) != 12 {
				return nil, fmt.Errorf("invalid object id %q", v)
			}
			ins = append(ins, bson.ObjectId(d))
		}
		e = bson.M{"_id": bson.M{"$in": ins}}
	}
	if len(fields) == 0 && len(ids) == 0 {
		f = bson.M{"date": true, "price": true, "subject": true, "photo": bson.M{"$slice": 1}}
	} else {
		f = bson.M{"date": false, "photo": false, "description": false}
		for _, v := range fields {
			switch v {
			case "date", "photo", "description":
				delete(f, v)
			}
		}
	}
	v := make([]A, 0, len(ids))
	err := t.Find(e).Select(f).Sort(sorts...).Limit(limit).Skip(offset).All(&v)
	if err != nil {
		return nil, err
	} else if len(v) == 0 && len(ids) == 1 {
		return nil, mgo.ErrNotFound
	} else {
		return v, nil
	}
}

func (s *Storage) Store(price float32, subject string, description string, photos []string) (interface{}, error) {
	t := s.table(defaultCollection)
	defer t.release()
	n := time.Now()
	v := A{
		Id:          bson.NewObjectId(),
		Date:        &n,
		Price:       price,
		Subject:     subject,
		Description: description,
		Photo:       photos,
	}
	if err := t.Insert(v); err != nil {
		return nil, err
	}
	return v, nil
}
