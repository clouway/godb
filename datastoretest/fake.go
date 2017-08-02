package datastoretest

import (
	"github.com/clouway/godb"
	"github.com/stretchr/testify/mock"
)

type FakeDatabase struct {
	FakeBulk       *FakeBulk
	FakeQuery      *FakeQuery
	FakeCollection *FakeCollection
	mock.Mock
}

func NewFakeDatabase() *FakeDatabase {
	bulk := new(FakeBulk)
	query := new(FakeQuery)

	return &FakeDatabase{
		FakeBulk:  bulk,
		FakeQuery: query,
		FakeCollection: &FakeCollection{
			FakeBulk:  bulk,
			FakeQuery: query,
		},
	}
}

func (d *FakeDatabase) Close() {}

func (d *FakeDatabase) Collection(name string) godb.Collection {
	return d.FakeCollection
}

func (d *FakeDatabase) Collections() ([]godb.Collection, error) {
	return nil, nil
}

func (d *FakeDatabase) Indexer(cname string) godb.Indexer {
	return nil
}

func (d *FakeDatabase) Ping() error {
	return d.Called().Error(0)
}

func (d *FakeDatabase) DropDatabase() error {
	return d.Called().Error(0)
}

type FakeCollection struct {
	FakeBulk  *FakeBulk
	FakeQuery *FakeQuery
	mock.Mock
}

func (c *FakeCollection) Find(query interface{}) godb.Query {
	return c.FakeQuery
}

func (c *FakeCollection) FindID(id interface{}) godb.Query {

	return c.FakeQuery
}

func (c *FakeCollection) Insert(doc interface{}) error {
	return c.Called().Error(0)
}

func (c *FakeCollection) Update(selector interface{}, update interface{}) error {
	return c.Called().Error(0)
}

func (c *FakeCollection) UpdateAll(selector interface{}, update interface{}) (*godb.ChangeInfo, error) {
	args := c.Called()
	return args.Get(0).(*godb.ChangeInfo), args.Error(1)
}

func (c *FakeCollection) Upsert(selector interface{}, update interface{}) (*godb.ChangeInfo, error) {
	args := c.Called()
	return args.Get(0).(*godb.ChangeInfo), args.Error(1)
}

func (c *FakeCollection) Remove(selector interface{}) error {
	return c.Called().Error(0)
}

func (c *FakeCollection) RemoveID(id interface{}) error {
	return c.Called().Error(0)
}

func (c *FakeCollection) RemoveAll(selector interface{}) (*godb.ChangeInfo, error) {
	args := c.Called()
	return args.Get(0).(*godb.ChangeInfo), args.Error(1)
}

func (c *FakeCollection) Bulk() godb.Bulk {
	return c.FakeBulk
}

func (c *FakeCollection) Clean() error {
	args := c.Called()
	return args.Error(0)
}

type FakeQuery struct {
	mock.Mock
}

func (q *FakeQuery) All(result interface{}) error {
	return q.Called().Error(0)
}

func (q *FakeQuery) Apply(change godb.Change, result interface{}) (*godb.ChangeInfo, error) {
	args := q.Called()
	return args.Get(0).(*godb.ChangeInfo), args.Error(1)
}

func (q *FakeQuery) One(result interface{}) error {
	return q.Called().Error(0)
}

func (q *FakeQuery) Iter() godb.Iter {
	return q.Called().Get(0).(godb.Iter)
}

func (q *FakeQuery) Count() (int, error) {
	args := q.Called()
	return args.Get(0).(int), args.Error(1)
}

func (q *FakeQuery) Limit(n int) godb.Query {
	return q
}

func (q *FakeQuery) Sort(fields ...string) godb.Query {
	return q
}

func (q *FakeQuery) Select(selector interface{}) godb.Query {
	return q
}

type FakeBulk struct {
	mock.Mock
}

func (b *FakeBulk) Run() (*godb.BulkResult, error) {
	args := b.Called()
	return args.Get(0).(*godb.BulkResult), args.Error(1)
}

func (b *FakeBulk) Insert(docs ...interface{}) {}

func (b *FakeBulk) Update(pairs ...interface{}) {}

func (b *FakeBulk) Upsert(pairs ...interface{}) {}
