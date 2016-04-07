package git

import "github.com/essentier/nomockutil/cmd"

func CreateGitCmdRunner(projectDir string, cmdRunner cmd.CmdRunner) cmd.CmdRunner {
	return &gitCmdRunner{
		CmdRunner:  cmdRunner,
		projectDir: projectDir,
	}
}

// not safe for concurrent access
type gitCmdRunner struct {
	cmd.CmdRunner
	projectDir string
}

func (g *gitCmdRunner) RunInNewRunner(name string, args ...string) (string, cmd.CmdError) {
	return g.CmdRunner.RunInNewRunner("git", g.appendGitArgs(name, args)...)
}

func (g *gitCmdRunner) RunCmd(name string, args ...string) (string, cmd.CmdError) {
	return g.CmdRunner.RunCmd("git", g.appendGitArgs(name, args)...)
}

func (g *gitCmdRunner) appendGitArgs(name string, args []string) []string {
	return append([]string{"-C", g.projectDir, name}, args...)
}
