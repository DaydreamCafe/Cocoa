package tests

import (
	"testing"

	"github.com/DaydreamCafe/Cocoa/api/permission_api"
)

func TestResetAllPermission(t *testing.T) {
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	err = permission_api.RestAllLevel()
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != 5 {
		t.Error(err)
		t.Error(level)
	}
}

func TestAddPermission(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != 1919810 {
		t.Error(err)
		t.Error(level)
	}
}

func TestRemovePermission(t *testing.T) {
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	err = permission_api.RemovePermission(114514)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != 5 {
		t.Error(err)
		t.Error(level)
	}
}

func TestUpdatePermission(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	err = permission_api.UpdatePermission(114514, 060303)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != 060303 {
		t.Error(err)
		t.Error(level)
	}
}

func TestQueryPermission(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != 1919810 {
		t.Error(err)
		t.Error(level)
	}
}

func TestCheckPermissions(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	result_1, err := permission_api.CheckPermissions(114514, 5)
	if err != nil {
		t.Error(err)
		t.Error(result_1)
	}
	if !result_1 {
		t.Errorf("Result_1 expected to be true, but it's false")
	}
	result_2, err := permission_api.CheckPermissions(114514, 1919810)
	if err != nil {
		t.Error(err)
		t.Error(result_2)
	}
	if !result_2 {
		t.Errorf("Result_2 expected to be true, but it's false")
	}
	result_3, err := permission_api.CheckPermissions(114514, 1919811)
	if err != nil {
		t.Error(err)
		t.Error(result_2)
	}
	if result_3 {
		t.Errorf("Result_3 expected to be false, but it's true")
	}
}

func TestResetUserLevel(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.ResetUserLevel(114514)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != 5 {
		t.Error(err)
		t.Error(level)
	}
}

func TestBanUser(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.BanUser(114514)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != -1 {
		t.Error(err)
		t.Error(level)
	}
	err = permission_api.BanUser(114514)
	if err == nil {
		t.Error(err)
	}
}

func TestIsBaned(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, -1)
	if err != nil {
		t.Error(err)
	}
	result_1, err := permission_api.IsBaned(114514)
	if err != nil || result_1 != true {
		t.Error(err)
		t.Error(result_1)
	}
	err = permission_api.UpdatePermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	result_2, err := permission_api.IsBaned(114514)
	if err != nil || result_2 != false {
		t.Error(err)
		t.Error(result_2)
	}
}

func TestUnbanUser(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}
	err = permission_api.BanUser(114514)
	if err != nil {
		t.Error(err)
	}
	level, err := permission_api.QueryPermission(114514)
	if err != nil || level != -1 {
		t.Error(err)
		t.Error(level)
	}
	err = permission_api.UnbanUser(114514)
	if err != nil {
		t.Error(err)
	}
	level, err = permission_api.QueryPermission(114514)
	if err != nil || level != 1919810 {
		t.Error(err)
		t.Error(level)
	}
	err = permission_api.UnbanUser(114514)
	if err == nil {
		t.Error(err)
	}
}

func TestQueryLastPermission(t *testing.T) {
	defer permission_api.RemovePermission(114514)
	err := permission_api.AddPermission(114514, 1919810)
	if err != nil {
		t.Error(err)
	}

	err = permission_api.UpdatePermission(114514, 5)
	if err != nil {
		t.Error(err)
	}
	last_level, err := permission_api.QueryLastPermission(114514)
	if err != nil{
		t.Error(err)
	}
	current_level, err := permission_api.QueryPermission(114514)
	if err != nil || current_level != 5 || last_level != 1919810 {
		t.Error(err)
		t.Error(current_level)
	}
}
