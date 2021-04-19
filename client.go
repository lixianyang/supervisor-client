package supervisor

import (
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/kolo/xmlrpc"
)

type Client struct {
	url    string
	client *xmlrpc.Client
}

// New Create new supervisor xml rpc client
func New(url string, timeout time.Duration) (*Client, error) {
	transport := &http.Transport{ResponseHeaderTimeout: timeout}
	rpcClient, err := xmlrpc.NewClient(url+"/RPC2", transport)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		client: rpcClient,
		url:    url,
	}
	return cli, nil
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// ListMethods Return an array listing the available method names
func (c *Client) ListMethods() ([]string, error) {
	methods := make([]string, 0)
	err := c.call(SystemNamespace, "listMethods", nil, &methods)
	return methods, err
}

// MethodHelp Return a string showing the method's documentation
func (c *Client) MethodHelp(name string) (string, error) {
	var help string
	args := []interface{}{name}
	err := c.call(SystemNamespace, "methodHelp", args, &help)
	return help, err
}

// GetAPIVersion Return the version of the RPC API used by supervisord
func (c *Client) GetAPIVersion() (string, error) {
	var version string
	err := c.call(DefaultNamespace, "getAPIVersion", nil, &version)
	return version, err
}

// GetSupervisorVersion Return the version of the supervisor package in use by supervisord
func (c *Client) GetSupervisorVersion() (string, error) {
	var version string
	err := c.call(DefaultNamespace, "getSupervisorVersion", nil, &version)
	return version, err
}

// GetIdentification Return identifying string of supervisord
func (c *Client) GetIdentification() (string, error) {
	var identifier string
	err := c.call(DefaultNamespace, "getIdentification", nil, &identifier)
	return identifier, err
}

// GetState Return current state of supervisord as a struct
func (c *Client) GetState() (ServerState, error) {
	var state ServerState
	err := c.call(DefaultNamespace, "getState", nil, &state)
	return state, err
}

// GetPID Return the PID of supervisord
func (c *Client) GetPID() (int, error) {
	var pid int
	err := c.call(DefaultNamespace, "getPID", nil, &pid)
	return pid, err
}

// ReadLog Read length bytes from the main log starting at offset
func (c *Client) ReadLog(offset, length int) (string, error) {
	var content string
	args := []interface{}{offset, length}
	err := c.call(DefaultNamespace, "readLog", args, &content)
	return content, err
}

// ClearLog Clear the main log
func (c *Client) ClearLog() (bool, error) {
	var flag bool
	err := c.call(DefaultNamespace, "clearLog", nil, &flag)
	return flag, err
}

// Shutdown Shut down the supervisor process
func (c *Client) Shutdown() error {
	var flag bool
	err := c.call(DefaultNamespace, "shutdown", nil, &flag)
	return err
}

// Restart Restart the supervisor process
func (c *Client) Restart() error {
	var flag bool
	err := c.call(DefaultNamespace, "restart", nil, &flag)
	return err
}

// ReloadConfig Reload the configuration
func (c *Client) ReloadConfig() (added []string, changed []string, removed []string, err error) {
	result := make([][][]string, 0)
	err = c.call(DefaultNamespace, "reloadConfig", nil, &result)
	if err != nil {
		return nil, nil, nil, err
	}
	return result[0][0], result[0][1], result[0][2], nil
}

// AddProcessGroup Update the config for a running process from config file
func (c *Client) AddProcessGroup(name string) (bool, error) {
	var flag bool
	err := c.call(DefaultNamespace, "addProcessGroup", name, &flag)
	return flag, err
}

// RemoveProcessGroup Remove a stopped process from the active configuration
func (c *Client) RemoveProcessGroup(name string) (bool, error) {
	var flag bool
	err := c.call(DefaultNamespace, "removeProcessGroup", name, &flag)
	return flag, err
}

