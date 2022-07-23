package kiku_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/hapoon/kiku"
	"github.com/stretchr/testify/assert"
)

func Test_GetStaffParam_EncodeURL(t *testing.T) {
	target := "baz"
	staffID := 123
	page := 2
	tests := map[string]struct {
		g      kiku.GetStaffParam
		expect string
		err    error
	}{
		"Normal scenario (minimum)": {
			g: kiku.GetStaffParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
			},
			expect: "/foo/staffs?token=bar",
			err:    nil,
		},
		"Normal scenario (all)": {
			g: kiku.GetStaffParam{
				LoginCompanyCode: "foo",
				Token:            "bar",
				Target:           &target,
				StaffID:          &staffID,
				Page:             &page,
			},
			expect: "/foo/staffs/123?page=2&target=baz&token=bar",
			err:    nil,
		},
		"Error scenario: LoginCompanyCode is empty": {
			g: kiku.GetStaffParam{
				LoginCompanyCode: "",
				Token:            "foo",
			},
			expect: "",
			err:    errors.New("LoginCompanyCode must be set"),
		},
		"Error scenario: Token is empty": {
			g: kiku.GetStaffParam{
				LoginCompanyCode: "foo",
				Token:            "",
			},
			expect: "",
			err:    errors.New("Token must be set"),
		},
	}

	for scenario, test := range tests {
		actual, err := test.g.EncodeURL()
		switch test.err {
		case nil:
			assert.NoError(t, err, fmt.Sprintf("Error occurred in %s", scenario))
		default:
			assert.Error(t, err, fmt.Sprintf("%s: %v", scenario, err))
			assert.Equal(t, test.err, err)
		}
		assert.Equal(t, test.expect, actual, scenario)
	}
}

func Test_GetStaffResponse_DecodeFrom(t *testing.T) {
	tests := map[string]struct {
		g      io.Reader
		expect *kiku.GetStaffResponse
		err    error
	}{
		"Normal scenario": {
			g: strings.NewReader(`
				{
					"success":true,
					"response":{
						"login_company_code":"foo",
						"Count":1,
						"TotalCount":1,
						"staffs":[
							{
								"staffId":1,
								"lastName":"愛",
								"firstName":"上大",
								"lastNameKana":"あい",
								"firstNameKana":"うえお",
								"organization":{},
								"subgroups":[],
								"employmentCategory":{},
								"tag":"bar",
								"staffNum":"123",
								"idmNum":"456",
								"cardTypeId":123,
								"remarks":"baz",
								"permissionGroup":{},
								"managedOrganizations":[]
							}
						]
					},
					"errors":[]
				}`),
			expect: &kiku.GetStaffResponse{
				LoginCompanyCode: "foo",
				Count:            1,
				TotalCount:       1,
				Staffs: []kiku.Staff{
					{
						ID:                   1,
						LastName:             "愛",
						FirstName:            "上大",
						LastNameKana:         "あい",
						FirstNameKana:        "うえお",
						Organization:         kiku.Organization{},
						SubGroups:            []kiku.Organization{},
						EmploymentCategory:   kiku.EmploymentCategory{},
						Tag:                  "bar",
						StaffNum:             "123",
						IDmNum:               "456",
						CardTypeID:           123,
						Remarks:              "baz",
						PermissionGroup:      kiku.PermissionGroup{},
						ManagedOrganizations: []kiku.Organization{},
					},
				},
			},
			err: nil,
		},
		"Error scenario: decode error": {
			g: strings.NewReader(`{
				"success":"foo"
			}`),
			expect: &kiku.GetStaffResponse{},
			err:    errors.New("Unmarshal error: field: success, value: string"),
		},
		"Error scenario: API failed": {
			g: strings.NewReader(`{
				"success":false
			}`),
			expect: &kiku.GetStaffResponse{},
			err:    errors.New("AKASHI API failed"),
		},
	}

	for scenario, test := range tests {
		actual := &kiku.GetStaffResponse{}
		err := actual.DecodeFrom(test.g)
		switch test.err {
		case nil:
			assert.NoError(t, err, scenario)
		default:
			assert.Error(t, err, scenario)
			assert.Equal(t, test.err, err, scenario)
		}
		assert.Equal(t, test.expect, actual, scenario)
	}
}
