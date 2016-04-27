package git

import (
	"log"
	"strings"

	"github.com/essentier/nomockutil/cmd"
	"github.com/essentier/servicebuilder/scm"
	"github.com/go-errors/errors"
)

func CreateDefaultGitProject(projectDir string) scm.Project {
	return &gitProject{cmdRunner: CreateGitCmdRunner(projectDir, cmd.CreateFailFastCmdRunner(cmd.CreateCmdConsole()))}
}

func CreateGitProject(projectDir string, cmdRunner cmd.CmdRunner) scm.Project {
	return &gitProject{cmdRunner: CreateGitCmdRunner(projectDir, cmdRunner)}
}

// We will have a gitProject instance for each source service.
// Not safe for concurrent access.
type gitProject struct {
	cmdRunner cmd.CmdRunner
}

func (g *gitProject) PushCode(repoUrl string) error {
	g.cmdRunner.SetError(cmd.CmdError{})
	stashCount := 0
	originalBranch := g.getCurrentBranch()
	stashCount = g.stashAll(stashCount)
	if !g.cmdRunner.HasError() {
		defer g.deferredPopStashed(stashCount)
	}

	g.branch("nomock")
	if !g.cmdRunner.HasError() {
		defer g.deferredDeleteBranch("nomock")
	}

	g.checkout("nomock")
	if !g.cmdRunner.HasError() {
		defer g.deferredCheckout(originalBranch)
	}

	g.pull(repoUrl, "nomock")
	g.applyStash(stashCount)
	g.addAll()
	g.commit("'done by nomock'")
	g.push(repoUrl, "nomock")
	lastError := g.cmdRunner.LastError()
	if lastError.Err != nil {
		return lastError
	} else {
		return nil
	}
}

func (g *gitProject) init() {
	g.cmdRunner.RunCmd("init")
}

func (g *gitProject) pull(remoteUrl string, branchName string) {
	if g.cmdRunner.HasError() {
		return
	}

	// This will fail if the project is pushed to essentier nomock the first time. That is okay.
	g.cmdRunner.RunInNewRunner("pull", "-s", "ours", remoteUrl, branchName)
}

func (g *gitProject) applyStash(stashCount int) {
	if stashCount <= 0 {
		return
	}

	g.cmdRunner.RunCmd("stash", "apply")
}

func (g *gitProject) checkout(branchName string) {
	g.cmdRunner.RunCmd("checkout", branchName)
}

func (g *gitProject) branch(branchName string) {
	g.cmdRunner.RunCmd("branch", branchName)
}

func (g *gitProject) addAll() {
	g.cmdRunner.RunCmd("add", ".", "-A")
}

func (g *gitProject) push(remoteUrl string, branchName string) {
	g.cmdRunner.RunCmd("push", remoteUrl, branchName)
}

func (g *gitProject) commit(message string) {
	if g.cmdRunner.HasError() {
		return
	}

	stdout, cmdErr := g.cmdRunner.RunInNewRunner("commit", "-m", message)
	if cmdErr.Err != nil {
		if !strings.Contains(stdout, "clean") {
			g.cmdRunner.SetError(cmdErr)
		}
	}
}

func (g *gitProject) stashAll(stashCount int) int {
	stdout, _ := g.cmdRunner.RunCmd("stash", "save", "-u")
	if g.cmdRunner.HasError() {
		return stashCount
	}

	if strings.Contains(stdout, "HEAD") {
		return stashCount + 1
	} else {
		return stashCount
	}
}

func parseForCurrentBranch(stdout string) string {
	currentBranch := ""
	branchs := strings.Split(stdout, "\n")
	for _, branch := range branchs {
		if strings.HasPrefix(branch, "*") {
			currentBranch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
			break
		}
	}
	log.Printf("current branch: [ %v ]", currentBranch)
	return currentBranch
}

func (g *gitProject) getCurrentBranch() string {
	stdout, _ := g.cmdRunner.RunCmd("branch")
	if g.cmdRunner.HasError() {
		return ""
	}

	currentBranch := parseForCurrentBranch(stdout)
	if currentBranch == "" {
		g.cmdRunner.SetError(cmd.CmdError{Err: errors.Errorf("Failed to find current git branch.")})
	}
	return currentBranch
}

func (g *gitProject) deferredPopStashed(stashCount int) int {
	if stashCount <= 0 {
		return stashCount
	}

	_, cmdErr := g.cmdRunner.RunInNewRunner("stash", "pop")
	if cmdErr.Err != nil {
		logCmdErrorIfAny(cmdErr, "Error when trying to pop stashed.")
		return stashCount
	} else {
		return stashCount - 1
	}
}

func (g *gitProject) deferredDeleteBranch(branchName string) {
	_, cmdErr := g.cmdRunner.RunInNewRunner("branch", "-D", branchName)
	logCmdErrorIfAny(cmdErr, "Error when trying to delete the nomock branch.")
}

func (g *gitProject) deferredCheckout(originalBranch string) {
	_, cmdErr := g.cmdRunner.RunInNewRunner("checkout", originalBranch)
	logCmdErrorIfAny(cmdErr, "Error when trying to checkout the original branch.")
}

func logCmdErrorIfAny(cmdErr cmd.CmdError, message string) {
	if cmdErr.Err != nil {
		log.Printf("%v Error: %v", message, cmdErr.Error())
	}
}
