package tests

import (
	"testing"

	"github.com/DaydreamCafe/Cocoa/plugins/nickname"
)

func init() {
	nickname.ResetAllNicknames()
}

func TestAddNickname(t *testing.T) {
	defer nickname.ResetAllNicknames()
	err := nickname.AddNickname(114514, "田所浩二", "野兽先辈")
	if err != nil {
		t.Error(err)
	}
	nickName, err := nickname.QueryNickname(114514)
	if err != nil {
		t.Error(err)
	}
	if nickName != "野兽先辈" {
		t.Error(err)
	}
}

func TestUpdateNickname(t *testing.T) {
	defer nickname.ResetAllNicknames()

	err := nickname.UpdateNickname(114514, "田所浩二", "野兽先辈")
	if err != nil {
		t.Error(err)
	}
	nickName, err := nickname.QueryNickname(114514)
	if err != nil {
		t.Error(err)
	}
	if nickName != "野兽先辈" {
		t.Error(err)
	}

	err = nickname.UpdateNickname(114514, "田所浩二", "李田所")
	if err != nil {
		t.Error(err)
	}
	nickName, err = nickname.QueryNickname(114514)
	if err != nil {
		t.Error(err)
	}
	if nickName != "李田所" {
		t.Error(err)
	}
}

func TestRemoveNickname(t *testing.T) {
	defer nickname.ResetAllNicknames()

	err := nickname.AddNickname(114514, "田所浩二", "野兽先辈")
	if err != nil {
		t.Error(err)
	}
	nickName, err := nickname.QueryNickname(114514)
	if err != nil {
		t.Error(err)
	}
	if nickName != "野兽先辈" {
		t.Error(err)
	}

	err = nickname.RemoveNickname(114514)
	if err != nil {
		t.Error(err)
	}
	nickName, err = nickname.QueryNickname(114514)
	if err == nil {
		t.Error()
	}
}

func TestQueryNickname(t *testing.T) {
	defer nickname.ResetAllNicknames()
	err := nickname.AddNickname(114514, "田所浩二", "野兽先辈")
	if err != nil {
		t.Error(err)
	}
	nickName, err := nickname.QueryNickname(114514)
	if err != nil {
		t.Error(err)
	}
	if nickName != "野兽先辈" {
		t.Error(err)
	}
}
