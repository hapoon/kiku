package kiku

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Stamp is the struct representing stamp information.
type Stamp struct {
	StampedAt  *AkTime        `json:"stamped_at"` // 打刻日時
	Type       StampType      `json:"type"`       // 打刻種別
	LocalTime  *AkTime        `json:"local_time"` // 打刻機のローカル打刻時刻
	Timezone   string         `json:"timezone"`   // 打刻機のタイムゾーン
	Attributes StampAttribute `json:"attributes"` // 実績参照結果
}

const (
	// DateFormat yyyymmddHHMMSS形式
	DateFormat = "20060102150405"
	// ReturnDateFormat yyyy/mm/dd HH:MM:SS形式
	ReturnDateFormat = "2006/01/02 15:04:05"
)

// AkTime is the time for AKASHI.
type AkTime struct {
	time.Time
}

// UnmarshalJSON is the function that extends unmarshalJSON.
func (a *AkTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return
	}
	t, err := time.Parse(`"`+ReturnDateFormat+`"`, string(data))
	*a = AkTime{t}
	return
}

// StampType is the integer represents stamp's type.
type StampType int

const (
	// StampTypeUnknown 打刻種別:不明
	StampTypeUnknown StampType = 0
	// StampTypeGoToWork 打刻種別:出勤
	StampTypeGoToWork StampType = 11
	// StampTypeLeaveWork 打刻種別:退勤
	StampTypeLeaveWork StampType = 12
	// StampTypeGoStraight 打刻種別:直行
	StampTypeGoStraight StampType = 21
	// StampTypeBounce 打刻種別:直帰
	StampTypeBounce StampType = 22
	// StampTypeBreak 打刻種別:休憩入
	StampTypeBreak StampType = 31
	// StampTypeBreakReturn 打刻種別:休憩戻
	StampTypeBreakReturn StampType = 32
)

func (s StampType) String() string {
	switch s {
	case StampTypeGoToWork:
		return "出勤"
	case StampTypeLeaveWork:
		return "退勤"
	case StampTypeGoStraight:
		return "直行"
	case StampTypeBounce:
		return "直帰"
	case StampTypeBreak:
		return "休憩入"
	case StampTypeBreakReturn:
		return "休憩戻"
	default:
		return ""
	}
}

// StampAttribute is the struct represents 打刻実績参照結果
type StampAttribute struct {
	Method      int     `json:"method"`       // 打刻方法
	OrgID       int     `json:"org_id"`       // 組織ID
	WorkplaceID int     `json:"workplace_id"` // 勤務地ID
	Latitude    float32 `json:"latitude"`     // 緯度
	Longitude   float32 `json:"longitude"`    // 経度
	IP          string  `json:"ip"`           // 打刻機のIPアドレス
}

// GetStampParam is the struct represents for the request parameters of GET stamp API.
type GetStampParam struct {
	LoginCompanyCode string     // AKASHI企業ID
	Token            string     // アクセストークン
	StartDate        *time.Time // 打刻取得期間の開始日時
	EndDate          *time.Time // 打刻取得期間の終了日時
	StaffID          int        // 取得対象の従業員ID
}

// IsValid is the function to verify GetStampParam is correct.
func (g GetStampParam) IsValid() (err error) {
	switch {
	case g.LoginCompanyCode == "":
		err = errors.New("LoginCompanyCode must be set")
	case g.Token == "":
		err = errors.New("Token must be set")
	case g.StartDate == nil:
		err = errors.New("StartDate must be set")
	case g.EndDate == nil:
		err = errors.New("EndDate must be set")
	}
	return
}

// EncodeURL is the function to encode URL.
func (g GetStampParam) EncodeURL() (encodedURL string) {
	encodedURL = fmt.Sprintf("/%s/stamps", g.LoginCompanyCode)
	if g.StaffID != 0 {
		encodedURL = fmt.Sprintf("%s/%d", encodedURL, g.StaffID)
	}

	uv := url.Values{}
	if g.Token != "" {
		uv.Add("token", g.Token)
	}
	if g.StartDate != nil {
		uv.Add("start_date", g.StartDate.Format(DateFormat))
	}
	if g.EndDate != nil {
		uv.Add("end_date", g.EndDate.Format(DateFormat))
	}
	q := uv.Encode()
	if q != "" {
		encodedURL += "?" + q
	}

	return
}

