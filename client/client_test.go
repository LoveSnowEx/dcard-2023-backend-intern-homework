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
	c pb.PageListServiceClient
)

func TestMain(m *testing.M) {
	go func() {
		_, err := db.MockConnet()
		if err != nil {
			log.Fatalf("failed to connect db: %v", err)
		}
		defer db.MockClose()
		err = service.RunGrpc(":50051")
		if err != nil {
			log.Fatalf("failed to run grpc: %v", err)
		}
	}()

	conn, err := grpc.Dial(":50051", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	c = pb.NewPageListServiceClient(conn)

	m.Run()
}

func TestNewList(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestListEnd(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	it, err := c.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, it)
}

func TestListBegin(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	itBegin, err := c.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, itBegin)

	itEnd, err := c.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, itEnd)

	require.Equal(t, itBegin.Key, itEnd.Key)
}

func CompareList(t *testing.T, l *list.List, listKey string) {
	it, err := c.Begin(context.Background(), &pb.BeginRequest{ListKey: listKey})
	require.NoError(t, err)
	require.NotNil(t, it)

	itEnd, err := c.End(context.Background(), &pb.EndRequest{ListKey: listKey})
	require.NoError(t, err)
	require.NotNil(t, itEnd)

	for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
		require.Equal(t, e.Value, it.PageId)
		it, err = c.Next(context.Background(), &pb.NextRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	require.Equal(t, itEnd.Key, it.Key)
}

func TestListPushBack0(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	CompareList(t, l, res.Key)
}

func TestListPushBack1(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	val := uint32(1)

	l.PushBack(val)

	it, err := c.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
	require.NoError(t, err)
	require.NotNil(t, it)

	CompareList(t, l, res.Key)
}

func TestListPushBack100(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := c.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestListPushBackRandom100(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := rand.Uint32()/1000 + 1

		l.PushBack(val)

		it, err := c.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestListPushFront0(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	CompareList(t, l, res.Key)
}

func TestListPushFront1(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	val := uint32(1)

	l.PushFront(val)

	it, err := c.PushFront(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
	require.NoError(t, err)
	require.NotNil(t, it)

	CompareList(t, l, res.Key)
}

func TestListPushFront100(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushFront(val)

		it, err := c.PushFront(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestListPushFrontRandom100(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := rand.Uint32()/1000 + 1

		l.PushFront(val)

		it, err := c.PushFront(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	CompareList(t, l, res.Key)
}

func TestListPrev(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})
	require.NoError(t, err)

	l := list.New()

	for i := 0; i < 100; i++ {
		val := uint32(i + 1)

		l.PushBack(val)

		it, err := c.PushBack(context.Background(), &pb.PushRequest{ListKey: res.Key, PageId: val})
		require.NoError(t, err)
		require.NotNil(t, it)
	}

	it, err := c.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, it)

	itEnd, err := c.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.NoError(t, err)
	require.NotNil(t, itEnd)

	for i, e := 0, l.Front(); e != nil; i, e = i+1, e.Next() {
		require.Equal(t, e.Value, it.PageId)
		itOld := it
		it, err = c.Next(context.Background(), &pb.NextRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, it)

		itPrev, err := c.Prev(context.Background(), &pb.PrevRequest{IterKey: it.Key})
		require.NoError(t, err)
		require.NotNil(t, itPrev)

		require.Equal(t, itOld.Key, itPrev.Key)
	}

	require.Equal(t, itEnd.Key, it.Key)
}

func TestDeleteList(t *testing.T) {
	res, err := c.New(context.Background(), &pb.Empty{})

	_, err = c.Delete(context.Background(), &pb.DeleteRequest{ListKey: res.Key})
	require.NoError(t, err)

	_, err = c.Begin(context.Background(), &pb.BeginRequest{ListKey: res.Key})
	require.Error(t, err)

	_, err = c.End(context.Background(), &pb.EndRequest{ListKey: res.Key})
	require.Error(t, err)
}
