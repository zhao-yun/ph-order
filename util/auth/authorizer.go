package auth

import "github.com/gin-gonic/gin"

//type Authorizer struct {
//	authCacheDB *repository.AuthCacheDB
//}
//
//func NewAuthorizer(cacheDB *repository.AuthCacheDB) *Authorizer {
//	return &Authorizer{
//		authCacheDB: cacheDB,
//	}
//}
//
//func (a *Authorizer) Verify(token string) (*User, error) {
//	user, err := a.authCacheDB.GetUserCache(token)
//	if err != nil {
//		return nil, err
//	}
//	return &User{
//		CognitoID: user.CognitoID,
//		UserName:  user.UserName,
//		Email:     user.Email,
//	}, nil
//}
//
//func (a *Authorizer) Invalidate(token string) error {
//	return a.authCacheDB.InvalidateUserCache(token)
//}

func GetUserID(c *gin.Context) (string, error) {
	return "fc8d0548-8021-7060-4b96-6919963dfb07", nil
}

func GetUserIDFromToken(c *gin.Context) (string, error) {
	return "", nil
}

func GetSitterID(c *gin.Context) (string, error) {
	return "ec5d3518-b0d1-7084-9840-4304844857f9", nil
}
