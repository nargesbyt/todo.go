package handler

type Problem struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func NewProblem(code int, detail string) Problem {
	p := Problem{
		Code:   code,
		Detail: detail,
	}
	return p

}
