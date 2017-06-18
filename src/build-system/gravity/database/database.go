package database

import (
  //"time"

  //"github.com/boltdb/bolt"
  "github.com/asdine/storm"
  "github.com/blevesearch/bleve"

  // Codec Options
  //"github.com/asdine/storm/codec/gob"
  //"github.com/asdine/storm/codec/json"
  //"github.com/asdine/storm/codec/sereal"
  //"github.com/asdine/storm/codec/protobuf"
  //"github.com/asdine/storm/codec/msgpack"
)

type Context struct {
  Opened        bool
  Path          string
  //Database    *bolt.DB
  Database      *storm.DB
  Indexes       map[string]SearchIndex
}

type SearchIndex struct {
  Name          string
  Path          string
  Index         bleve.Index
}

//self.Index, _  = bleve.Open(self.Path)
func (self *Context) InitiateSearchIndex(path string) (SearchIndex) {
  index, _       := bleve.New(path, bleve.NewIndexMapping())
  return SearchIndex{
    Path:     path,
    Index:    index,
  }
}

// TODO: Add going through the bleve.SearchResult and return a map or slice
// so the other package doesn't need to know about bleve
func (self *Context) SearchIndex(name, query string) {
  searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
  _, _ = self.Indexes[name].Index.Search(searchRequest)
}

func (self *Context) Open(path string, searchIndex bool) {
  var err error
  //var gobDb,      err = storm.Open("gob.db", storm.Codec(gob.Codec))
  //var jsonDb,     err = storm.Open("json.db", storm.Codec(json.Codec))
  //var serealDb,   err = storm.Open("sereal.db", storm.Codec(sereal.Codec))
  //var msgpackDb,  err = storm.Open("msgpack.db", storm.Codec(msgpack.Codec))
  //self.Bolt, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
  self.Database, err = storm.Open(path, storm.AutoIncrement(), storm.Batch())
  if err == nil {
    self.Opened = true
  }
  if searchIndex {
    self.Indexes = make(map[string]SearchIndex)
  }
}

func (self *Context) Close() error {
  self.Opened = false
  return self.Database.Close()
}
