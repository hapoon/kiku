package kiku_test

import (
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hapoon/kiku"
	"github.com/stretchr/testify/assert"
)

func Test_PostTokenReissueParam_IsValid(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		p      kiku.PostTokenReissueParam
		expect error
	}{
		"Normal": {
			p: kiku.PostTokenReissueParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
			},
			expect: nil,
		},
		"Error LoginCompanyCode is empty": {
			p: kiku.PostTokenReissueParam{
				LoginCompanyCode: "",
				Token:            "bar",
			},
			expect: errors.New("LoginCompanyCode must be set"),
		},
		"Error Token is empty": {
			p: kiku.PostTokenReissueParam{
				LoginCompanyCode: "foo",
				Token:            "",
			},
			expect: errors.New("Token must be set"),
		},
	}

	for scenario, test := range tests {
		actual := test.p.IsValid()
		assert.Equal(t, test.expect, actual, scenario)
	}
}

func Test_PostTokenReissueResponse_Decode(t *testing.T) {
	t.Parallel()

	expiredDate := kiku.AkTime{
		Time: time.Date(2000, time.January, 2, 3, 4, 5, 0, time.UTC),
	}

	tests := map[string]struct {
		input  io.Reader
		expect kiku.PostTokenReissueResponse
		err    error
	}{
		"Success": {
			input: strings.NewReader(`{"success":true,"response":{"login_company_code":"foo","staff_id":123,"agency_manager_id":456,"token":"new_token","expired_at":"2000/01/02 03:04:05"}}`),
			expect: kiku.PostTokenReissueResponse{
				LoginCompanyCode: "foo",
				StaffId:          123,
				AgencyManagerId:  456,
				Token:            "new_token",
				ExpiredAt:        &expiredDate,
			},
			err: nil,
		},
		"Failed": {
			input:  strings.NewReader(`{"success":false,"response":{},"errors":[{"code":"error","message":"error"}]}`),
			expect: kiku.PostTokenReissueResponse{},
			err:    errors.New("Requesting Token Reissue API failed"),
		},
	}

	for scenario, test := range tests {
		var actual kiku.PostTokenReissueResponse
		err := actual.Decode(test.input)
		assert.Equal(t, test.err, err, scenario)
		assert.Equal(t, test.expect, actual, scenario)
	}
}
