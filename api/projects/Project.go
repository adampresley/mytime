package projects

type Project struct {
	ProjectID         int    `json:"projectID"`
	Name              string `json:"name"`
	Code              string `json:"code"`
	ClientID          int    `json:"clientID"`
	DefaultCategoryID int    `json:"defaultCategoryID"`
	Archived          bool   `json:"archived"`
}

func (p Project) ID() (string, interface{}) {
	return "projectID", p.ProjectID
}

type ProjectCollection []Project

type ProjectSearch struct {
	Archived bool
	Client   string
	Name     string
}
