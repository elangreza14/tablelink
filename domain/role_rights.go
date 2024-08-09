package domain

type RoleRights struct {
	ID      int
	RoleID  int
	Route   string
	Section string
	Path    string
	RCreate int
	RRead   int
	RUpdate int
	RDelete int
}
