package utils

import (
	"time"
	mm_model "github.com/mattermost/mattermost-server/v6/model"
)

// GetMillis 는 Epoch 이후 밀리초를 가져오는 편리한 메서드입니다.
func GetMillis() int64 {
	return mm_model.GetMillis()
}

// GetMillisForTime 은 제공된 Time에 대한 epoch 이후 밀리초를 가져오는 편리한 메소드입니다.
func GetMillisForTime(thisTime time.Time) int64 {
	return mm_model.GetMillisForTime(thisTime)
}

// SecondsToMillis 는 초를 밀리초로 변환하는 편리한 방법입니다. func SecondsToMillis(초 int64) int64 { 반환 초 * 1000 }
func SecondsToMillis(seconds int64) int64 {
	return seconds * 1000
}
