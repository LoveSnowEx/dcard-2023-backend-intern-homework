package pagelistServer

import (
	"context"
	"fmt"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/pagelist"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/proto/pb"
	"github.com/google/uuid"
)

type Server struct {
	pb.UnimplementedPageListServiceServer
}

var dbConn *db.DB

func init() {
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}
	dbConn = conn
}

func (s *Server) New(context.Context, *pb.Empty) (*pb.PageList, error) {
	pl := pagelist.New()
	err := dbConn.CreatePageList(pl)
	if err != nil {
		return nil, fmt.Errorf("failed to create: %v", err)
	}
	return &pb.PageList{Key: pl.Key.String()}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.Empty, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	err = dbConn.DeletePageList(u)
	if err != nil {
		return nil, fmt.Errorf("failed to delete: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) Begin(ctx context.Context, req *pb.BeginRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	end, err := dbConn.GetPageListEnd(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get begin: %v", err)
	}
	return &pb.PageIterator{Key: end.NextKey.String()}, nil
}

func (s *Server) End(ctx context.Context, req *pb.EndRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	end, err := dbConn.GetPageListEnd(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get end: %v", err)
	}
	return &pb.PageIterator{Key: end.Key.String()}, nil
}

func (s *Server) Next(ctx context.Context, req *pb.NextRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.IterKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.IterKey)
	}
	next, err := dbConn.GetPageNodeNext(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get next: %v", err)
	}
	return &pb.PageIterator{Key: next.Key.String()}, nil
}

func (s *Server) Prev(ctx context.Context, req *pb.PrevRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.IterKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.IterKey)
	}
	prev, err := dbConn.GetPageNodePrev(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get prev: %v", err)
	}
	return &pb.PageIterator{Key: prev.Key.String()}, nil
}

func (s *Server) Clear(ctx context.Context, pl *pb.PageList) (*pb.Empty, error) {
	u, err := uuid.Parse(pl.Key)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", pl.Key)
	}
	err = dbConn.DeletePageList(u)
	if err != nil {
		return nil, fmt.Errorf("failed to clear: %v", err)
	}
	return &pb.Empty{}, err
}

func (s *Server) Insert(ctx context.Context, req *pb.InsertRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.IterKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.IterKey)
	}
	node, err := dbConn.InsertPageNode(u, uint(req.PageId))
	if err != nil {
		return nil, fmt.Errorf("failed to insert: %v", err)
	}
	return &pb.PageIterator{Key: node.Key.String(), PageId: uint32(node.PageID)}, nil
}

func (s *Server) Erase(ctx context.Context, req *pb.EraseRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.IterKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.IterKey)
	}
	node, err := dbConn.ErasePageNode(u)
	if err != nil {
		return nil, fmt.Errorf("failed to erase: %v", err)
	}
	return &pb.PageIterator{Key: node.Key.String(), PageId: uint32(node.PageID)}, nil
}

func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.IterKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.IterKey)
	}
	node, err := dbConn.SetPageNode(u, uint(req.PageId))
	if err != nil {
		return nil, fmt.Errorf("failed to set: %v", err)
	}
	return &pb.PageIterator{Key: node.Key.String(), PageId: uint32(node.PageID)}, nil
}

func (s *Server) PushBack(ctx context.Context, req *pb.PushRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	node, err := dbConn.PushBackPageList(u, uint(req.PageId))
	if err != nil {
		return nil, fmt.Errorf("failed to push back: %v", err)
	}
	return &pb.PageIterator{Key: node.Key.String(), PageId: uint32(node.PageID)}, nil
}

func (s *Server) PopBack(ctx context.Context, req *pb.PopRequest) (*pb.Empty, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	err = dbConn.PopBackPageList(u)
	if err != nil {
		return nil, fmt.Errorf("failed to pop back: %v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Server) PushFront(ctx context.Context, req *pb.PushRequest) (*pb.PageIterator, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	node, err := dbConn.PushFrontPageList(u, uint(req.PageId))
	if err != nil {
		return nil, fmt.Errorf("failed to push front: %v", err)
	}
	return &pb.PageIterator{Key: node.Key.String(), PageId: uint32(node.PageID)}, nil
}

func (s *Server) PopFront(ctx context.Context, req *pb.PopRequest) (*pb.Empty, error) {
	u, err := uuid.Parse(req.ListKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %s", req.ListKey)
	}
	err = dbConn.PopFrontPageList(u)
	if err != nil {
		return nil, fmt.Errorf("failed to pop front: %v", err)
	}
	return &pb.Empty{}, nil
}
