package client

import (
	"container/list"
	"context"
	"log"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/internal/service"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/proto/pb"
	"google.golang.org/grpc"
)

var (
	client pb.PageListServiceClient
)

func TestMain(m *testing.M) {
	_, err := db.MockConnet()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.MockClose()

	go func() {
		err := service.RunGrpc(":50051")
		if err != nil {
			log.Fatalf("failed to run grpc: %v", err)
		}
	}()

	conn, err := grpc.Dial(":50051", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client = pb.NewPageListServiceClient(conn)

	m.Run()
}

func TestNew(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestEnd(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	it, err := client.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, it)
}

func TestBegin(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	beginIt, err := client.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, beginIt)

	endIt, err := client.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, endIt)

	require.Equal(t, beginIt.Key, endIt.Key)
}

func CompareList(t *testing.T, l *list.List, listKey string) {
	it, err := client.Begin(context.Background(), &pb.BeginRequest{ListKey: listKey})
	require.NoError(t, err)
	require.NotNil(t, it)

	endIt, err := client.End(context.Background(), &pb.EndRequest{ListKey: listKey})
	require.NoError(t, err)
	require.NotNil(t, endIt)

	for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
		require.Equal(t, e.Value, it.PageId)
		it, err = client.Next(context.Background(), &pb.NextRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	require.Equal(t, endIt.Key, it.Key)
}

func TestPushBack0(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	CompareList(t, l, res.Key)
}

func TestPushBack1(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	val := uint32(1)

	l.PushBack(val)

	it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
	require.NoError(t, err)
	require.NotNil(t, it)

	CompareList(t, l, res.Key)
}

func TestPushBack100(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestPushBackRandom100(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := rand.Uint32()/1000 + 1

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestPushFront0(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	CompareList(t, l, res.Key)
}

func TestPushFront1(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	val := uint32(1)

	l.PushFront(val)

	it, err := client.PushFront(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
	require.NoError(t, err)
	require.NotNil(t, it)

	CompareList(t, l, res.Key)
}

func TestPushFront100(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushFront(val)

		it, err := client.PushFront(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestPushFrontRandom100(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := rand.Uint32()/1000 + 1

		l.PushFront(val)

		it, err := client.PushFront(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestPrev(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	it, err := client.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, it)

	endIt, err := client.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, endIt)

	for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
		require.Equal(t, e.Value, it.PageId)
		oldIt := it
		it, err = client.Next(context.Background(), &pb.NextRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, it)

		prevIt, err := client.Prev(context.Background(), &pb.PrevRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, prevIt)

		require.Equal(t, oldIt.Key, prevIt.Key)
	}

	require.Equal(t, endIt.Key, it.Key)
}

func TestDelete(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	_, err = client.Delete(context.Background(), &pb.DeleteRequest{ListKey: res.Key})
	require.NoError(t, err)

	_, err = client.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.Error(t, err)

	_, err = client.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.Error(t, err)
}

func TestClear(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	_, err = client.Clear(context.Background(), &pb.ClearRequest{ListKey: res.Key})
	require.NoError(t, err)

	it, err := client.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, it)

	endIt, err := client.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, endIt)

	require.Equal(t, endIt.Key, it.Key)
}

func RandIter(t *testing.T, l *list.List, key string) (*list.Element, *pb.PageIterator) {
	it, err := client.Begin(context.Background(), &pb.BeginRequest{ListKey: key})
	require.NoError(t, err)
	require.NotNil(t, it)

	r := rand.Intn(l.Len())
	e := l.Front()
	for i := 0; i < r; i, e = i+1, e.Next() {
		it, err = client.Next(context.Background(), &pb.NextRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, it)
	}
	return e, it
}

func TestInsert(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		e, it := RandIter(t, l, res.Key)

		l.InsertBefore(val, e)

		it, err = client.Insert(context.Background(), &pb.InsertRequest{IterKey: it.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestErase(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	for i := 0; i < 100; i++ {
		e, it := RandIter(t, l, res.Key)

		l.Remove(e)

		it, err = client.Erase(context.Background(), &pb.EraseRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestSet(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	for i := 0; i < 100; i++ {
		e, it := RandIter(t, l, res.Key)

		val := uint32(i + 1)

		e.Value = val

		it, err = client.Set(context.Background(), &pb.SetRequest{IterKey: it.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestPopBack(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	for i := 0; i < 100; i++ {
		l.Remove(l.Back())

		it, err := client.PopBack(context.Background(), &pb.PopRequest{ListKey: res.Key})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestPopFront(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	for i := 0; i < 100; i++ {
		l.Remove(l.Front())

		it, err := client.PopFront(context.Background(), &pb.PopRequest{ListKey: res.Key})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestClone(t *testing.T) {
	res, err := client.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := client.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	res, err = client.Clone(context.Background(), &pb.CloneRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, res)

	CompareList(t, l, res.Key)
}
