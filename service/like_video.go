package service

import (
	"context"
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
)

// 点赞和取消点赞视频
func (s *VideoServer) LikeVideo(ctx context.Context, req *pb.LikeVideoRequest) (*pb.LikeVideoResponse, error) {
	video_id := req.GetId()
	token_user_id := req.GetTokenUserId()
	action_type := req.GetActionType()
	// 根据action_type执行点赞或取消点赞操作
	if action_type == "1" {
		err := s.r.LikeVideo(token_user_id, video_id)
		if err != nil {
			return nil, fmt.Errorf("点赞视频失败")
		}
	} else if action_type == "2" {
		err := s.r.DislikeVideo(token_user_id, video_id)
		if err != nil {
			return nil, fmt.Errorf("取消点赞失败")
		}
	} else {
		return nil, fmt.Errorf("请求参数错误")
	}
	return &pb.LikeVideoResponse{}, nil
}
