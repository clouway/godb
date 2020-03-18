package datastoretest

import (
	"github.com/stretchr/testify/mock"

	"github.com/clouway/godb"
)

type FakeDatabase struct {
	FakeBulk       *FakeBulk
	FakePipe       *FakePipe
	FakeQuery      *FakeQuery
	FakeCollection *FakeCollection
	mock.Mock
}

func NewFakeDatabase() *FakeDatabase {
	bulk := &FakeBulk{}
	pipe := &FakePipe{}

	query := &FakeQuery{
		FakeIter: &FakeIter{},
	}

	return &FakeDatabase{
		FakeBulk:  bulk,
		FakePipe:  pipe,
		FakeQuery: query,
		FakeCollection: &FakeCollection{
			FakeBulk:  bulk,
			FakePipe:  pipe,
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
	FakePipe  *FakePipe
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

func (c *FakeCollection) Pipe(pipeline interface{}) godb.Pipe {
	return c.FakePipe
}

func (c *FakeCollection) Clean() error {
	args := c.Called()
	return args.Error(0)
}

type FakeQuery struct {
	FakeIter *FakeIter
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
	return q.FakeIter
}

func (q *FakeQuery) Count() (int, error) {
	args := q.Called()
	return args.Get(0).(int), args.Error(1)
}

func (q *FakeQuery) Distinct(key string, result interface{}) error {
	args := q.Called(key, result)
	return args.Error(0)
}

func (q *FakeQuery) Limit(n int) godb.Query {
	return q
}

func (q *FakeQuery) Skip(n int) godb.Query {
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

type FakePipe struct {
	mock.Mock
}

func (b *FakePipe) All(result interface{}) error {
	return b.Called().Error(0)
}

func (b *FakePipe) One(result interface{}) error {
	return b.Called().Error(0)
}

type FakeIter struct {
	mock.Mock
}

func (i *FakeIter) Err() error {
	args := i.Called()
	return args.Error(0)
}

func (i *FakeIter) Next(result interface{}) bool {
	args := i.Called(result)
	return args.Bool(0)
}

func (i *FakeIter) Done() bool {
	args := i.Called()
	return args.Bool(0)
}

func (i *FakeIter) Close() error {
	args := i.Called()
	return args.Error(0)
}

func (i *FakeIter) Timeout() bool {
	args := i.Called()
	return args.Bool(0)
}
