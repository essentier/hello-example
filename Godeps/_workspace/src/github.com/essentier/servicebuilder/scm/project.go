package scm

// This interface represents a source project with version control.
// It knows how to push code to its corresponding repository in the version control system.
type Project interface {
	PushCode(repoUrl string) error
}
