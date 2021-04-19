package supervisor

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

type Namespace string

const (
	SystemNamespace  Namespace = "system"
	DefaultNamespace Namespace = "supervisor"
)

type Status int

const (
	StatusUnknownMethod        Status = 1  // UNKNOWN_METHOD
	StatusIncorrectParameters  Status = 2  // INCORRECT_PARAMETERS
	StatusBadArguments         Status = 3  // BAD_ARGUMENTS
	StatusSignatureUnsupported Status = 4  // SIGNATURE_UNSUPPORTED
	StatusShutdownState        Status = 6  // SHUTDOWN_STATE
	StatusBadName              Status = 10 // BAD_NAME
	StatusBadSignal            Status = 11 // BAD_SIGNAL
	StatusNoFile               Status = 20 // NO_FILE
	StatusNotExecutable        Status = 21 // NOT_EXECUTABLE
	StatusFailed               Status = 30 // FAILED
	StatusAbnormalTermination  Status = 40 // ABNORMAL_TERMINATION
	StatusSpawnError           Status = 50 // SPAWN_ERROR
	StatusAlreadyStarted       Status = 60 // ALREADY_STARTED
	StatusNotRunning           Status = 70 // NOT_RUNNING
	StatusSuccess              Status = 80 // SUCCESS
	StatusAlreadyAdded         Status = 90 // ALREADY_ADDED
	StatusStillRunning         Status = 91 // STILL_RUNNING
	StatusCantReread           Status = 92 // CANT_REREAD
)

// ProcessState process state
type ProcessState int

const (
	ProcessStopped  ProcessState = 0    // STOPPED
	ProcessStarting ProcessState = 10   // STARTING
	ProcessRunning  ProcessState = 20   // RUNNING
	ProcessBackoff  ProcessState = 30   // BACKOFF
	ProcessStopping ProcessState = 40   // STOPPING
	ProcessExited   ProcessState = 100  // EXITED
	ProcessFatal    ProcessState = 200  // FATAL
	ProcessUnknown  ProcessState = 1000 // UNKNOWN
)

type State int

const (
	ServerFatal      State = 2  // FATAL
	ServerRunning    State = 1  // RUNNING
	ServerRestarting State = 0  // RESTARTING
	ServerShutdown   State = -1 // SHUTDOWN
)

var (
	ErrIncorrectParameters = errors.New("INCORRECT_PARAMETERS")
	ErrNotRunning          = errors.New("NOT_RUNNING")
	ErrNoFile              = errors.New("NO_FILE")
)

type ServerState struct {
	Code State  `xmlrpc:"statecode"`
	Name string `xmlrpc:"statename"`
}

type ActionStatus struct {
	Name        string `xmlrpc:"name"`
	Group       string `xmlrpc:"group"`
	Status      Status `xmlrpc:"status"`
	Description string `xmlrpc:"description"`
}

type ProcessInfo struct {
	Name          string       `xmlrpc:"name"`
	Group         string       `xmlrpc:"group"`
	Start         int          `xmlrpc:"start"`
	Stop          int          `xmlrpc:"stop"`
	Now           int          `xmlrpc:"now"`
	State         ProcessState `xmlrpc:"state"`
	StateName     string       `xmlrpc:"statename"`
	SpawnErr      string       `xmlrpc:"spawnerr"`
	ExitStatus    int          `xmlrpc:"exitstatus"`
	Logfile       string       `xmlrpc:"logfile"`
	StdoutLogfile string       `xmlrpc:"stdout_logfile"`
	StderrLogfile string       `xmlrpc:"stderr_logfile"`
	Pid           int          `xmlrpc:"pid"`
}

func (pi ProcessInfo) String() string {
	return printStruct(pi)
}

type TailResult struct {
	Content  string
	Offset   int64
	Overflow bool
}

type ProgramConfig struct {
	Name                  string         `xmlrpc:"name"`
	Group                 string         `xmlrpc:"group"`
	Command               string         `xmlrpc:"command"`
	InUse                 bool           `xmlrpc:"inuse"`
	Autostart             bool           `xmlrpc:"autostart"`
	StartSeconds          int            `xmlrpc:"startsecs"`
	StartRetries          int            `xmlrpc:"startretries"`
	StopSignal            syscall.Signal `xmlrpc:"stopsignal"`
	StopWaitSeconds       int            `xmlrpc:"stopwaitsecs"`
	RedirectStderr        bool           `xmlrpc:"redirect_stderr"`
	ExitCodes             []int          `xmlrpc:"exitcodes"`
	ProcessPriority       int            `xmlrpc:"process_prio"`
	GroupPriority         int            `xmlrpc:"group_prio"`
	KillAsGroup           bool           `xmlrpc:"killasgroup"`
	StdoutLogfile         string         `xmlrpc:"stdout_logfile"`
	StderrLogfile         string         `xmlrpc:"stderr_logfile"`
	StderrLogfileBackups  int            `xmlrpc:"stderr_logfile_backups"`
	StdoutLogfileBackups  int            `xmlrpc:"stdout_logfile_backups"`
	StdoutLogfileMaxBytes int64          `xmlrpc:"stdout_logfile_maxbytes"`
	StderrLogfileMaxBytes int64          `xmlrpc:"stderr_logfile_maxbytes"`
	StdoutCaptureMaxBytes int64          `xmlrpc:"stdout_capture_maxbytes"`
	StderrCaptureMaxBytes int64          `xmlrpc:"stderr_capture_maxbytes"`
	StdoutEventsEnabled   bool           `xmlrpc:"stdout_events_enabled"`
	StderrEventsEnabled   bool           `xmlrpc:"stderr_events_enabled"`
}

func (pc ProgramConfig) String() string {
	return printStruct(pc)
}

func printStruct(s interface{}) string {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	b := strings.Builder{}
	maxLength := 0
	for i := 0; i < v.NumField(); i++ {
		if v.CanInterface() {
			if l := len(t.Field(i).Name); l > maxLength {
				maxLength = l
			}
		}
	}
	kvFormat := "%-" + strconv.Itoa(maxLength) + "s: %v\n"
	for i := 0; i < v.NumField(); i++ {
		if v.CanInterface() {
			b.WriteString(fmt.Sprintf(kvFormat, t.Field(i).Name, v.Field(i)))
		}
	}
	return b.String()
}
