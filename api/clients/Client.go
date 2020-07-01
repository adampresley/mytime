package clients

type Client struct {
	ClientID int    `json:"clientID"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Archived bool   `json:"archived"`
}

type ClientCollection []Client

func (c Client) ID() (string, interface{}) {
	return "clientID", c.ClientID
}
