package utils

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type Snowflake struct {
	mutex       sync.Mutex
	epoch       int64
	machineID   int64
	sequence    int64
	lastTime    int64
	maxSequence int64
}

var lock = &sync.Mutex{}
var instance *Snowflake

func GetSnowflakeInstance() *Snowflake {

	if instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if instance == nil {
			instance = &Snowflake{
				epoch:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
				machineID:   1,
				sequence:    0,
				lastTime:    -1,
				maxSequence: 4095,
			}
		}
	}

	return instance
}

func (s *Snowflake) GenerateId() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixMilli()
	if now < s.lastTime {
		return "", errors.New("clock moved backwards")
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & s.maxSequence
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	id := ((now - s.epoch) << 22) | (s.machineID << 12) | s.sequence
	idStr := strconv.FormatInt(id, 10)
	return idStr, nil
}
