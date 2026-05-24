package svc

import (
	"context"
	"encoding/json"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/config"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type CacheService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewCacheService(c config.Config) *CacheService {
	return &CacheService{
		DB:  utils.InitDB(c.Mysql.DataSource),
		RDB: utils.InitRDB(c.Cache[0].Host, c.Cache[0].Pass),
	}
}

// SetCache 定义缓存读写服务
func (s *CacheService) SetCache(key string, value interface{}, expiration time.Duration) error {
	return s.RDB.Set(context.Background(), key, value, expiration).Err()
}

func (s *CacheService) GetCache(key string) (string, error) {
	return s.RDB.Get(context.Background(), key).Result()
}

func (s *CacheService) GetUserCache(user *models.User, key string) error {
	data, err := s.RDB.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), user)
}

func (s *CacheService) SetUserCache(user *models.User, key string) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.SetCache(key, userData, time.Hour*24)
}

// GetUserByPhone 业务逻辑层
func (s *CacheService) GetUserByPhone(user *models.User, phone string) error {
	cacheKey := "user_phone:" + phone
	if err := s.GetUserCache(user, cacheKey); err == nil {
		return nil
	}

	if err := s.DB.Where("phone = ?", phone).First(user).Error; err != nil {
		return err
	}
	return s.SetUserCache(user, cacheKey)
}

func (s *CacheService) GetUserByIds(users *[]models.User, ids []string) error {
	if len(ids) == 0 {
		*users = []models.User{}
		return nil
	}

	orderedUsers := make([]models.User, 0, len(ids))
	for _, id := range ids {
		cacheKey := "user_id:" + id
		var user models.User
		if err := s.GetUserCache(&user, cacheKey); err == nil {
			orderedUsers = append(orderedUsers, user)
			continue
		}

		err := s.DB.Where("id = ?", id).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return err
		}
		if err := s.SetUserCache(&user, cacheKey); err != nil {
			return err
		}
		orderedUsers = append(orderedUsers, user)
	}
	*users = orderedUsers
	return nil
}

func (s *CacheService) GetUserByName(user *models.User, name string) error {
	cacheKey := "user_name:" + name
	if err := s.GetUserCache(user, cacheKey); err == nil {
		return nil
	}
	if err := s.DB.Where("nickname = ?", name).First(user).Error; err != nil {
		return err
	}
	return s.SetUserCache(user, cacheKey)
}

func (s *CacheService) UpdateUserProfile(userID, nickname string, sex int8, email string) error {
	var u models.User
	if err := s.DB.Where("id = ?", userID).First(&u).Error; err != nil {
		return err
	}

	oldNickname := u.Nickname
	updatedSex := utils.ConvertToInt8(sex)
	updates := map[string]interface{}{
		"nickname": nickname,
		"sex":      updatedSex,
		"email":    email,
	}
	if err := s.DB.Model(&u).Updates(updates).Error; err != nil {
		return err
	}

	u.Nickname = nickname
	u.Sex = updatedSex
	u.Email = email

	ctx := context.Background()
	cacheKeys := []string{"user_id:" + u.ID, "user_phone:" + u.Phone}
	if oldNickname != "" {
		cacheKeys = append(cacheKeys, "user_name:"+oldNickname)
	}
	if u.Nickname != "" && u.Nickname != oldNickname {
		cacheKeys = append(cacheKeys, "user_name:"+u.Nickname)
	}

	for _, cacheKey := range cacheKeys {
		if err := s.RDB.Del(ctx, cacheKey).Err(); err != nil {
			logx.Errorf("failed to invalidate user cache after profile update, userID=%s cacheKey=%s err=%v", userID, cacheKey, err)
		}
	}
	return nil
}

func (s *CacheService) UpdateUserPassword(userID, passwordHash string) error {
	ctx := context.Background()
	var u models.User
	if err := s.DB.Where("id = ?", userID).First(&u).Error; err != nil {
		return err
	}
	if err := s.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", passwordHash).Error; err != nil {
		return err
	}
	cacheKeys := []string{"user_id:" + userID}
	if u.Phone != "" {
		cacheKeys = append(cacheKeys, "user_phone:"+u.Phone)
	}
	for _, cacheKey := range cacheKeys {
		if err := s.RDB.Del(ctx, cacheKey).Err(); err != nil {
			logx.Errorf("failed to invalidate user cache after password update, userID=%s cacheKey=%s err=%v", userID, cacheKey, err)
		}
	}
	return nil
}

// CreateUser 保存
func (s *CacheService) CreateUser(user *models.User) error {
	cacheKey := "user_phone:" + user.Phone
	err := s.DB.Create(user).Error
	if err != nil {
		return err
	}
	s.RDB.Del(context.Background(), cacheKey)
	// 延时双删策略
	time.Sleep(1 * time.Second)
	s.RDB.Del(context.Background(), cacheKey)
	return nil
}
