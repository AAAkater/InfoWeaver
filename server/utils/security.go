package utils

import (
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sony/sonyflake"
)

func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	sf     *sonyflake.Sonyflake
	sfOnce sync.Once
)

// GenerateSnowID 生成一个雪花 ID
func GenerateSnowID() (uint64, error) {
	init := func() {
		startTime := time.Date(2006, 4, 11, 0, 0, 0, 0, time.UTC) // 自定义起始时间（epoch）

		machineID := func() (uint16, error) {
			return 1, nil // 生产环境中应确保每台机器 ID 唯一
		}

		sf = sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
			MachineID: machineID,
		})

		if sf == nil {
			panic("failed to create Sonyflake instance")
		}
	}
	sfOnce.Do(init)
	return sf.NextID()
}
