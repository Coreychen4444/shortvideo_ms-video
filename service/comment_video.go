package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-video/model"
)

// 评论视频
func (s *VideoServer) CommentVideo(ctx context.Context, req *pb.CommentVideoRequest) (*pb.CommentVideoResponse, error) {
	id := req.GetId()
	content := req.GetContent()
	token_user_id := req.GetTokenUserId()
	// 创建评论
	comment := &model.Comment{
		Video_id:   int64(id),
		UserID:     token_user_id,
		Content:    content,
		CreateDate: time.Now().Format("01-02"),
	}
	com, err := s.r.CreateComment(comment)
	if err != nil {
		if err.Error() == "评论发表成功, 但返回评论信息失败" {
			return nil, err
		}
		return nil, fmt.Errorf("评论失败")
	}
	return &pb.CommentVideoResponse{Comment: com}, nil
}
