package user

import "context"

type userDataKeys int

const (
	userModelCtxKey userDataKeys = iota
)

type UserModel struct {
	Name string
}

func ContextWithUser(ctx context.Context, user UserModel) context.Context {
	return context.WithValue(ctx, userModelCtxKey, user)
}

func UserFromContext(ctx context.Context) (UserModel, bool) {
	user, ok := ctx.Value(userModelCtxKey).(UserModel)
	return user, ok
}
