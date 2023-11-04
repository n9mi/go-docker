package exception

//

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return e.Entity + " doesn't exists"
}

//

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

//

type NotFoundValidate struct {
	Entity string
}

func (e *NotFoundValidate) Error() string {
	return e.Entity + " doesn't exists"
}

//

type TokenError struct {
}

func (e *TokenError) Error() string {
	return "invalid token"
}
