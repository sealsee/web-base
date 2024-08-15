package context

import (
	"time"

	"github.com/sealsee/web-base/public/basemodel"
)

type SessionUser struct {
	UserId   int64          `json:"userId,string"`
	UserName string         `json:"userName"`
	NickName string         `json:"nickName"`
	DeptId   int64          `json:"deptId,string"`
	Avatar   string         `json:"avatar" `
	Password string         `json:"-"`
	IsUsed   bool           `json:"isUsed"`
	Ext      map[string]any `json:"ext,omitempty"`

	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`

	LoginDevice string             `json:"loginDevice"`
	LoginIp     string             `json:"loginIp"`
	LoginTime   basemodel.BaseTime `json:"loginTime"`
	Token       string             `json:"token"`
	ExpireTime  int64              `json:"expire_time"`
}

func NewSessionUser() *SessionUser {
	return &SessionUser{LoginTime: basemodel.BaseTime(time.Now())}
}

func (u *SessionUser) AddExt(k string, v any) {
	if k == "" || v == nil {
		return
	}
	if u.Ext == nil {
		u.Ext = make(map[string]any)
	}
	u.Ext[k] = v
}
