package routes

type ViewData struct {
	Context string
	Name    string
}

type Auction struct {
	ID       string
	Seller   string
	Type     string
	Form     string
	MinPrice uint
}