// GetStampResponse is the struct represents for the response of GET stamp API.
type GetStampResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // ログイン企業ID
	StaffID          int     `json:"staff_id"`           // 従業員ID
	Count            int     `json:"count"`              // 従業員の打刻数
	Stamps           []Stamp `json:"stamps"`             // 打刻データの配列
}

// Decode is the function to decode GetStampResponse.
func (g *GetStampResponse) Decode(r io.Reader) (err error) {
	var gsr struct {
		Success  bool             `json:"success"`
		Response GetStampResponse `json:"response"`
		Errors   []Error          `json:"errors"`
	}

	err = json.NewDecoder(r).Decode(&gsr)
	switch {
	case !gsr.Success:
		err = errors.New("Requesting Stamps API failed")
		fallthrough
	case err != nil:
		return
	default:
		*g = gsr.Response
	}
	return
}

// GetStamps is the function that retrieves stamp information from AKASHI.
func GetStamps(ctx context.Context, param GetStampParam) (response GetStampResponse, err error) {
	if err = param.IsValid(); err != nil {
		return
	}

	endpoint := param.EncodeURL()

	cli := newClient()
	res, err := cli.Get(ctx, endpoint)
	switch {
	case res.StatusCode != http.StatusOK:
		err = fmt.Errorf("Status code=%d", res.StatusCode)
		fallthrough
	case err != nil:
		return
	}

	err = response.Decode(res.Body)

	return
}

type PostStampParam struct {
	LoginCompanyCode string    // AKASHI企業ID
	Token            string    `json:"token"`               // アクセストークン
	Type             StampType `json:"type,omitempty"`      // 打刻種別
	StampedAt        *AkTime   `json:"stampedAt,omitempty"` // クライアントでの打刻日時
	Timezone         string    `json:"timezone,omitempty"`  // クライアントでのタイムゾーン
}

func (p PostStampParam) IsValid() (err error) {
	switch {
	case p.LoginCompanyCode == "":
		err = errors.New("LoginCompanyCode must be set")
	case p.Token == "":
		err = errors.New("Token must be set")
	}
	return
}

func (p PostStampParam) EncodeURL() (encodedURL string, err error) {
	encodedURL = fmt.Sprintf("/%s/stamps", p.LoginCompanyCode)
	return
}

type PostStampResponse struct {
	LoginCompanyCode string    `json:"login_company_code"` // AKASHI企業ID
	StaffID          int       `json:"staff_id"`           // 従業員ID
	Type             StampType `json:"type"`               // 打刻種別
	StampedAt        *AkTime   `json:"stampedAt"`          // サーバ側での打刻日時
}

func (p *PostStampResponse) Decode(r io.Reader) (err error) {
	var psr struct {
		Success  bool              `json:"success"`
		Response PostStampResponse `json:"response"`
		Errors   []Error           `json:"errors"`
	}

	err = json.NewDecoder(r).Decode(&psr)
	switch {
	case !psr.Success:
		err = errors.New("Requesting Stamp API failed")
		fallthrough
	case err != nil:
		return
	default:
		*p = psr.Response
	}
	return
}

func PostStamp(ctx context.Context, param PostStampParam) (response PostStampResponse, err error) {
	if err = param.IsValid(); err != nil {
		return
	}

	endpoint, err := param.EncodeURL()
	if err != nil {
		return
	}

	cli := newClient()
	res, err := cli.Post(ctx, endpoint, param)
	switch {
	case res.StatusCode != http.StatusOK:
		err = fmt.Errorf("Status code=%d", res.StatusCode)
		fallthrough
	case err != nil:
		return
	}

	err = response.Decode(res.Body)
	fmt.Printf("response: %+v\n", response)

	return
}
