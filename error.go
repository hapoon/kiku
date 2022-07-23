package kiku

// Error 失敗の原因となったエラーオブジェクト
type Error struct {
	Code    string `json:"code"`    // エラーコード
	Message string `json:"message"` // エラーメッセージ
}
