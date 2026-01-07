package serverctx

type contextKey string

const (
	AccessTokenKey contextKey = "access-token"
	UserIDKey      contextKey = "user-id"
)
