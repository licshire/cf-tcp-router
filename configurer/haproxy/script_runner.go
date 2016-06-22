package haproxy

//go:generate counterfeiter -o fakes/fake_script_runner.go . ScriptRunner
type ScriptRunner interface {
	Run() error
}
