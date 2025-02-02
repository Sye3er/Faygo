package commandimport (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)
type LinuxCommand struct {
}
func NewLinuxCommand() *LinuxCommand {
	return &LinuxCommand{}
}
func (lc *LinuxCommand) Exec(args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command(os.Getenv("SHELL"), args...)	cmd.SysProcAttr = &syscall.SysProcAttr{}	outpip, err := cmd.StdoutPipe()
	defer outpip.Close()	if err != nil {
		return 0, "", err
	}	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}	return cmd.Process.Pid, string(out), nil
}
func (lc *LinuxCommand) ExecAsync(stdout chan string, args ...string) int {
	var pidChan = make(chan int, 1)	go func() {
		args = append([]string{"-c"}, args...)
		cmd := exec.Command(os.Getenv("SHELL"), args...)		cmd.SysProcAttr = &syscall.SysProcAttr{}		outpip, err := cmd.StdoutPipe()
		defer outpip.Close()		if err != nil {
			panic(err)
		}		err = cmd.Start()
		if err != nil {
			panic(err)
		}		pidChan <- cmd.Process.Pid		out, err := ioutil.ReadAll(outpip)
		if err != nil {
			panic(err)
		}		stdout <- string(out)
	}()	return <-pidChan
}
func (lc *LinuxCommand) ExecIgnoreResult(args ...string) error {	args = append([]string{"-c"}, args...)
	cmd := exec.Command(os.Getenv("SHELL"), args...)	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}	err := cmd.Run()	return err
}
