package app

import "appengine"
import "appengine/datastore"
import "crypto/sha256"
import "encoding/hex"
import "errors"

type entity struct {
  Data []byte `datastore:"data,noindex"`
}

func hashOK(key string, data []byte) bool {
  sha := sha256.New()
  sha.Write(data)
  sum := sha.Sum(nil)
  calculatedHash := hex.EncodeToString(sum)
  return calculatedHash == key
}

func makeDatastoreKey(c appengine.Context, key string) *datastore.Key {
  return datastore.NewKey(c, "data", key, 0, nil)
}

func shared(c appengine.Context, key string) (bool, error) {
  query := datastore.NewQuery("data").KeysOnly().Filter("__key__ =", makeDatastoreKey(c, key))
  count, e := query.Count(c)
  return count > 0, e
}

func share(c appengine.Context, key string, data []byte) error {
  if !hashOK(key, data) {
    return errors.New("hash doesn't match")
  }
  _, e := datastore.Put(c, makeDatastoreKey(c, key), &entity{data})
  return e
}

func fetch(c appengine.Context, key string) (*[]byte, error) {
  entity := entity{}
  e := datastore.Get(c, makeDatastoreKey(c, key), &entity)
  if e == datastore.ErrNoSuchEntity {
    return nil, nil
  } else if e != nil {
    return nil, e
  } else {
    return &entity.Data, nil
  }
}
