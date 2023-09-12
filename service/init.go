package service

import (
	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-video/repo"
)

type VideoServer struct {
	r *repo.DbRepository
	pb.UnimplementedVideoServiceServer
}

func NewVideoServer(r *repo.DbRepository) *VideoServer {
	return &VideoServer{r: r}
}
