package gsb_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v2"
	gsb "github.com/jtorz/gorilla-sessions-bigcache"
)

func ExampleNewBigcacheStore() {
	bigcacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}
	store := gsb.NewBigcacheStore(bigcacheClient, "session_prefix_", []byte("secret"))
	//store := gsb.NewBigCacherStoreWithValueStorer(gsb.NewGoBigcacher(bigcacheClient), &gsb.HeaderStorer{HeaderFieldName: "X-CUSTOM-HEADER"}, "session_prefix_", []byte("secret-key-goes-here"))
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
