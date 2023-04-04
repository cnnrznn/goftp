package client

type Client struct {
	filename string
	toURI    string
}

func New(
	fn string,
	toURI string,
) *Client {
	return &Client{
		filename: fn,
		toURI:    toURI,
	}
}

func (s *Client) Run(stopChan chan error) {
	// do stuff
}
