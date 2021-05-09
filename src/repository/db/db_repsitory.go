package db

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/sial-soft/oauth-api/src/client/cassandra"
	"github.com/sial-soft/oauth-api/src/domain/access_token"
	"github.com/sial-soft/oauth-api/src/utils/rest_errors"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES(?, ?, ?, ?);"
	queryUpdateExpire      = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(token access_token.AccessToken) *rest_errors.RestErr
	UpdateExpirationTime(token access_token.AccessToken) *rest_errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(accessTokenId string) (*access_token.AccessToken, *rest_errors.RestErr) {

	var accessToken access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, accessTokenId).Scan(
		&accessToken.AccessToken, &accessToken.UserId, &accessToken.ClientId, &accessToken.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no access token found with id:%s", accessTokenId))
		}
		return nil, rest_errors.NewInternalError(err.Error())
	}

	return &accessToken, nil
}

func (r *dbRepository) Create(token access_token.AccessToken) *rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryCreateAccessToken, token.AccessToken, token.UserId,
		token.ClientId, token.Expires).Exec(); err != nil {
		return rest_errors.NewInternalError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(token access_token.AccessToken) *rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryUpdateExpire, token.Expires, token.AccessToken).Exec(); err != nil {
		return rest_errors.NewInternalError(err.Error())
	}
	return nil
}

func NewDbRepository() DbRepository {
	return &dbRepository{}
}
