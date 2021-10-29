package auth

import (
	"testing"

	"github.com/script-development/RT-CV/mock"
	"github.com/stretchr/testify/assert"
)

func TestAuthHelper(t *testing.T) {
	db := mock.NewMockDB()

	helper := NewHelper(db)

	testCases := []struct {
		name      string
		id        string
		key       string
		expectErr error
	}{
		// Valid
		{"valid mock key1", mock.Key1.ID.Hex(), mock.Key1.Key, nil},
		{"valid mock key1 from cache", mock.Key1.ID.Hex(), mock.Key1.Key, nil},
		{"valid mock key2", mock.Key2.ID.Hex(), mock.Key2.Key, nil},
		{"valid mock key3", mock.Key3.ID.Hex(), mock.Key3.Key, nil},
		{"valid mock dashboard key", mock.DashboardKey.ID.Hex(), mock.DashboardKey.Key, nil},
		// Invalid
		{"to short key", mock.Key1.ID.Hex()[:23], mock.Key1.Key, ErrAuthHeaderHasInvalidLen},
		{"not a hex id", mock.Key1.ID.Hex()[:20] + "++++", mock.Key1.Key, ErrAuthHeaderInvalidFormat},
		{"id does not exists", mock.Key1.ID.Hex()[:10] + "00000000000000", mock.Key1.Key, ErrAuthHeaderInvalid},
		{"key is invalid", mock.Key1.ID.Hex(), "this is a invalid key", ErrAuthHeaderInvalid},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			credentialsHeader := GenAuthHeaderKey(testCase.id, testCase.key)

			res, err := helper.Valid(credentialsHeader)
			if testCase.expectErr != nil {
				assert.Nil(t, res)
				assert.NotNil(t, err)
				assert.Equal(t, testCase.expectErr.Error(), err.Error())
			} else {
				assert.NotNil(t, res)
				assert.Nil(t, err)
			}
		})
	}
}
