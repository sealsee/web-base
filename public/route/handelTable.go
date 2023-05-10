package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/utils/set"
)

type MethodType int

const (
	GET        MethodType = 1
	POST       MethodType = 2
	GROUP_ROOT            = "-1"
)

var GroupsHander = make([]*groupTable, 0, 20)
var groupsPath = make(map[string]*groupTable)
var UrlName = make(map[string]string)

// var NoLoginHanders = make([]handelTable, 0, 20)
// var LoginHanders = make([]handelTable, 0, 100)

var methods = set.Set[MethodType]{}

func init() {
	methods.Add(GET)
	methods.Add(POST)
}

type groupTable struct {
	Group  string
	Handel []handelTable
	Mark   string
}

type handelTable struct {
	Method    MethodType
	Pattern   string
	Action    func(*gin.Context)
	NeedLogin bool
	Mark      string
}

func (h *groupTable) IsDefault() bool {
	return h.Group == GROUP_ROOT
}

func (h *handelTable) IsGet() bool {
	return h.Method == GET
}

func UseDefaultGroup() (g *groupTable) {
	return AddGroup(GROUP_ROOT, "")
}

func AddGroup(group string, mark string) (g *groupTable) {
	if group == "" {
		group = GROUP_ROOT
	}

	g, ok := groupsPath[group]
	if !ok {
		g = &groupTable{Group: group, Handel: make([]handelTable, 0, 10), Mark: mark}
		groupsPath[group] = g
		GroupsHander = append(GroupsHander, g)
	}
	return
}

func (g *groupTable) AddHandel(method MethodType, pattern string, action func(*gin.Context), needLogin bool, mark string) *groupTable {
	if pattern == "" || action == nil {
		return g
	}
	if !methods.Contains(method) {
		method = POST
	}

	ht := handelTable{Method: method, Pattern: pattern, Action: action, NeedLogin: needLogin}
	g.Handel = append(g.Handel, ht)

	absolutePath := ""
	if g.Group == "" || g.IsDefault() {
		absolutePath = pattern
	} else {
		if g.Group[len(g.Group)-1:] == "/" {
			absolutePath = g.Group[:len(g.Group)-1] + pattern
		} else {
			absolutePath = g.Group + pattern
		}
	}

	pathName := mark
	if g.Mark != "" {
		pathName = g.Mark + "->" + mark
	}
	UrlName[absolutePath] = pathName
	return g
}

func (g *groupTable) AddGETHandel(pattern string, action func(*gin.Context), needLogin bool, mark string) *groupTable {
	return g.AddHandel(GET, pattern, action, needLogin, mark)
}

func (g *groupTable) AddPOSTHandel(pattern string, action func(*gin.Context), needLogin bool, mark string) *groupTable {
	return g.AddHandel(POST, pattern, action, needLogin, mark)
}

func GetUrlMark(url string) string {
	return UrlName[url]
}
