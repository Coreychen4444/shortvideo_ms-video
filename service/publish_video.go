package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-video/model"
)

// 发布视频
func (s *VideoServer) PublishVideo(ctx context.Context, req *pb.PublishVideoRequest) (*pb.PublishVideoResponse, error) {
	content := req.GetContent()
	title := req.GetTitle()
	token_user_id := req.GetTokenUserId()
	// 生成唯一文件名
	uniqueFileName := fmt.Sprintf("%d.mp4", time.Now().UnixNano())
	// 保存视频文件
	// 保存视频文件到本地
	if _, err := os.Stat("./public/videofile"); os.IsNotExist(err) {
		err = os.MkdirAll("./public/videofile", 0755)
		if err != nil {
			return nil, err
		}
	}
	video_path := "./public/videofile/" + uniqueFileName
	if err := os.WriteFile(video_path, content, 0644); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// 制定视频封面
	uniqueNameWithoutExt := uniqueFileName[0 : len(uniqueFileName)-len(filepath.Ext(uniqueFileName))]
	// 保存封面图片到本地
	// 创建封面图片存放路径
	if _, err := os.Stat("./public/cover"); os.IsNotExist(err) {
		err = os.MkdirAll("./public/cover", 0755)
		if err != nil {
			return nil, err
		}
	}
	coverPath := "./public/cover/" + uniqueNameWithoutExt + ".jpg"
	// 生成封面图片
	err := generateCoverFromVideo(video_path, coverPath)
	if err != nil {
		return nil, fmt.Errorf("生成封面失败")
	}
	// 创建视频
	// https://storage.googleapis.com/<bucket-name>/<object-path>
	video := &model.Video{
		Title:  title,
		UserID: token_user_id,
		/* 		PlayURL:  fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName),
		   		CoverURL: fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectCoverName), */
		PlayURL:     fmt.Sprintf("http://192.168.3.82:8080/public/videofile/%s", uniqueFileName),
		CoverURL:    fmt.Sprintf("http://192.168.3.82:8080/public/cover/%s.jpg", uniqueNameWithoutExt),
		PublishedAt: time.Now().UTC().Unix(),
	}
	err = s.r.CreateVideo(video)
	if err != nil {
		return nil, fmt.Errorf("发布视频失败")
	}
	return &pb.PublishVideoResponse{}, nil
}

// 生成视频封面
func generateCoverFromVideo(videoPath string, coverPath string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath, // 输入文件路径
		"-ss", "00:00:01", // 开始时间，这里设置为视频的第1秒
		"-vframes", "1", // 只输出1帧图片
		"-f", "image2", // 输出格式
		coverPath, // 输出文件路径
	)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate cover: %w", err)
	}
	return nil
}
