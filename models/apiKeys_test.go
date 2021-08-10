package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestApiKeyRole(t *testing.T) {
	seneraios := []struct {
		name              string
		userRole          ApiKeyRole
		contains          ApiKeyRole
		expectedMatchAll  bool
		expectedMatchSome bool
	}{
		{
			"no roles should not match anything",
			ApiKeyRole(0),
			ApiKeyRoleAdmin,
			false,
			false,
		},
		{
			"single role should match equal single role",
			ApiKeyRoleController,
			ApiKeyRoleController,
			true,
			true,
		},
		{
			"multiple roles should match other matching role",
			ApiKeyRoleScraper | ApiKeyRoleInformationObtainer,
			ApiKeyRoleInformationObtainer,
			true,
			true,
		},
		{
			"single role mismatch",
			ApiKeyRoleScraper,
			ApiKeyRoleInformationObtainer,
			false,
			false,
		},
		{
			"multiple roles mismatch",
			ApiKeyRoleScraper | ApiKeyRoleInformationObtainer,
			ApiKeyRoleController | ApiKeyRoleAdmin,
			false,
			false,
		},
		{
			"multiple roles contain one match",
			ApiKeyRoleScraper | ApiKeyRoleInformationObtainer,
			ApiKeyRoleInformationObtainer | ApiKeyRoleAdmin,
			false,
			true,
		},
	}

	for _, s := range seneraios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			Equal(t, s.expectedMatchAll, s.userRole.ContainsAll(s.contains))
			Equal(t, s.expectedMatchSome, s.userRole.ContainsSome(s.contains))
		})
	}
}
