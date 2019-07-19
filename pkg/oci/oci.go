package oci

import (
	"bytes"
	"fmt"
	"io"
	"syscall"
	"time"

	spec "github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/net/context"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	// ContainerStateCreated represents the created state of a container
	ContainerStateCreated = "created"
	// ContainerStatePaused represents the paused state of a container
	ContainerStatePaused = "paused"
	// ContainerStateRunning represents the running state of a container
	ContainerStateRunning = "running"
	// ContainerStateStopped represents the stopped state of a container
	ContainerStateStopped = "stopped"
	// ContainerCreateTimeout represents the value of container creating timeout
	ContainerCreateTimeout = 240 * time.Second

	// CgroupfsCgroupsManager represents cgroupfs native cgroup manager
	CgroupfsCgroupsManager = "cgroupfs"
	// SystemdCgroupsManager represents systemd native cgroup manager
	SystemdCgroupsManager = "systemd"

	// killContainerTimeout is the timeout that we wait for the container to
	// be SIGKILLed.
	KillContainerTimeout = 2 * time.Minute

	// minCtrStopTimeout is the minimal amount of time in seconds to wait
	// before issuing a timeout regarding the proper termination of the
	// container.
	minCtrStopTimeout = 30
)

// ExecSyncError wraps command's streams, exit code and error on ExecSync error.
type ExecSyncError struct {
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int32
	Err      error
}

func (e *ExecSyncError) Error() string {
	return fmt.Sprintf("command error: %+v, stdout: %s, stderr: %s, exit code %d", e.Err, e.Stdout.Bytes(), e.Stderr.Bytes(), e.ExitCode)
}

// RuntimeImpl is an interface used by the caller to interact with the
// container runtime. The purpose of this interface being to abstract
// implementations and their associated assumptions regarding the way to
// interact with containers. This will allow for new implementations of
// this interface, especially useful for the case of VM based container
// runtimes. Assumptions based on the fact that a container process runs
// on the host will be limited to the RuntimeOCI implementation.
type RuntimeImpl interface {
	CreateContainer(*Container, string) error
	StartContainer(*Container) error
	ExecContainer(*Container, []string, io.Reader, io.WriteCloser, io.WriteCloser,
		bool, <-chan remotecommand.TerminalSize) error
	ExecSyncContainer(*Container, []string, int64) (*ExecSyncResponse, error)
	UpdateContainer(*Container, *spec.LinuxResources) error
	StopContainer(context.Context, *Container, int64) error
	DeleteContainer(*Container) error
	UpdateContainerStatus(*Container) error
	PauseContainer(*Container) error
	UnpauseContainer(*Container) error
	ContainerStats(*Container) (*ContainerStats, error)
	SignalContainer(*Container, syscall.Signal) error
	AttachContainer(*Container, io.Reader, io.WriteCloser, io.WriteCloser, bool, <-chan remotecommand.TerminalSize) error
	PortForwardContainer(*Container, int32, io.ReadWriter) error
	ReopenContainerLog(*Container) error
	WaitContainerStateStopped(context.Context, *Container) error
}

// ContainerStats contains the statistics information for a running container
type ContainerStats struct {
	Container   string
	CPU         float64
	CPUNano     uint64
	SystemNano  int64
	MemUsage    uint64
	MemLimit    uint64
	MemPerc     float64
	NetInput    uint64
	NetOutput   uint64
	BlockInput  uint64
	BlockOutput uint64
	PIDs        uint64
}

// ExecSyncResponse is returned from ExecSync.
type ExecSyncResponse struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int32
}

//func MetricsToCtrStats(c *oci.Container, m *cgroups.Metrics) *oci.ContainerStats {
//	var (
//		cpu         float64
//		cpuNano     uint64
//		memUsage    uint64
//		memLimit    uint64
//		memPerc     float64
//		netInput    uint64
//		netOutput   uint64
//		blockInput  uint64
//		blockOutput uint64
//		pids        uint64
//	)
//
//	if m != nil {
//		pids = m.Pids.Current
//
//		cpuNano = m.CPU.Usage.Total
//		cpu = genericCalculateCPUPercent(cpuNano, m.CPU.Usage.PerCPU)
//
//		memUsage = m.Memory.Usage.Usage
//		memLimit = getMemLimit(m.Memory.Usage.Limit)
//		memPerc = float64(memUsage) / float64(memLimit)
//
//		for _, entry := range m.Blkio.IoServiceBytesRecursive {
//			switch strings.ToLower(entry.Op) {
//			case "read":
//				blockInput += entry.Value
//			case "write":
//				blockOutput += entry.Value
//			}
//		}
//	}
//
//	return &oci.ContainerStats{
//		Container:   c.ID(),
//		CPU:         cpu,
//		CPUNano:     cpuNano,
//		SystemNano:  time.Now().UnixNano(),
//		MemUsage:    memUsage,
//		MemLimit:    memLimit,
//		MemPerc:     memPerc,
//		NetInput:    netInput,
//		NetOutput:   netOutput,
//		BlockInput:  blockInput,
//		BlockOutput: blockOutput,
//		PIDs:        pids,
//	}
//}
