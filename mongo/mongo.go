package mongo

import (
	"fmt"
	"time"

	"github.com/clouway/godb"
	"gopkg.in/mgo.v2"
)

type database struct {
	mgoSess *mgo.Session
	mgoDB   *mgo.Database
}

// NewDatabase establishes a new database connection using
// configuration options in the provided config.
func NewDatabase(config *godb.Config) godb.Database {
	info := &mgo.DialInfo{Addrs: config.Addrs, Timeout: 5 * time.Second}
	sess, err := mgo.DialWithInfo(info)

	if err != nil {
		panic(fmt.Errorf("unable to connect to host '%s', failed with: %v", info.Addrs, err))
	}

	db := sess.DB(config.Database)

	return &database{sess, db}
}

func (db *database) Close() {
	db.mgoSess.Close()
}

func (db *database) Collection(name string) godb.Collection {
	coll := db.mgoDB.C(name)
	return &collection{db.mgoSess, coll}
}

type collection struct {
	sess *mgo.Session
	coll *mgo.Collection
}

func (c *collection) Find(criteria interface{}) godb.Query {
	sess, coll := c.refresh()
	mgoQuery := coll.Find(criteria)

	return &query{sess, mgoQuery}
}

func (c *collection) FindID(id interface{}) godb.Query {
	sess, coll := c.refresh()
	mgoQuery := coll.FindId(id)

	return &query{sess, mgoQuery}
}

func (c *collection) Insert(doc interface{}) error {
	sess, coll := c.refresh()
	defer sess.Close()

	return coll.Insert(doc)
}

func (c *collection) Update(selector interface{}, update interface{}) error {
	sess, coll := c.refresh()
	defer sess.Close()

	return coll.Update(selector, update)
}

func (c *collection) Upsert(selector interface{}, update interface{}) (*godb.ChangeInfo, error) {
	sess, coll := c.refresh()
	defer sess.Close()

	info, err := coll.Upsert(selector, update)
	return adaptChangeInfo(info), err
}

func (c *collection) Remove(selector interface{}) error {
	sess, coll := c.refresh()
	defer sess.Close()

	return coll.Remove(selector)
}

func (c *collection) RemoveID(id interface{}) error {
	sess, coll := c.refresh()
	defer sess.Close()

	return coll.RemoveId(id)
}

func (c *collection) RemoveAll(selector interface{}) (*godb.ChangeInfo, error) {
	sess, coll := c.refresh()
	defer sess.Close()

	info, err := coll.RemoveAll(selector)
	return adaptChangeInfo(info), err
}

func (c *collection) Bulk() godb.Bulk {
	sess, coll := c.refresh()
	mgoBulk := coll.Bulk()

	return &bulk{sess, mgoBulk}
}

func (c *collection) refresh() (*mgo.Session, *mgo.Collection) {
	sess := c.sess.Copy()
	coll := c.coll.With(sess)

	return sess, coll
}

type query struct {
	sess  *mgo.Session
	query *mgo.Query
}

func (q *query) All(result interface{}) error {
	defer q.sess.Close()
	return q.query.All(result)
}

func (q *query) Apply(change godb.Change, result interface{}) (*godb.ChangeInfo, error) {
	defer q.sess.Close()

	mgoChange := mgo.Change{Update: change.Update, Upsert: change.Upsert, Remove: change.Remove, ReturnNew: change.ReturnNew}
	info, err := q.query.Apply(mgoChange, result)

	return adaptChangeInfo(info), err
}

func (q *query) One(result interface{}) error {
	defer q.sess.Close()
	return q.query.One(result)
}

func (q *query) Iter() godb.Iter {
	return &iter{q.sess, q.query.Iter()}
}

func (q *query) Count() (int, error) {
	defer q.sess.Close()
	return q.query.Count()
}

func (q *query) Limit(n int) godb.Query {
	q.query = q.query.Limit(n)
	return q
}

func (q *query) Sort(fields ...string) godb.Query {
	q.query = q.query.Sort(fields...)
	return q
}

func (q *query) Select(selector interface{}) godb.Query {
	q.query = q.query.Select(selector)
	return q
}

type iter struct {
	sess *mgo.Session
	iter *mgo.Iter
}

func (i *iter) Next(result interface{}) bool {
	hasNext := i.iter.Next(result)

	if !hasNext {
		i.sess.Close()
	}

	return hasNext
}

type bulk struct {
	sess *mgo.Session
	bulk *mgo.Bulk
}

func (b *bulk) Run() (*godb.BulkResult, error) {
	defer b.sess.Close()

	r, err := b.bulk.Run()

	return adaptBulkResult(r), err
}

func (b *bulk) Update(pairs ...interface{}) {
	b.bulk.Update(pairs...)
}

func adaptChangeInfo(info *mgo.ChangeInfo) *godb.ChangeInfo {
	if info == nil {
		return nil
	}

	return &godb.ChangeInfo{Updated: info.Updated, Removed: info.Removed, UpsertedID: info.UpsertedId}
}

func adaptBulkResult(r *mgo.BulkResult) *godb.BulkResult {
	if r == nil {
		return nil
	}

	return &godb.BulkResult{Matched: r.Matched, Modified: r.Modified}
}
