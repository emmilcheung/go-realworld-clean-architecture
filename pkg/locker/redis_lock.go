package locker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type Locker interface {
	ObtainLock(ctx context.Context, key string) *redislock.Lock
}

type locker struct {
	Redis *redis.Client
	Lock  *redislock.Client
}

var backoff = redislock.LinearBackoff(50 * time.Millisecond)

func LockerInit(redisClient *redis.Client) Locker {
	lockClient := redislock.New(redisClient)
	return &locker{Redis: redisClient, Lock: lockClient}
}

func (l *locker) ObtainLock(ctx context.Context, key string) *redislock.Lock {
	lockCtx, cancel := context.WithDeadline(ctx, time.Now().Add(2*time.Minute))
	defer cancel()

	// Obtain lock with retry + custom deadline
	lock, err := l.Lock.Obtain(lockCtx,
		fmt.Sprintf("realworld:%s", key),
		30*time.Second,
		&redislock.Options{RetryStrategy: backoff})

	if err == redislock.ErrNotObtained {
		fmt.Println("Could not obtain lock!")
	} else if err != nil {
		log.Fatalln(err)
	}

	// Don't forget to defer Release.
	// fmt.Println("I have a lock!")
	return lock
}
