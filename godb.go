package godb

type Database interface {
	Close()

	Collection(name string) Collection
}

type Collection interface {
	Find(query interface{}) Query

	FindID(id interface{}) Query

	Insert(doc interface{}) error

	Update(selector interface{}, update interface{}) error

	Upsert(selector interface{}, update interface{}) (*ChangeInfo, error)

	Remove(selector interface{}) error

	RemoveID(id interface{}) error

	RemoveAll(selector interface{}) (*ChangeInfo, error)

	Bulk() Bulk
}

type Query interface {
	All(result interface{}) error

	Apply(change Change, result interface{}) (*ChangeInfo, error)

	One(result interface{}) error

	Iter() Iter

	Count() (int, error)

	Limit(n int) Query

	Sort(fields ...string) Query

	Select(selector interface{}) Query
}

type Iter interface {
	Next(result interface{}) bool
}

type Bulk interface {
	Run() (*BulkResult, error)

	Update(pairs ...interface{})
}

type Config struct {
	Addrs []string

	Database string
}

type Change struct {
	Update    interface{}
	Upsert    bool
	Remove    bool
	ReturnNew bool
}

type ChangeInfo struct {
	Updated    int
	Removed    int
	UpsertedID interface{}
}

type BulkResult struct {
	Matched  int
	Modified int
}
