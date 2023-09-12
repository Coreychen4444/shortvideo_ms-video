package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// 获取视频列表
func (s *VideoServer) GetVideoFlow(ctx context.Context, req *pb.VideoFlowRequest) (*pb.VideoFlowResponse, error) {
	// 获取最新时间
	latestTime := req.GetLatestTime()
	videos, next_time, err := s.r.GetVideoList(latestTime)
	if err != nil {
		return nil, fmt.Errorf("获取视频失败")
	}
	token := req.GetToken()
	user_id, err := VerifyToken(token)
	if err != nil {
		return &pb.VideoFlowResponse{VideoList: videos, NextTime: next_time}, nil
	}
	// 返回点赞个性化信息
	videoId, err := s.r.GetUserLikeId(user_id)
	if err != nil {
		return &pb.VideoFlowResponse{VideoList: videos, NextTime: next_time}, nil
	}
	isLike := make(map[int64]bool)
	for _, id := range videoId {
		isLike[id] = true
	}
	for i := range videos {
		videos[i].IsFavorite = isLike[videos[i].Id]
	}
	return &pb.VideoFlowResponse{VideoList: videos, NextTime: next_time}, nil
}

// 验证token

func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		SECRET_KEY, err := getSecretKey()
		if err != nil {
			return nil, err
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		userId, ok := (*claims)["userId"].(float64) // 注意：因为JSON的解码机制，整数通常会被解码为float64
		if !ok {
			return -1, errors.New("token is invalid")
		}
		return int64(userId), nil
	}

	return -1, errors.New("token is not valid")
}

func getSecretKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	return os.Getenv("SECRET_KEY"), nil
}
