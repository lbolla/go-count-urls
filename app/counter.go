package counter

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"time"
)

type Counter struct {
	Path string
	Count int
	Timestamp time.Time
}

func getEmptyCounter(path string) Counter {
	return Counter{Path: path, Count: 0, Timestamp: time.Now()}
}

func inc(c appengine.Context, key *datastore.Key, path string) (Counter, error) {
	var x Counter

	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return getEmptyCounter(path), err
	}

	x.Path = path
	x.Count++
	x.Timestamp = time.Now()

	if _, err := datastore.Put(c, key, &x); err != nil {
		return getEmptyCounter(path), err
	}

	return x, nil
}

func handle(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Path[1:]
	if key == "" {
		http.NotFound(w, r)
		return
	}

	var count Counter
	c := appengine.NewContext(r)
	err := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var err1 error
		count, err1 = inc(c, datastore.NewKey(c, key, "singleton", 0, nil), r.URL.Path)
		return err1
	}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Path=%s, Count=%d, When=%s", count.Path, count.Count, count.Timestamp)
}

func init() {
	http.HandleFunc("/", handle)
}
