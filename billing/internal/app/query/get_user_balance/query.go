package get_user_balance

import "context"

type Query struct{}

func New() *Query {
	return &Query{}
}

func (q Query) Handle(_ context.Context) error {
	return nil
}
