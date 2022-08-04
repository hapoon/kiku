package kiku_test

import (
	"errors"
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
		err    error
	}{
		"Normal scenario": {
			g: kiku.GetStampParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				StartDate:        &startDate,
				EndDate:          &endDate,
				StaffID:          1,
			},
			expect: "",
			err:    nil,
		},
	}

	for scenario, test := range tests {
		actual, err := test.g.EncodeURL()
		assert.Equal(t, test.expect, actual, scenario)
		assert.Equal(t, test.err, err, scenario)
	}
}
