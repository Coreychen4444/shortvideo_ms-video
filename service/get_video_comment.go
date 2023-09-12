package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

// 获取视频评论列表
func (s *VideoServer) GetCommentList(ctx context.Context, req *pb.VideoCommentRequest) (*pb.VideoCommentResponse, error) {
	// 获取视频评论列表
	comments, err := s.r.GetCommentList(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("获取评论列表失败")
	}
	return &pb.VideoCommentResponse{CommentList: comments}, nil
}
