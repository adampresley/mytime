package categories

type Category struct {
	CategoryID int     `json:"categoryID"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Rate       float64 `json:"rate"`
	Archived   bool    `json:"archived"`
}

func (c Category) ID() (string, interface{}) {
	return "categoryID", c.CategoryID
}

type CategoryCollection []Category

type CategorySearch struct {
	Archived bool
	Name     string
}
