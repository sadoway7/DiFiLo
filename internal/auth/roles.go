package auth

// User roles.
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleGeneral = "general"
)

// UserRole provides the role and identity of a user for permission
// checks. It is satisfied by user types in the db package, which keeps
// this leaf package free of any dependency on them.
type UserRole interface {
	GetRole() string
	GetID() int64
}

// CommentOwner provides the owning user of a comment for permission
// checks. It is satisfied by comment types in the db package.
type CommentOwner interface {
	GetUserID() int64
}

// CanDeleteComment reports whether the user may delete the given
// comment. Admins and managers may delete any comment; general users
// may delete only their own. A nil user or nil comment is never
// permitted.
func CanDeleteComment(user UserRole, comment CommentOwner) bool {
	if user == nil || comment == nil {
		return false
	}
	switch user.GetRole() {
	case RoleAdmin, RoleManager:
		return true
	case RoleGeneral:
		return user.GetID() == comment.GetUserID()
	}
	return false
}
