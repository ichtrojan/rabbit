package rabbit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Config struct {
	Conn    *redis.Client
	AppName string
	Job     string
	Queue   string
	Delay   uint64
}

type Param struct {
	Type  string
	Name  string
	Value string
}

func (c Config) Dispatch(params ...Param) error {
	jobId := uuid.New().String()

	pushedAt := time.Now().Unix()

	if c.AppName == "" {
		return errors.New("AppName is required in config")
	}

	if c.Queue == "" {
		return errors.New("queue is required in config")
	}

	if c.Job == "" {
		return errors.New("job is required in config")
	}

	appName := strings.ToLower(c.AppName)

	payload, err := json.Marshal(map[string]interface{}{
		"uuid":          jobId,
		"displayName":   c.Job,
		"job":           "Illuminate\\Queue\\CallQueuedHandler@call",
		"maxTries":      nil,
		"maxExceptions": nil,
		"failOnTimeout": false,
		"backoff":       nil,
		"timeout":       nil,
		"retryUntil":    nil,
		"data": map[string]string{
			"commandName": c.Job,
			"command":     fmt.Sprintf("O:%d:\"%s\":%d:{%ss:3:\"job\";N;s:10:\"connection\";N;s:5:\"queue\";s:7:\"%s\";s:15:\"chainConnection\";N;s:10:\"chainQueue\";N;s:19:\"chainCatchCallbacks\";N;s:5:\"delay\";%ss:11:\"afterCommit\";N;s:10:\"middleware\";a:0:{}s:7:\"chained\";a:0:{}}", len(c.Job), c.Job, 10+len(params), c.generateProperty(params...), c.Queue, c.generateDelay()),
		},
		"id":       jobId,
		"attempts": 0,
		"type":     "job",
		"tags":     []interface{}{},
		"silenced": false,
		"pushedAt": pushedAt,
	})

	if err != nil {
		return err
	}

	if c.Delay == 0 {
		err = c.Conn.RPush(fmt.Sprintf("queues:%s", c.Queue), string(payload)).Err()

		if err != nil {
			return err
		}

		err = c.Conn.RPush(fmt.Sprintf("queues:%s:notify", c.Queue), 1).Err()

		if err != nil {
			return err
		}
	} else {
		delayUnixTime := time.Now().Add(time.Second * time.Duration(c.Delay)).Unix()

		err = c.Conn.ZAdd(fmt.Sprintf("queues:%s:delayed", c.Queue), redis.Z{
			Score:  float64(delayUnixTime),
			Member: string(payload),
		}).Err()

		if err != nil {
			return err
		}
	}

	hashItems := map[string]interface{}{
		"created_at": pushedAt,
		"connection": "redis",
		"updated_at": pushedAt,
		"name":       c.Job,
		"id":         jobId,
		"queue":      c.Queue,
		"payload":    string(payload),
		"status":     "pending",
	}

	for field, hashItem := range hashItems {
		err = c.Conn.HSet(fmt.Sprintf("%s_horizon:%s", c.AppName, jobId), field, hashItem).Err()

		if err != nil {
			return err
		}

		err = c.Conn.Expire(fmt.Sprintf("%s_horizon:%s", c.AppName, jobId), time.Hour).Err()

		if err != nil {
			return err
		}
	}

	err = c.Conn.ZAdd(fmt.Sprintf("%s_horizon:pending_jobs", appName), redis.Z{
		Score:  float64(-pushedAt),
		Member: jobId,
	}).Err()

	if err != nil {
		return err
	}

	err = c.Conn.ZAdd(fmt.Sprintf("%s_horizon:recent_jobs", appName), redis.Z{
		Score:  float64(-pushedAt),
		Member: jobId,
	}).Err()

	if err != nil {
		return err
	}

	return nil
}

func (c Config) generateProperty(params ...Param) string {
	var generatedParams string

	for _, param := range params {
		switch param.Type {
		case "private":
			generatedParams = generatedParams + fmt.Sprintf("s:%d:\"\u0000%s\u0000%s\";s:%d:\"%s\";", len(c.Job)+2, c.Job, param.Name, len(param.Value), param.Value)
			break
		case "protected":
			generatedParams = generatedParams + fmt.Sprintf("s:%d:\"\u0000*\u0000%s\";s:%d:\"%s\";", len(param.Name)+3, param.Name, len(param.Value), param.Value)
			break
		case "public":
			generatedParams = generatedParams + fmt.Sprintf("s:%d:\"%s\";s:%d:\"%s\";", len(param.Name), param.Name, len(param.Value), param.Value)
			break
		default:
			generatedParams = generatedParams + fmt.Sprintf("s:%d:\"%s\";s:%d:\"%s\";", len(param.Name), param.Name, len(param.Value), param.Value)
			break
		}
	}

	return generatedParams
}

func (c Config) generateDelay() string {
	if c.Delay == 0 {
		return "N;"
	}

	return fmt.Sprintf("i:%d;", c.Delay)
}
