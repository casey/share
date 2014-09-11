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
  datastoreKey := makeDatastoreKey(c, key)
  query := datastore.NewQuery("data").KeysOnly().Filter("__key__ =", datastoreKey)
  count, e := query.Count(c)
  return count > 0, e
}

func share(c appengine.Context, key string, data []byte) error {
  if !hashOK(key, data) {
    return errors.New("hash doesn't match")
  }

  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{data}
  _, e := datastore.Put(c, datastoreKey, &entity)
  return e
}

func fetch(c appengine.Context, key string) (*[]byte, error) {
  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{}
  e := datastore.Get(c, datastoreKey, &entity)
  if e == datastore.ErrNoSuchEntity {
    return nil, nil
  } else if e != nil {
    return nil, e
  } else {
    return &entity.Data, nil
  }
}
