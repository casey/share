package app

import "crypto/sha256"
import "encoding/hex"
import "appengine"
import "appengine/datastore"

type entity struct {
  Data string `datastore:"data,noindex"`
}

func stringID(key string) string {
  sha := sha256.New()
  sha.Write([]byte(key))
  sum := sha.Sum(nil)
  return hex.EncodeToString(sum)
}

func makeDatastoreKey(c appengine.Context, key string) *datastore.Key {
  return datastore.NewKey(c, "data", stringID(key), 0, nil)
}

func getData(c appengine.Context, key string) (*string, error) {
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

func putData(c appengine.Context, key string, data string) (*string, error) {
  datastoreKey := makeDatastoreKey(c, key)
  entity := entity{data}
  _, e := datastore.Put(c, datastoreKey, &entity)
  if e == nil {
    return &entity.Value, nil
  } else {
    return nil, e
  }
}
