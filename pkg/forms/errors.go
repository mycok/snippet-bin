package forms

type errors map[string][]string

func (e errors) Add(field, errMessage string) {
	e[field] = append(e[field], errMessage)
}

func (e errors) Get(field string) string {
	errString := e[field]
	if len(errString) == 0 {
		return ""
	}
	return errString[0]
}
