package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"high-concurrency-api/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DataDAO struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewDataDAO(db *gorm.DB, redis *redis.Client) *DataDAO {
	return &DataDAO{
		db:    db,
		redis: redis,
	}
}

func (d *DataDAO) Create(ctx context.Context, data *models.Data) error {
	// 使用事务确保数据一致性
	tx := d.db.Begin()
	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 将数据存入Redis缓存
	jsonData, _ := json.Marshal(data)
	key := fmt.Sprintf("data:%s", data.ID)
	d.redis.Set(ctx, key, jsonData, 30*time.Minute)

	return nil
}

func (d *DataDAO) Update(ctx context.Context, id string, data *models.Data) error {
	// 使用分布式锁确保并发安全
	lockKey := fmt.Sprintf("lock:data:%s", id)
	lock := d.redis.SetNX(ctx, lockKey, "1", 10*time.Second)
	if !lock.Val() {
		return fmt.Errorf("resource is locked")
	}
	defer d.redis.Del(ctx, lockKey)

	tx := d.db.Begin()
	if err := tx.Model(&models.Data{}).Where("id = ?", id).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 更新缓存
	jsonData, _ := json.Marshal(data)
	key := fmt.Sprintf("data:%s", id)
	d.redis.Set(ctx, key, jsonData, 30*time.Minute)

	return nil
}

func (d *DataDAO) Delete(ctx context.Context, id string) error {
	// 使用分布式锁确保并发安全
	lockKey := fmt.Sprintf("lock:data:%s", id)
	lock := d.redis.SetNX(ctx, lockKey, "1", 10*time.Second)
	if !lock.Val() {
		return fmt.Errorf("resource is locked")
	}
	defer d.redis.Del(ctx, lockKey)

	tx := d.db.Begin()
	if err := tx.Delete(&models.Data{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 删除缓存
	key := fmt.Sprintf("data:%s", id)
	d.redis.Del(ctx, key)

	return nil
}

func (d *DataDAO) Get(ctx context.Context, id string) (*models.Data, error) {
	var data models.Data

	// 先从Redis获取
	key := fmt.Sprintf("data:%s", id)
	val, err := d.redis.Get(ctx, key).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &data)
		return &data, nil
	}

	// Redis没有，从数据库获取
	if err := d.db.First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// 写入Redis
	jsonData, _ := json.Marshal(data)
	d.redis.Set(ctx, key, jsonData, 30*time.Minute)

	return &data, nil
} 