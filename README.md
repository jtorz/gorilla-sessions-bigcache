gorilla-sessions-bigcache
=========================

Fork github.com/bradleypeabody/gorilla-sessions-memcache

[Bigcache](github.com/allegro/bigcache) session support for Gorilla Web Toolkit.

Dependencies
------------

The usual gorilla stuff:

    go get github.com/gorilla/sessions

For an ASCII bigcache client:

    go get "github.com/allegro/bigcache/v2"

Usage
-----

```go
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/allegro/bigcache/v2"
  gsb "github.com/jtorz/gorilla-sessions-bigcache"
)

func main() {
  bigcacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
  if err != nil {
    panic(err)
  }
  store := gsb.NewBigcacheStore(bigcacheClient, "session_prefix_", []byte("secret"))
  runServer(store)
}

func runServer(store *gsb.BigcacheStore) {
  http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
      fmt.Fprintf(w, "Error: %v", err)
      return
    }
    fmt.Fprintf(w, "Got: %v", session)
  })

  http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
      fmt.Fprintf(w, "Error: %v", err)
      return
    }
    session.Values["foo"] = "bar"
    session.Values[42] = 43
    session.Save(r, w)
    fmt.Fprint(w, "ok")
  })
  http.ListenAndServe(":8080", nil)
}
```

You can also setup a BigcacheStore, which does not rely on the browser accepting cookies.
this means, your client has to extract and send a configurable http Headerfield manually.

```go
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/allegro/bigcache/v2"
  gsb "github.com/jtorz/gorilla-sessions-bigcache"
)

func main() {
  bigcacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
  if err != nil {
    panic(err)
  }
  store := gsb.NewBigCacherStoreWithValueStorer(gsb.NewGoBigcacher(bigcacheClient), &gsb.HeaderStorer{HeaderFieldName: "X-CUSTOM-HEADER"}, "session_prefix_", []byte
  runServer(store)
}

func runServer(store *gsb.BigcacheStore) {
  http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
      fmt.Fprintf(w, "Error: %v", err)
      return
    }
    fmt.Fprintf(w, "Got: %v", session)
  })

  http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-name")
    if err != nil {
      fmt.Fprintf(w, "Error: %v", err)
      return
    }
    session.Values["foo"] = "bar"
    session.Values[42] = 43
    session.Save(r, w)
    fmt.Fprint(w, "ok")
  })
  http.ListenAndServe(":8080", nil)
}
```

Storage Methods
---------------

I've added a few different methods of storage of the session data in bigcache.  You
use them by setting the StoreMethod field.

* SecureCookie - uses the default securecookie encoding.  Values are more secure
  as they are not readable from bigcache without the secret key.
* Gob - uses the Gob encoder directly without any post processing.  Faster.
  Result is Gob's usual binary gibber (not human readable)
* Json - uses the Json Marshaller.  Result is human readable, slower but still
  pretty fast.  Be careful - it will munch your data into stuff that works
  with JSON, and the keys must be strings.  Example: you put in an int64 value
  and you'll get back a float64.

Example:

```go
store := gsb.NewBigCacherStore(bigcacheClient, "session_prefix_", []byte("..."))
// do one of these:
store.StoreMethod = gsb.StoreMethodSecureCookie // default, more secure
store.StoreMethod = gsb.StoreMethodGob // faster
store.StoreMethod = gsb.StoreMethodJson // human readable
              // (but watch out, it munches your types
              // to JSON compatible stuff)
```

Logging
-------

Logging is available by setting the Logging field to > 0 after making your BigcacheStore.

```go
store := gsb.NewBigCacherStore(bigcacheClient, "session_prefix_", []byte("..."))
store.Logging = 1
```

That will output (using `log.Printf`) data about each session read/written from/to bigcache.
Useful for debugging

<!--

Things to Know
--------------
 markdownlint-disable MD000
* No official release has been done of this package but it should be stable for production use.

* You can also call NewDumbMemorySessionStore() for local development without a bigcache server (it's a stub that just stuffs your session data in a map - definitely do not use this for anything but local dev and testing).
-->