// StartProcess Start a process
// string name Process name (or ``group:name``, or ``group:*``)
// bool wait Wait for process to be fully started
func (c *Client) StartProcess(name string, wait bool) error {
	var flag bool
	args := []interface{}{name, wait}
	err := c.call(DefaultNamespace, "startProcess", args, &flag)
	return err
}

// StartProcessGroup Start all processes in the group named 'name'
// string name The group name
// bool wait Wait for process to be fully started
func (c *Client) StartProcessGroup(name string, wait bool) ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	args := []interface{}{name, wait}
	err := c.call(DefaultNamespace, "startProcessGroup", args, &infos)
	return infos, err
}

// StartAllProcesses Start all processes listed in the configuration file
func (c *Client) StartAllProcesses(wait bool) ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	args := []interface{}{wait}
	err := c.call(DefaultNamespace, "startAllProcesses", args, &infos)
	return infos, err
}

// StopProcess Stop a process named by name
// string name Process name (or ``group:name``, or ``group:*``)
// bool wait Wait for process to be fully stopped
func (c *Client) StopProcess(name string, wait bool) error {
	var flag bool
	args := []interface{}{name, wait}
	err := c.call(DefaultNamespace, "stopProcess", args, &flag)
	return err
}

// StopProcessGroup Stop all processes in the group named 'name'
// string name The group name
// bool wait Wait for process to be fully stopped
func (c *Client) StopProcessGroup(name string, wait bool) ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	args := []interface{}{name, wait}
	err := c.call(DefaultNamespace, "stopProcessGroup", args, &infos)
	return infos, err
}

// StopAllProcesses Stop all processes listed in the configuration file
func (c *Client) StopAllProcesses(wait bool) ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	args := []interface{}{wait}
	err := c.call(DefaultNamespace, "stopAllProcesses", args, &infos)
	return infos, err
}

// SignalProcess Send an arbitrary UNIX signal to the process named by name
func (c *Client) SignalProcess(name string, signal syscall.Signal) error {
	var flag bool
	args := []interface{}{name, signal}
	err := c.call(DefaultNamespace, "signalProcess", args, &flag)
	return err
}

// SignalProcessGroup Send a signal to all processes in the group named 'name'
func (c *Client) SignalProcessGroup(name string, signal syscall.Signal) ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	args := []interface{}{name, signal}
	err := c.call(DefaultNamespace, "signalProcessGroup", args, &infos)
	return infos, err
}

// SignalAllProcesses Send a signal to all processes in the process list
func (c *Client) SignalAllProcesses(signal syscall.Signal) ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	args := []interface{}{signal}
	err := c.call(DefaultNamespace, "signalAllProcesses", args, &infos)
	return infos, err
}

// GetAllConfigInfo Get info about all available process configurations. Each struct represents a single process (i.e. groups get flattened).
func (c *Client) GetAllConfigInfo() ([]ProgramConfig, error) { // fixme: should return struct, not interface
	configs := make([]ProgramConfig, 0)
	err := c.call(DefaultNamespace, "getAllConfigInfo", nil, &configs)
	return configs, err
}

// GetProcessInfo Get info about a process named name
func (c *Client) GetProcessInfo(name string) (ProcessInfo, error) {
	info := ProcessInfo{}
	args := []interface{}{name}
	err := c.call(DefaultNamespace, "getProcessInfo", args, &info)
	return info, err
}

// GetAllProcessInfo Get info about all processes
func (c *Client) GetAllProcessInfo() ([]ProcessInfo, error) {
	list := make([]ProcessInfo, 0)
	err := c.call(DefaultNamespace, "getAllProcessInfo", nil, &list)
	return list, err
}

// ReadProcessStdoutLog Read length bytes from name's stdout log starting at offset
func (c *Client) ReadProcessStdoutLog(name string, offset, length int) (string, error) {
	var content string
	args := []interface{}{name, offset, length}
	err := c.call(DefaultNamespace, "readProcessStdoutLog", args, &content)
	return content, err
}

