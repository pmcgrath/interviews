package main

func authenticate(userName, password string) bool {
	// Should be going to a store, but this is good enough for now
	return userName == "ted" && password == "toe"
}
