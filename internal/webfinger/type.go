package webfinger

type Webfinger struct {
	Subject string          `json:"subject"`
	Aliases []string        `json:"aliases"`
	Links   []WebfingerLink `json:"links"`
}

func (w *Webfinger) GetLink(findType string) (link *WebfingerLink, ok bool) {
	if w.Links == nil {
		return
	}
	for _, v := range w.Links {
		if v.Type == findType {
			link = &v
			ok = true
			return
		}
	}
	return
}

type WebfingerLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}
