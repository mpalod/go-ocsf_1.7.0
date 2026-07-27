package main

import (
	"fmt"
	"testing"

	"github.com/telophasehq/go-ocsf/ocsf/v1_8_0"
)

type validatesObservableV180 interface {
	ValidateObservables() error
}

func strPtrV180(v string) *string {
	return &v
}

func TestValidateObservablesV180(t *testing.T) {
	tests := []struct {
		name string
		in   validatesObservableV180
		err  error
	}{
		{
			name: "vulnerability finding with no observable fields set",
			in:   &v1_8_0.VulnerabilityFinding{},
			err:  nil,
		},
		{
			name: "vulnerability finding missing user observable in list",
			in: &v1_8_0.VulnerabilityFinding{
				Resources: []v1_8_0.ResourceDetails{
					{Owner: &v1_8_0.User{Name: strPtrV180("jack")}},
				},
			},
			err: fmt.Errorf("non-null observable user(21) not found in observables array"),
		},
		{
			name: "vulnerability finding includes user observable in list",
			in: &v1_8_0.VulnerabilityFinding{
				Resources: []v1_8_0.ResourceDetails{
					{Owner: &v1_8_0.User{Name: strPtrV180("jack")}},
				},
				Observables: []v1_8_0.Observable{
					{Name: strPtrV180("user"), TypeId: 21},
				},
			},
			err: nil,
		},
		{
			name: "account change validates required user observable",
			in: &v1_8_0.AccountChange{
				User: v1_8_0.User{Name: strPtrV180("jane")},
				Observables: []v1_8_0.Observable{
					{Name: strPtrV180("user"), TypeId: 21},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.ValidateObservables()
			if tt.err == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.err != nil {
				if err == nil {
					t.Fatalf("got nil error, want %v", tt.err)
				}
				if err.Error() != tt.err.Error() {
					t.Fatalf("got %v, want %v", err, tt.err)
				}
			}
		})
	}
}
