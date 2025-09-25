// pkg/queue/redis_queue.go
package queue

import (
	"context"
	"encoding/json"
	"llmapisrv/pkg/logger"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(client *redis.Client) *RedisQueue {
	return &RedisQueue{
		client: client,
	}
}

// Push 推送消息到队列
func (q *RedisQueue) Push(queue string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	logger.Infof("RedisQueue Push, body: %v", string(jsonData))

	return q.client.RPush(context.Background(), queue, jsonData).Err()
}

// ProcessQueue 处理队列中的消息
func (q *RedisQueue) ProcessQueue(queue string, handler func([]byte) error) {
	for {
		// 阻塞式获取消息
		result, err := q.client.BLPop(context.Background(), 5*time.Second, queue).Result()
		if err != nil {
			if err != redis.Nil {
				log.Printf("Error getting message from queue %s: %v", queue, err)
			}
			continue
		}
		logger.Infof("ProcessQueue, body data: %v", result)

		if len(result) < 2 {
			continue
		}

		// 处理消息
		data := []byte(result[1])
		if err := handler(data); err != nil {
			log.Printf("Error processing message from queue %s: %v", queue, err)
		}
	}
}

// StartWorker 启动队列处理工作
func (q *RedisQueue) StartWorker(queue string, handler func([]byte) error) {
	go q.ProcessQueue(queue, handler)
}
