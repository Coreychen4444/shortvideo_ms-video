package repo

import (
	"fmt"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-video/model"
	"gorm.io/gorm"
)

// 获取视频列表

func (r *DbRepository) GetVideoList(latest_time int64) ([]*pb.Video, int64, error) {
	var videos []*pb.Video
	err := r.db.Model(&model.Video{}).Preload("Author").Where("published_at < ?", latest_time).Order("published_at desc").Limit(10).Find(&videos).Error
	if err != nil {
		return nil, -1, err
	}
	var nextTime int64
	if len(videos) == 0 {
		nextTime = -1
		return videos, nextTime, nil
	}
	err = r.db.Model(&model.Video{}).Where("id = ?", videos[len(videos)-1].Id).Pluck("published_at", &nextTime).Error
	if err != nil {
		return nil, -1, err
	}
	return videos, nextTime, nil
}

// 创建视频
func (r *DbRepository) CreateVideo(video *model.Video) error {
	tx := r.db.Begin()
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		txErr = err
		return err
	}
	if err := r.db.Create(video).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户投稿数
	if err := r.db.Model(&model.User{}).Where("id = ?", video.UserID).Update("work_count", gorm.Expr("work_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 获取某一用户投稿的视频列表
func (r *DbRepository) GetVideoListByUserId(user_id int64) ([]*pb.Video, error) {
	_, err := r.GetUserById(user_id)
	if err != nil {
		return nil, err
	}
	var videos []*pb.Video
	err = r.db.Model(model.Video{}).Preload("Author").Where("user_id = ?", user_id).Order("published_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 根据用户id获取用户
func (r *DbRepository) GetUserById(id int64) (*pb.User, error) {
	var user pb.User
	err := r.db.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// 点赞
func (r *DbRepository) LikeVideo(user_id int64, video_id int64) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	var videoLike model.VideoLike
	videoLike.UserID = user_id
	videoLike.VideoID = video_id
	if err := r.db.Create(&videoLike).Error; err != nil {
		txErr = err
		return err
	}
	// 更新视频点赞数
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户点赞数
	if err := r.db.Model(&model.User{}).Where("id = ?", user_id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户获赞数
	var authorId int64
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Pluck("user_id", &authorId).Error; err != nil {
		txErr = err
		return err
	}
	if err := r.db.Model(&model.User{}).Where("id = ?", authorId).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 取消点赞
func (r *DbRepository) DislikeVideo(user_id int64, video_id int64) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	var videoLike model.VideoLike
	if err := r.db.Where("user_id = ? and video_id = ?", user_id, video_id).Delete(&videoLike).Error; err != nil {
		txErr = err
		return err
	}
	// 更新视频点赞数
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户点赞数
	if err := r.db.Model(&model.User{}).Where("id = ?", user_id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	// 更新用户获赞数
	var authorId int64
	if err := r.db.Model(&model.Video{}).Where("id = ?", video_id).Pluck("user_id", &authorId).Error; err != nil {
		txErr = err
		return err
	}
	if err := r.db.Model(&model.User{}).Where("id = ?", authorId).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 获取用户点赞的视频列表
func (r *DbRepository) GetUserLike(user_id int64) ([]*pb.Video, error) {
	videoId, err := r.GetUserLikeId(user_id)
	if err != nil {
		return nil, err
	}
	var videos []*pb.Video
	err = r.db.Model(&model.Video{}).Preload("Author").Where("id in (?)", videoId).Order("published_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 获取用户喜欢的视频ID列表
func (r *DbRepository) GetUserLikeId(user_id int64) ([]int64, error) {
	var videoId []int64
	err := r.db.Model(&model.VideoLike{}).Where("user_id = ?", user_id).Pluck("video_id", &videoId).Error
	if err != nil {
		return nil, err
	}
	return videoId, nil
}

// 发布评论
func (r *DbRepository) CreateComment(comment *model.Comment) (*pb.Comment, error) {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	if err := r.db.Create(comment).Error; err != nil {
		txErr = err
		return nil, err
	}
	// 更新视频评论数
	if err := r.db.Model(&model.Video{}).Where("id = ?", comment.Video_id).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		txErr = err
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	var com pb.Comment

	id := comment.ID
	err := r.db.Model(&model.Comment{}).Preload("User").Where("id = ?", id).First(&com).Error
	if err != nil {
		return nil, fmt.Errorf("评论发表成功, 但返回评论信息失败")
	}
	return &com, nil
}

// 删除评论
func (r *DbRepository) DeleteComment(comment_id int64) error {
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	var txErr error
	defer func() {
		if txErr != nil || tx.Commit().Error != nil {
			tx.Rollback()
		}
	}()
	var comment model.Comment
	if err := r.db.Where("id = ?", comment_id).Delete(&comment).Error; err != nil {
		txErr = err
		return err
	}
	// 更新视频评论数
	if err := r.db.Model(&model.Video{}).Where("id = ?", comment.Video_id).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
		txErr = err
		return err
	}
	return tx.Commit().Error
}

// 获取视频评论列表
func (r *DbRepository) GetCommentList(video_id int64) ([]*pb.Comment, error) {
	var comments []*pb.Comment
	err := r.db.Model(&model.Comment{}).Preload("User").Where("video_id = ?", video_id).Order("create_date desc").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
