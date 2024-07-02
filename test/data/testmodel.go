package data

type A1 struct {
	Id  int    `json:"id"`
	Str string `json:"str"`
}

type A2 struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	CTime int64  `json:"ctime"`
}

type A3 struct {
	Str   string `json:"str"`
	Num   int    `json:"num"`
	UTime int64  `json:"utime"`
}
