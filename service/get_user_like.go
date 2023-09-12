package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

// 获取用户点赞的视频列表
func (s *VideoServer) GetUserLike(ctx context.Context, req *pb.UserLikeRequest) (*pb.UserLikeResponse, error) {
	// 获取用户点赞的视频列表
	videos, err := s.r.GetUserLike(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("获取列表失败")
	}
	return &pb.UserLikeResponse{VideoList: videos}, nil
}