// ReadProcessStderrLog Read length bytes from name's stderr log starting at offset
func (c *Client) ReadProcessStderrLog(name string, offset, length int) (string, error) {
	var content string
	args := []interface{}{name, offset, length}
	err := c.call(DefaultNamespace, "readProcessStderrLog", args, &content)
	return content, err
}

// TailProcessStdoutLog Provides a more efficient way to tail the (stdout) log than ReadProcessStdoutLog().  Use ReadProcessStdoutLog() to read chunks and TailProcessStdoutLog() to tail.
/*
   Requests (length) bytes from the (name)'s log, starting at
   (offset).  If the total log size is greater than (offset +
   length), the overflow flag is set and the (offset) is
   automatically increased to position the buffer at the end of
   the log.  If less than (length) bytes are available, the
   maximum number of available bytes will be returned.  (offset)
   returned is always the last offset in the log +1.
*/
func (c *Client) TailProcessStdoutLog(name string, offset, length int) (*TailResult, error) {
	result := make([]interface{}, 0, 3)
	args := []interface{}{name, offset, length}
	err := c.call(DefaultNamespace, "tailProcessStdoutLog", args, &result)
	if err != nil {
		return nil, err
	}
	tail := &TailResult{
		Content:  result[0].(string),
		Offset:   result[1].(int64),
		Overflow: result[2].(bool),
	}
	return tail, nil
}

// TailProcessStderrLog Provides a more efficient way to tail the (stderr) log than ReadProcessStderrLog().  Use ReadProcessStderrLog() to read chunks and TailProcessStderrLog() to tail.
/*
   Requests (length) bytes from the (name)'s log, starting at
   (offset).  If the total log size is greater than (offset +
   length), the overflow flag is set and the (offset) is
   automatically increased to position the buffer at the end of
   the log.  If less than (length) bytes are available, the
   maximum number of available bytes will be returned.  (offset)
   returned is always the last offset in the log +1.
*/
func (c *Client) TailProcessStderrLog(name string, offset, length int) (*TailResult, error) {
	result := make([]interface{}, 0, 3)
	args := []interface{}{name, offset, length}
	err := c.call(DefaultNamespace, "tailProcessStderrLog", args, &result)
	if err != nil {
		return nil, err
	}
	tail := &TailResult{
		Content:  result[0].(string),
		Offset:   result[1].(int64),
		Overflow: result[2].(bool),
	}
	return tail, nil
}

// ClearProcessLogs Clear the stdout and stderr logs for the named process and reopen them.
func (c *Client) ClearProcessLogs(name string) error {
	var flag bool
	args := []interface{}{name}
	err := c.call(DefaultNamespace, "clearProcessLogs", args, &flag)
	return err
}

// ClearAllProcessLogs Clear all process log files
func (c *Client) ClearAllProcessLogs() ([]ActionStatus, error) {
	infos := make([]ActionStatus, 0)
	err := c.call(DefaultNamespace, "clearAllProcessLogs", nil, &infos)
	return infos, err
}

// SendProcessStdin Send a string of chars to the stdin of the process name.
//        If non-7-bit data is sent (unicode), it is encoded to utf-8
//        before being sent to the process' stdin.  If chars is not a
//        string or is not unicode, return ErrIncorrectParameters.  If the
//        process is not running, return ErrNotRunning.  If the process'
//        stdin cannot accept input (e.g. it was closed by the child
//        process), return ErrNoFile.
func (c *Client) SendProcessStdin(name, chars string) error {
	var flag bool
	args := []interface{}{name, chars}
	err := c.call(DefaultNamespace, "sendProcessStdin", args, &flag)
	return err
}

func (c *Client) call(ns Namespace, method string, args interface{}, relay interface{}) error {
	return c.client.Call(fmt.Sprintf("%s.%s", ns, method), args, relay)
}
