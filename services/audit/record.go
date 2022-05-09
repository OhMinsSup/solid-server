package audit

import "github.com/mattermost/mattermost-server/v6/shared/mlog"

// Meta 레코드에 이름/값 쌍으로 추가할 수 있는 메타데이터를 나타냅니다.
type Meta struct {
	K string
	V interface{}
}

// FuncMetaTypeConv 메타 데이터 유형을 무언가로 변환할 수 있는 함수를 정의합니다.
// 레코드에 대해 직렬화
type FuncMetaTypeConv func(val interface{}) (newVal interface{}, converted bool)

// Record 모든 로깅에 사용되는 일관된 필드 집합을 제공합니다.
type Record struct {
	APIPath   string
	Event     string
	Status    string
	UserID    string
	SessionID string
	Client    string
	IPAddress string
	Meta      []Meta
	metaConv  []FuncMetaTypeConv
}

// Success 레코드 상태를 성공으로 표시합니다.
func (rec *Record) Success() {
	rec.Status = Success
}

// Success 레코드 상태를 실패로 표시합니다.
func (rec *Record) Fail() {
	rec.Status = Fail
}

// AddMeta  레코드의 메타데이터에 단일 이름/값 쌍을 추가합니다.
func (rec *Record) AddMeta(name string, val interface{}) {
	if rec.Meta == nil {
		rec.Meta = []Meta{}
	}

	// 0개 이상의 변환 함수를 통해
	//val을 직렬화에 더 적합한 것으로 변환할 수 있습니다.
	for _, conv := range rec.metaConv {
		converted, wasConverted := conv(val)
		if wasConverted {
			val = converted
			break
		}
	}

	lc, ok := val.(mlog.LogCloner)
	if ok {
		val = lc.LogClone()
	}

	rec.Meta = append(rec.Meta, Meta{K: name, V: val})
}

// AddMetaTypeConverter 메타 필드 유형을 변환할 수 있는 기능 추가
func (rec *Record) AddMetaTypeConverter(f FuncMetaTypeConv) {
	rec.metaConv = append(rec.metaConv, f)
}

