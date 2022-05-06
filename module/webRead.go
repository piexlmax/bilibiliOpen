package module

type WebCMD struct {
	CMD  string `json:"cmd"`
	Data string `json:"data"`
}

type RegisterData struct {
	Name string
}
