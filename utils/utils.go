package utils

import (
	mm_model "github.com/mattermost/mattermost-server/v6/model"
	"time"
)

type IDType byte

const (
	IDTypeNone    IDType = '7'
	IDTypeTeam    IDType = 't'
	IDTypeBoard   IDType = 'b'
	IDTypeCard    IDType = 'c'
	IDTypeView    IDType = 'v'
	IDTypeSession IDType = 's'
	IDTypeUser    IDType = 'u'
	IDTypeToken   IDType = 'k'
	IDTypeBlock   IDType = 'a'
)

// NewID 는 전역적으로 고유한 식별자입니다.
// 27자 길이의 [A-Z0-9] 문자열입니다. 패딩이 제거된
// zbased32로 인코딩된 UUID 버전 4 Guid와 엔티티의 유형을 나타내는 1자 알파 접두사 또는 알 수 없는 유형인 경우 '7'입니다.
func NewID(idType IDType) string {
	return string(idType) + mm_model.NewId()
}

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
