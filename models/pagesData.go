package models

type PageData struct {
	Title string
	Loged bool
	Username string
}

type IndexData struct {
	PageData
	Posts []PostI
}

type PostI struct { // lista de post del index
	IdPost int
	Titulo string
	IdCategoria int
}