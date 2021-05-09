package access_token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenConstant(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime)
}

func TestGetAccessToken(t *testing.T) {
	at := GetAccessToken()
	assert.NotNil(t, at)
	assert.False(t, at.IsExpired(), "new token should not be expired")
	assert.EqualValues(t, "", at.AccessToken, "brand new access token should be empty")
	assert.EqualValues(t, int64(0), at.UserId, "new access token should be zero")
}

func TestAccessToken_IsExpired(t *testing.T) {
	at := AccessToken{}
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "access token expiring in 3 hours should be expired")
}
