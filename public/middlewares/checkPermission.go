package middlewares

import (
	"fmt"
	"strings"

	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/utils/stringUtils"
)

type PermissionCheck struct{}

var permissionCheck *PermissionCheck

func NewPermissionCheck() *PermissionCheck {
	return permissionCheck
}

// 实现权限校验
//
//	needRoles:
//	"role1,role2" 需要所有列出的角色
//	"any:role1,role2,role3" 需要所列角色里的其中一个
//	"not:role1,role2,role3" 需要不具备role1/2/3任一角色，也就是除了这3个角色都可以访问
//
//	needPermissions:
//	"p1,p2,p3" 需要拥有这2个权限才能访问
//	"any:p1,p2,p3" 拥有p1/2/3任一个权限就可以访问
//	"not:p1,p2,p3" 拥有p1/2/3任一个权限的用户不能访问
func (p *PermissionCheck) Check(path, needRoles, needPermissions string, sessionUser *context.SessionUser) bool {
	if needRoles == "" && needPermissions == "" {
		return true
	}
	userRoles := sessionUser.Roles
	userPermissions := sessionUser.Permissions
	// fmt.Printf("------------->>>check permission-------------\n path: %v \n needRoles: %v \n needPermissions: %v \n roles: %v \n permissions: %v \n",
	// 	path,
	// 	needRoles,
	// 	needPermissions,
	// 	userRoles,
	// 	userPermissions)
	if len(userRoles) == 0 || len(userPermissions) == 0 {
		return false
	}

	ANY, NOT, SEP, ALLPERMISSION := "any:", "not:", ",", "*:*:*"
	// 角色校验
	if needRoles != "" {
		if strings.HasPrefix(needRoles, ANY) {
			needRoleList := strings.Split(needRoles[len(ANY):], SEP)
			hasRole := false
			for _, role := range needRoleList {
				if stringUtils.ContainsStr(userRoles, role) {
					hasRole = true
					break
				}
			}
			if !hasRole {
				fmt.Println("role check fail, with any.")
				return false
			}
		} else if strings.HasPrefix(needRoles, NOT) {
			needRoleList := strings.Split(needRoles[len(NOT):], SEP)
			for _, role := range needRoleList {
				if stringUtils.ContainsStr(userRoles, role) {
					fmt.Println("role check fail, with not.")
					return false
				}
			}
		} else {
			needRoleList := strings.Split(needRoles, SEP)
			for _, role := range needRoleList {
				if !stringUtils.ContainsStr(userRoles, role) {
					fmt.Println("role check fail.")
					return false
				}
			}
		}
	}

	// 权限校验
	if needPermissions != "" {
		if stringUtils.ContainsStr(userPermissions, ALLPERMISSION) {
			fmt.Println("permission all pass.")
			return true
		}
		if strings.HasPrefix(needPermissions, ANY) {
			needPermList := strings.Split(needPermissions[len(ANY):], SEP)
			hasPerm := false
			for _, perm := range needPermList {
				if stringUtils.ContainsStr(userPermissions, perm) {
					hasPerm = true
					break
				}
			}
			if !hasPerm {
				fmt.Println("permission check fail, with any.")
				return false
			}
		} else if strings.HasPrefix(needPermissions, NOT) {
			needPermList := strings.Split(needPermissions[len(NOT):], SEP)
			for _, perm := range needPermList {
				if stringUtils.ContainsStr(userPermissions, perm) {
					fmt.Println("permission check fail, with not.")
					return false
				}
			}
		} else {
			needPermList := strings.Split(needPermissions, SEP)
			for _, perm := range needPermList {
				if !stringUtils.ContainsStr(userPermissions, perm) {
					fmt.Println("permission check fail.")
					return false
				}
			}
		}
	}

	return true
}
