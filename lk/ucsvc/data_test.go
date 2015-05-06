package main

import "testing"

func TestNewId(t *testing.T) {
	id1 := NewId()
	id2 := NewId()

	if len(id1) != 36 {
		t.Errorf("id1 [%s] unexpected length of %d, expected 36", id1, len(id1))
	}
	if id1 == id2 {
		t.Errorf("id1 [%s] should not be == id2 [%s]", id1, id2)
	}
}

func TestId_IsValid(t *testing.T) {
	testCases := []struct {
		Id       Id
		Expected bool
	}{
		{Id(""), false},
		{Id("f47ac10b-58cc-0372-8567-0e02b2c3d479-a"), false},
		{Id("f47ac10b-58cc-0372-8567-0e02b2c3d479"), true},
		{NewId(), true},
	}

	for _, testCase := range testCases {
		actual := testCase.Id.IsValid()
		if actual != testCase.Expected {
			t.Errorf("For %s expected %t, actual %t", testCase.Id, testCase.Expected, actual)
		}
	}
}

func TestNewUser_IsValid(t *testing.T) {
	testCases := []struct {
		User     *NewUser
		Expected bool
	}{
		{&NewUser{Name: ""}, false},
		{&NewUser{Name: "  "}, false},
		{&NewUser{Name: "ted"}, true},
	}

	for _, testCase := range testCases {
		actual := testCase.User.IsValid()
		if actual != testCase.Expected {
			t.Errorf("For %#v expected %t, actual %t", testCase.User, testCase.Expected, actual)
		}
	}
}

func TestUser_IsValid(t *testing.T) {
	testCases := []struct {
		User     *User
		Expected bool
	}{
		{&User{Id: Id(""), Name: ""}, false},
		{&User{Id: Id("f47ac10b-58cc-0372-8567-0e02b2c3d479-a"), Name: ""}, false},
		{&User{Id: Id("f47ac10b-58cc-0372-8567-0e02b2c3d479"), Name: "  "}, false},
		{&User{Id: NewId(), Name: "ted"}, true},
	}

	for _, testCase := range testCases {
		actual := testCase.User.IsValid()
		if actual != testCase.Expected {
			t.Errorf("For %#v expected %t, actual %t", testCase.User, testCase.Expected, actual)
		}
	}
}
