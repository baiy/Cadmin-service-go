package userGroupRelate

import (
	"errors"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/doug-martin/goqu/v9"
	"github.com/juliangruber/go-intersect"
)

type Model struct {
	models.Model
	AdminUserGroupId string `json:"admin_user_group_id"`
	AdminAuthId      string `json:"admin_auth_id"`
}

func AuthIds(userGroupIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_user_group_relate").Select("admin_auth_id").Where(goqu.Ex{
		"admin_user_group_id": userGroupIds,
	}).ScanVals(&ids)
	return ids
}

// 用户分组权限检查
func Check(authIds []int, userGroupIds []int) bool {
	if len(authIds) == 0 || len(userGroupIds) == 0 {
		return false
	}
	existAuthIds := AuthIds(userGroupIds)

	if len(existAuthIds) == 0 {
		return false
	}
	return len(intersect.Simple(existAuthIds, authIds).([]interface{})) != 0
}

func Remove(userGroupId, authId int) error {
	if userGroupId == 0 && authId == 0 {
		return errors.New("参数错误")
	}
	where := make(goqu.Ex)
	if userGroupId != 0 {
		where["admin_user_group_id"] = userGroupId
	}
	if authId != 0 {
		where["admin_auth_id"] = authId
	}
	_, err := models.Db.Delete("admin_user_group_relate").Where(where).Executor().Exec()
	return err
}
