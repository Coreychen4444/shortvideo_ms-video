package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
	"gorm.io/gorm"
)

// 获取用户视频列表
func (s *VideoServer) GetUserVideoList(ctx context.Context, req *pb.UserVideoListRequest) (*pb.UserVideoListResponse, error) {
	videos, err := s.r.GetVideoListByUserId(req.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("该用户不存在")
		}
		return nil, fmt.Errorf("获取视频失败")
	}
	return &pb.UserVideoListResponse{VideoList: videos}, nil
}
