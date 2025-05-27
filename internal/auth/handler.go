package auth

type Handle struct {
	authorization *AuthorizationServiceImpl
}

func NewHandle(authorization *AuthorizationServiceImpl) *Handle {
	return &Handle{
		authorization: authorization,
	}
}
