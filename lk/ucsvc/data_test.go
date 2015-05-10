package main

import "testing"

func TestNewID(t *testing.T) {
	id1 := newID()
	id2 := newID()

	if len(id1) != 36 {
		t.Errorf("id1 [%s] unexpected length of %d, expected 36", id1, len(id1))
	}
	if id1 == id2 {
		t.Errorf("id1 [%s] should not be == id2 [%s]", id1, id2)
	}
}

func TestID_IsValid(t *testing.T) {
	testCases := []struct {
		ID       ID
		Expected bool
	}{
		{emptyID, false},
		{ID("f47ac10b-58cc-0372-8567-0e02b2c3d479-a"), false},
		{ID("f47ac10b-58cc-0372-8567-0e02b2c3d479"), true},
		{newID(), true},
	}

	for _, testCase := range testCases {
		actual := testCase.ID.IsValid()
		if actual != testCase.Expected {
			t.Errorf("For %s expected %t, actual %t", testCase.ID, testCase.Expected, actual)
		}
	}
}

func TestNewUser_IsValid(t *testing.T) {
	testCases := []struct {
		User     *newUser
		Expected bool
	}{
		{&newUser{Name: ""}, false},
		{&newUser{Name: "  "}, false},
		{&newUser{Name: "ted"}, true},
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
		User     *user
		Expected bool
	}{
		{&user{ID: emptyID, Name: ""}, false},
		{&user{ID: ID("f47ac10b-58cc-0372-8567-0e02b2c3d479-a"), Name: ""}, false},
		{&user{ID: ID("f47ac10b-58cc-0372-8567-0e02b2c3d479"), Name: "  "}, false},
		{&user{ID: newID(), Name: "ted"}, true},
	}

	for _, testCase := range testCases {
		actual := testCase.User.IsValid()
		if actual != testCase.Expected {
			t.Errorf("For %#v expected %t, actual %t", testCase.User, testCase.Expected, actual)
		}
	}
}
