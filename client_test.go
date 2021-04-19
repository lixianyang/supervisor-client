package supervisor

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

var testClient *Client

func TestMain(m *testing.M) {
	var err error
	testClient, err = New("http://user:123@127.0.0.1:9001", 5*time.Second)
	if err != nil {
		fmt.Println("failed to create client", err)
		os.Exit(1)
	}
	infos, err := testClient.StartAllProcesses(true)
	if err != nil {
		fmt.Println("failed to start all processes", err)
	}
	for _, info := range infos {
		if info.Status != StatusSuccess && info.Status != StatusAlreadyStarted {
			fmt.Println("failed to start process", info.Name, info.Description)
			os.Exit(1)
		}
	}
	code := m.Run()
	os.Exit(code)
}

func TestListMethods(t *testing.T) {
	methods, err := testClient.ListMethods()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range methods {
		t.Log(m)
	}
	if l := len(methods); l == 0 {
		t.Fatal("invalid methods")
	}
}

func TestMethodHelp(t *testing.T) {
	help, err := testClient.MethodHelp("system.listMethods")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("system.listMethods:", help)
	if l := len(help); l == 0 {
		t.Fatal("invalid method help")
	}
}

func TestGetAPIVersion(t *testing.T) {
	version, err := testClient.GetAPIVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("APIVersion:", version)
	if l := len(version); l == 0 {
		t.Fatal("invalid api version")
	}
}

func TestGetSupervisorVersion(t *testing.T) {
	version, err := testClient.GetSupervisorVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("SupervisorVersion:", version)
	if l := len(version); l == 0 {
		t.Fatal("invalid api version")
	}
}

func TestGetIdentification(t *testing.T) {
	id, err := testClient.GetIdentification()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Identification:", id)
	if l := len(id); l == 0 {
		t.Fatal("invalid identification")
	}
}

func TestGetState(t *testing.T) {
	s, err := testClient.GetState()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("State code: %d name: %s", s.Code, s.Name)
	if s.Code != ServerFatal && s.Code != ServerRestarting && s.Code != ServerShutdown && s.Code != ServerRunning {
		t.Fatalf("invalid state code %d name %s", s.Code, s.Name)
	}
}

func TestGetPID(t *testing.T) {
	pid, err := testClient.GetPID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Supervisord pid:", pid)
	if pid == 0 {
		t.Fatal("invalid pid")
	}
}

func TestReadLog(t *testing.T) {
	content, err := testClient.ReadLog(0, 16)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Main log:", content)
	if len(content) == 0 {
		t.Fatal("invalid main log content")
	}
}

func TestStartStopProcess(t *testing.T) {
	name := "web"
	testExpectProcessState(t, name, ProcessRunning)
	err := testClient.StopProcess(name, true)
	if err != nil {
		t.Fatal(err)
	}
	testExpectProcessState(t, name, ProcessStopped)
	err = testClient.StartProcess(name, true)
	if err != nil {
		t.Fatal(err)
	}
	testExpectProcessState(t, name, ProcessRunning)
}

func testExpectProcessState(t *testing.T, name string, s ProcessState) {
	st, err := testClient.GetProcessInfo(name)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Process %s state %s", name, st.StateName)
	if st.State != s {
		t.Fatal("invalid process state", st.StateName, "expected", s, "but", st.State)
	}
}

func TestStartStopAllProcesses(t *testing.T) {
	testExpectedAllProcesses(t, ProcessRunning)
	testAllProcesses(t, testClient.StopAllProcesses)
	testExpectedAllProcesses(t, ProcessStopped)
	testAllProcesses(t, testClient.StartAllProcesses)
	testExpectedAllProcesses(t, ProcessRunning)
}

func testExpectedAllProcesses(t *testing.T, state ProcessState) {
	infos, err := testClient.GetAllProcessInfo()
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range infos {
		if info.State != state {
			t.Fatalf("Process %s:%s expected %d but %d", info.Group, info.Name, state, info.State)
		}
	}
}

func testAllProcesses(t *testing.T, f func(wait bool) ([]ActionStatus, error)) {
	infos, err := f(true)
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range infos {
		t.Logf("Process %s:%s %s", info.Group, info.Name, info.Description)
		if info.Status != StatusSuccess {
			t.Fatalf("Process %s:%s %s", info.Group, info.Name, info.Description)
		}
	}
}

func TestSignalProcess(t *testing.T) {
	name := "web"
	testExpectProcessState(t, name, ProcessRunning)
	err := testClient.SignalProcess(name, syscall.SIGINT)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	testExpectProcessState(t, name, ProcessExited)
	err = testClient.StartProcess(name, true)
	if err != nil {
		t.Fatal(err)
	}
	testExpectProcessState(t, name, ProcessRunning)
}

func TestGetProcessInfo(t *testing.T) {
	name := "web"
	info, err := testClient.GetProcessInfo(name)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Process info:\n%s", info)
	if info.State != ProcessRunning {
		t.Fatalf("Process %s:%s %s expected %d but %d", info.Group, info.Name, info.StateName, ProcessRunning, info.State)
	}
}

func TestGetAllProcessInfo(t *testing.T) {
	infos, err := testClient.GetAllProcessInfo()
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range infos {
		t.Logf("Process info:\n%s", info)
		if info.State != ProcessRunning {
			t.Fatalf("Process %s:%s %s expected %d but %d", info.Group, info.Name, info.StateName, ProcessRunning, info.State)
		}
	}
}

func TestReadProcessStdoutLog(t *testing.T) {
	name := "web"
	content, err := testClient.ReadProcessStdoutLog(name, 0, 16)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Process %s stdout log: %s", name, content)
	if len(content) == 0 {
		t.Fatal("invalid stdout log")
	}
}

func TestTailProcessStdoutLog(t *testing.T) {
	name := "web"
	result, err := testClient.TailProcessStdoutLog(name, 0, 16)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Tail process stdout log:\ncontent:%soffset:%d\noverflow:%t", result.Content, result.Offset, result.Overflow)
}

func TestReloadConfig(t *testing.T) {
	added, changed, removed, err := testClient.ReloadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(added) > 0 {
		t.Log("Added:", strings.Join(added, ","))
	}
	if len(changed) > 0 {
		t.Log("Changed:", strings.Join(changed, ","))
	}
	if len(removed) > 0 {
		t.Log("Removed:", strings.Join(removed, ","))
	}
}

func TestGetAllConfigInfo(t *testing.T) {
	configs, err := testClient.GetAllConfigInfo()
	if err != nil {
		t.Fatal(err)
	}
	for _, cfg := range configs {
		t.Logf("****************\n%s", cfg)
	}
}
