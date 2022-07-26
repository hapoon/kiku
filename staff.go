package kiku

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Staff is the struct representing employee information.
type Staff struct {
	ID                   int                `json:"staffId"`              // 従業員ID
	LastName             string             `json:"lastName"`             // 姓
	FirstName            string             `json:"firstName"`            // 名
	LastNameKana         string             `json:"lastNameKana"`         // カナ(姓)
	FirstNameKana        string             `json:"firstNameKana"`        // カナ(名)
	Organization         Organization       `json:"organization"`         // 組織(メイン)
	SubGroups            []Organization     `json:"subgroups"`            // 組織(サブグループ)の配列
	EmploymentCategory   EmploymentCategory `json:"employmentCategory"`   // 雇用区分
	Tag                  string             `json:"tag"`                  // タグ
	StaffNum             string             `json:"staffNum"`             // 従業員番号
	IDmNum               string             `json:"idmNum"`               // IDm番号
	CardTypeID           int                `json:"cardTypeId"`           // カード種別
	Remarks              string             `json:"remarks"`              // 鼻腔
	PermissionGroup      PermissionGroup    `json:"permissionGroup"`      // 権限グループ
	ManagedOrganizations []Organization     `json:"managedOrganizations"` // 管理対象組織
}

// EmploymentCategory is the struct representing an employee's employment category.
type EmploymentCategory struct {
	ID   int    `json:"employmentCategoryId"` // 雇用区分ID
	Name string `json:"Name"`                 // 雇用区分名称
}

// PermissionGroup is the struct representing authorisation information.
type PermissionGroup struct {
	ID   int    `json:"permissionGroupId"` // 権限グループID
	Type int    `json:"permissionType"`    // 権限種別(1:企業管理者,2:一般管理者,3:従業員)
	Name string `json:"name"`              // 権限グループ名
}

// GetStaffParam is the struct for the request parameters of GET Employee API.
type GetStaffParam struct {
	LoginCompanyCode string  // AKASHI企業ID
	Token            string  // アクセストークン
	Target           *string // 取得する従業員のトークン
	StaffID          *int    // 取得対象の従業員ID
	Page             *int    // 管理下にある従業員をすべて取得する場合のページ番号
}

// EncodeURL is the function that encodes the request parameter URL of GET Employee API.
func (g GetStaffParam) EncodeURL() (encodedURL string, err error) {
	switch {
	case g.LoginCompanyCode == "":
		err = errors.New("LoginCompanyCode must be set")
		return
	case g.Token == "":
		err = errors.New("Token must be set")
		return
	}

	encodedURL = fmt.Sprintf("/%s/staffs", g.LoginCompanyCode)
	if g.StaffID != nil {
		encodedURL += fmt.Sprintf("/%d", *g.StaffID)
	}

	uv := url.Values{}
	uv.Add("token", g.Token)
	switch {
	case g.Target != nil:
		uv.Add("target", *g.Target)
		fallthrough
	case g.Page != nil:
		uv.Add("page", strconv.FormatInt(int64(*g.Page), 10))
	}
	if q := uv.Encode(); q != "" {
		encodedURL += "?" + q
	}

	return
}

// GetStaffResponse is the struct representing the response of GET Employee API.
type GetStaffResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // AKASHI企業ID
	Count            int     `json:"Count"`              // 取得された従業員数
	TotalCount       int     `json:"TotalCount"`         // 取得することができる従業員数
	Staffs           []Staff `json:"staffs"`             // 取得した従業員情報の配列
}

// DecodeFrom is the function that decodes the response of GET Employee API.
func (g *GetStaffResponse) DecodeFrom(r io.Reader) (err error) {
	var decoded struct {
		Success  bool             `json:"success"`
		Response GetStaffResponse `json:"response"`
		Errors   []Error          `json:"errors"`
	}
	if err = json.NewDecoder(r).Decode(&decoded); err != nil {
		e := err.(*json.UnmarshalTypeError)
		err = fmt.Errorf("Unmarshal error: field: %s, value: %s", e.Field, e.Value)
		return
	}
	if !decoded.Success {
		err = errors.New("AKASHI API failed")
	}
	*g = decoded.Response
	return
}

// GetStaff is the function that retrieves employee information from AKASHI.
func GetStaff(ctx context.Context, param GetStaffParam) (response GetStaffResponse, err error) {
	endpointURL, err := param.EncodeURL()
	if err != nil {
		return
	}

	cli := newClient()
	res, err := cli.Get(ctx, endpointURL)
	defer res.Body.Close()
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Status code=%d", res.StatusCode)
		return
	}

	err = response.DecodeFrom(res.Body)
	return
}
