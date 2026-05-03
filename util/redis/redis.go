package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

const (
	RedisPrefix                    = "pet:order:"
	RedisPrefixCreateCode          = RedisPrefix + "create_code:"
	RedisPrefixFinishCode          = RedisPrefix + "finish_code:"
	RedisPrefixCreateCodeUserState = RedisPrefixCreateCode + "user_state:"
	RedisPrefixInterviewUpdateInfo = RedisPrefixCreateCode + "interview:update_info:"
)

// Init 初始化Redis连接
func Init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 100, // 连接池大小
	})

	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logrus.Errorf("redis start failed, err: %v", err)
		panic(err)
	}
	return nil
}

// Close 关闭Redis连接
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}

// ==================== 字符串操作 ====================

// Set 设置键值对，带过期时间
func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

// Get 获取字符串值

func Get(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		// 如果是key不存在的错误，返回"0"和nil
		if err == redis.Nil {
			return "0", nil
		}
		// 其他错误直接返回
		return "", err
	}
	return val, nil
}

// GetInt 获取整数值
func GetInt(key string) (int64, error) {
	return rdb.Get(ctx, key).Int64()
}

// GetBool 获取布尔值
func GetBool(key string) (bool, error) {
	return rdb.Get(ctx, key).Bool()
}

// MSet 批量设置键值对
func MSet(pairs map[string]interface{}) error {
	return rdb.MSet(ctx, pairs).Err()
}

// MGet 批量获取值
func MGet(keys ...string) ([]interface{}, error) {
	return rdb.MGet(ctx, keys...).Result()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return rdb.Incr(ctx, key).Result()
}

// IncrBy 按指定步长自增
func IncrBy(key string, value int64) (int64, error) {
	return rdb.IncrBy(ctx, key, value).Result()
}

// ==================== 哈希操作 ====================

// HSet 设置哈希字段值
func HSet(key string, field string, value interface{}) error {
	return rdb.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段值
func HGet(key string, field string) (string, error) {
	return rdb.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段和值
func HGetAll(key string) (map[string]string, error) {
	return rdb.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	return rdb.HDel(ctx, key, fields...).Err()
}

// HExists 检查哈希字段是否存在
func HExists(key string, field string) (bool, error) {
	return rdb.HExists(ctx, key, field).Result()
}

// ==================== 列表操作 ====================

// LPush 从列表左侧插入元素
func LPush(key string, values ...interface{}) error {
	return rdb.LPush(ctx, key, values...).Err()
}

// RPush 从列表右侧插入元素
func RPush(key string, values ...interface{}) error {
	return rdb.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出元素
func LPop(key string) (string, error) {
	return rdb.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出元素
func RPop(key string) (string, error) {
	return rdb.RPop(ctx, key).Result()
}

// LRange 获取列表范围内的元素
func LRange(key string, start, stop int64) ([]string, error) {
	return rdb.LRange(ctx, key, start, stop).Result()
}

// ==================== 集合操作 ====================

// SAdd 向集合添加元素
func SAdd(key string, members ...interface{}) error {
	return rdb.SAdd(ctx, key, members...).Err()
}

// SRem 从集合移除元素
func SRem(key string, members ...interface{}) error {
	return rdb.SRem(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
	return rdb.SMembers(ctx, key).Result()
}

// SIsMember 检查元素是否在集合中
func SIsMember(key string, member interface{}) (bool, error) {
	return rdb.SIsMember(ctx, key, member).Result()
}

// ==================== 有序集合操作 ====================

// ZAdd 向有序集合添加元素
func ZAdd(key string, members ...*redis.Z) error {
	return rdb.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围内的元素
func ZRange(key string, start, stop int64) ([]string, error) {
	return rdb.ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 获取有序集合范围内的元素(逆序)
func ZRevRange(key string, start, stop int64) ([]string, error) {
	return rdb.ZRevRange(ctx, key, start, stop).Result()
}

// ==================== 键操作 ====================

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	n, err := rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// Delete 删除键
func Delete(key string) error {
	return rdb.Del(ctx, key).Err()
}

// Expire 设置键的过期时间
func Expire(key string, expiration time.Duration) error {
	return rdb.Expire(ctx, key, expiration).Err()
}

// TTL 获取键的剩余生存时间
func TTL(key string) (time.Duration, error) {
	return rdb.TTL(ctx, key).Result()
}

// ==================== 高级功能 ====================

// Lock 获取分布式锁
func Lock(key string, value interface{}, expiration time.Duration) (bool, error) {
	return rdb.SetNX(ctx, key, value, expiration).Result()
}

// Unlock 释放分布式锁
func Unlock(key string) error {
	return rdb.Del(ctx, key).Err()
}

// Publish 发布消息到频道
func Publish(channel string, message interface{}) error {
	return rdb.Publish(ctx, channel, message).Err()
}
