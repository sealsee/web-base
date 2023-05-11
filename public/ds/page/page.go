package page

const (
	PAGE_SIZE     = 10
	MAX_PAGE_SIZE = 1000
)

type Page struct {
	CurPage   int `json:"curPage"`
	PageSize  int `json:"pageSize"`
	TotalSize int `json:"totalSize"`
	TotalPage int `json:"totalPage"`
	Offset    int `json:"-"`
	Limit     int `json:"-"`
}

func NewPage() *Page {
	return &Page{CurPage: 1, PageSize: PAGE_SIZE}
}

func (p *Page) SetTotalSize(totalSize int) {
	p.TotalSize = totalSize
	if totalSize != 0 {
		if p.PageSize > MAX_PAGE_SIZE {
			p.PageSize = MAX_PAGE_SIZE
		}
		if totalSize%p.PageSize == 0 {
			p.TotalPage = totalSize / p.PageSize
		} else {
			p.TotalPage = totalSize/p.PageSize + 1
		}
		if p.CurPage > p.TotalPage {
			p.CurPage = p.TotalPage
		}
	}
}

func (p *Page) SetPageSize(pageSize int) {
	p.PageSize = pageSize
	if pageSize != 0 {
		if p.PageSize > MAX_PAGE_SIZE {
			p.PageSize = MAX_PAGE_SIZE
		}
	}
}

func (p *Page) GetOffset() int {
	if p.CurPage > 0 {
		p.Offset = (p.CurPage - 1) * p.PageSize
	}
	return p.Offset
}

func (p *Page) GetLimit() int {
	p.Limit = p.PageSize
	return p.Limit
}
