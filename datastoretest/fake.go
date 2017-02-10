package datastoretest

import (
	"github.com/clouway/godb"
	"github.com/stretchr/testify/mock"
)

type FakeDatabase struct {
	FakeQuery      *FakeQuery
	FakeCollection *FakeCollection
}

func NewFakeDatabase() *FakeDatabase {
	query := new(FakeQuery)

	return &FakeDatabase{
		FakeQuery:      query,
		FakeCollection: &FakeCollection{FakeQuery: query},
	}
}

func (d *FakeDatabase) Close() {}

func (d *FakeDatabase) Collection(name string) godb.Collection {
	return d.FakeCollection
}

func (d *FakeDatabase) Indexer(cname string) godb.Indexer {
	return nil
}

type FakeCollection struct {
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
	return c.Called().Get(0).(godb.Bulk)
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
