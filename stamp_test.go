package kiku_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hapoon/kiku"
	"github.com/stretchr/testify/assert"
)

func Test_GetStampParam_IsValid(t *testing.T) {
	startDate := time.Date(2000, time.January, 2, 3, 4, 5, 0, time.UTC)
	endDate := time.Date(2000, time.February, 3, 4, 5, 6, 0, time.UTC)
	tests := map[string]struct {
		g      kiku.GetStampParam
		expect error
	}{
		"Normal scenario": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				StartDate:        &startDate,
				EndDate:          &endDate,
				StaffID:          1,
			},
			expect: nil,
		},
		"Error scenario: LoginCompanyCode is empty": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "",
				Token:            "bar",
				StartDate:        &startDate,
				EndDate:          &endDate,
				StaffID:          1,
			},
			expect: errors.New("LoginCompanyCode must be set"),
		},
		"Error scenario: Token is empty": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "",
				StartDate:        &startDate,
				EndDate:          &endDate,
				StaffID:          1,
			},
			expect: errors.New("Token must be set"),
		},
		"Error scenario: StartDate is empty": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				StartDate:        nil,
				EndDate:          &endDate,
				StaffID:          1,
			},
			expect: errors.New("StartDate must be set"),
		},
		"Error scenario: EndDate is empty": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				StartDate:        &startDate,
				EndDate:          nil,
				StaffID:          1,
			},
			expect: errors.New("EndDate must be set"),
		},
	}

	for scenario, test := range tests {
		actual := test.g.IsValid()
		assert.Equal(t, test.expect, actual, scenario)
	}
}

func Test_GetStampParam_EncodeURL(t *testing.T) {
	startDate := time.Date(2000, time.January, 2, 3, 4, 5, 0, time.UTC)
	endDate := time.Date(2000, time.February, 3, 4, 5, 6, 0, time.UTC)
	tests := map[string]struct {
		g      kiku.GetStampParam
		expect string
	}{
		"All parameter set": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				StartDate:        &startDate,
				EndDate:          &endDate,
				StaffID:          1,
			},
			expect: "/foo/stamps/1?end_date=20000203040506&start_date=20000102030405&token=bar",
		},
		"login_company_code,token,staff_id only set": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				StaffID:          1,
			},
			expect: "/foo/stamps/1?token=bar",
		},
		"login_company_code&token only set": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
			},
			expect: "/foo/stamps?token=bar",
		},
	}

	for scenario, test := range tests {
		actual := test.g.EncodeURL()
		assert.Equal(t, test.expect, actual, scenario)
	}
}

func Test_PostStampParam_IsValid(t *testing.T) {
	tests := map[string]struct {
		p   kiku.PostStampParam
		err error
	}{
		"Necessary parameter set": {
			p: kiku.PostStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
			},
		},
		"LoginCompanyCode is not set": {
			p: kiku.PostStampParam{
				LoginCompanyCode: "",
				Token:            "bar",
			},
			err: errors.New("LoginCompanyCode must be set"),
		},
		"Token is not set": {
			p: kiku.PostStampParam{
				LoginCompanyCode: "foo",
				Token:            "",
			},
			err: errors.New("Token must be set"),
		},
	}

	for scenario, test := range tests {
		err := test.p.IsValid()
		switch test.err {
		case nil:
			assert.NoError(t, err, fmt.Sprintf("Error occurred in %s", scenario))
		default:
			assert.Error(t, err, fmt.Sprintf("%s: %v", scenario, err))
			assert.Equal(t, test.err, err)
		}
	}
}

func Test_PostStampResponse_Decode(t *testing.T) {
	tests := map[string]struct {
		input  io.Reader
		expect kiku.PostStampResponse
		err    error
	}{
		"Success": {
			input: strings.NewReader(`{"success":true,"response":{"login_company_code":"foo","staff_id":1,"type":11}}`),
			expect: kiku.PostStampResponse{
				LoginCompanyCode: "foo",
				StaffID:          1,
				Type:             kiku.StampTypeGoToWork,
			},
			err: nil,
		},
		"Request fail": {
			input:  strings.NewReader(`{"success":false}`),
			expect: kiku.PostStampResponse{},
			err:    errors.New("Requesting Stamp API failed"),
		},
	}

	for scenario, test := range tests {
		var actual kiku.PostStampResponse
		err := actual.Decode(test.input)
		assert.Equal(t, test.err, err, scenario)
		assert.Equal(t, test.expect, actual, scenario)
	}
}
