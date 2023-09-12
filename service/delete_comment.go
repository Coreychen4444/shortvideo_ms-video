package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

// 删除评论
func (s *VideoServer) DeleteComment(ctx context.Context, req *pb.CommentVideoRequest) (*pb.CommentVideoResponse, error) {
	// 删除评论
	err := s.r.DeleteComment(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("删除评论失败")
	}
	return &pb.CommentVideoResponse{}, nil
}
