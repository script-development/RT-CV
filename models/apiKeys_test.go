package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestApiKeyRole(t *testing.T) {
	seneraios := []struct {
		name              string
		userRole          APIKeyRole
		contains          APIKeyRole
		expectedMatchAll  bool
		expectedMatchSome bool
	}{
		{
			"no roles should not match anything",
			APIKeyRole(0),
			APIKeyRoleAdmin,
			false,
			false,
		},
		{
			"single role should match equal single role",
			APIKeyRoleController,
			APIKeyRoleController,
			true,
			true,
		},
		{
			"multiple roles should match other matching role",
			APIKeyRoleScraper | APIKeyRoleInformationObtainer,
			APIKeyRoleInformationObtainer,
			true,
			true,
		},
		{
			"single role mismatch",
			APIKeyRoleScraper,
			APIKeyRoleInformationObtainer,
			false,
			false,
		},
		{
			"multiple roles mismatch",
			APIKeyRoleScraper | APIKeyRoleInformationObtainer,
			APIKeyRoleController | APIKeyRoleAdmin,
			false,
			false,
		},
		{
			"multiple roles contain one match",
			APIKeyRoleScraper | APIKeyRoleInformationObtainer,
			APIKeyRoleInformationObtainer | APIKeyRoleAdmin,
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
