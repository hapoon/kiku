package kiku

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// PostTokenReissueParam トークン再発行リクエストパラメータ
type PostTokenReissueParam struct {
	LoginCompanyCode string // AKASHI企業ID
	Token            string `json:"token"` // アクセストークン
}

// IsValid リクエストパラメータが正しいか確認
func (p PostTokenReissueParam) IsValid() (err error) {
	switch {
	case p.LoginCompanyCode == "":
		err = errors.New("LoginCompanyCode must be set")
	case p.Token == "":
		err = errors.New("Token must be set")
	}
	return
}

// EncodeUrl URLを生成
func (p PostTokenReissueParam) EncodeUrl() (encodedUrl string, err error) {
	encodedUrl = fmt.Sprintf("/token/reissue/%s", p.LoginCompanyCode)
	return
}

// PostTokenReissueResponse トークン再発行レスポンス
type PostTokenReissueResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // ログイン企業ID
	StaffId          int     `json:"staff_id"`           // 従業員ID
	AgencyManagerId  int     `json:"agency_manager_id"`  // -
	Token            string  `json:"token"`              // アクセストークン
	ExpiredAt        *AkTime `json:"expired_at"`         // アクセストークンの有効期限
}

func (p *PostTokenReissueResponse) Decode(r io.Reader) (err error) {
	var psr struct {
		Success  bool                     `json:"success"`
		Response PostTokenReissueResponse `json:"response"`
		Errors   []Error                  `json:"errors"`
	}

	err = json.NewDecoder(r).Decode(&psr)
	switch {
	case !psr.Success:
		err = errors.New("Requesting Token Reissue API failed")
		fallthrough
	case err != nil:
		return
	default:
		*p = psr.Response
	}
	return
}

func PostTokenReissue(ctx context.Context, param PostTokenReissueParam) (res PostTokenReissueResponse, err error) {
	if err = param.IsValid(); err != nil {
		return
	}

	endpoint, err := param.EncodeUrl()
	if err != nil {
		return
	}

	cli := newClient()
	r, err := cli.Post(ctx, endpoint, param)
	switch {
	case r.StatusCode != http.StatusOK:
		err = fmt.Errorf("Status code=%d", r.StatusCode)
		fallthrough
	case err != nil:
		return
	}

	err = res.Decode(r.Body)

	return
}
