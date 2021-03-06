package db

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/vilterp/instance-manager-toy/proto"
)

type TaskGraphID string

// what task graphs are running
type TaskGraphsDB interface {
	Insert(g *proto.TaskGraphSpec) *proto.TaskGraph
	List() []*proto.TaskGraph
	Get(id TaskGraphID) (*proto.TaskGraph, bool)
	GetState(id TaskGraphID) TasksDB
	MarkDone(id TaskGraphID)
}

type MockTaskGraphsDB struct {
	graphs      map[TaskGraphID]*proto.TaskGraph
	graphStates map[TaskGraphID]TasksDB
}

func (g *MockTaskGraphsDB) MarkDone(id TaskGraphID) {
	graph := g.graphs[id]
	graph.FinishedAt = ptypes.TimestampNow()
	graph.State = proto.TaskGraphState_TaskGraphSucceeded
}

var _ TaskGraphsDB = &MockTaskGraphsDB{}

func NewMockTaskGraphsDB() *MockTaskGraphsDB {
	return &MockTaskGraphsDB{
		graphs:      map[TaskGraphID]*proto.TaskGraph{},
		graphStates: map[TaskGraphID]TasksDB{},
	}
}

func (g *MockTaskGraphsDB) Get(id TaskGraphID) (*proto.TaskGraph, bool) {
	graph, ok := g.graphs[id]
	if !ok {
		return nil, ok
	}
	tasks := g.graphStates[id].List()
	graph.Tasks = tasks
	return graph, ok
}

func (g *MockTaskGraphsDB) List() []*proto.TaskGraph {
	var out []*proto.TaskGraph
	for _, tg := range g.graphs {
		out = append(out, tg)
	}
	return out
}

func (g *MockTaskGraphsDB) GetState(id TaskGraphID) TasksDB {
	return g.graphStates[id]
}

func (g *MockTaskGraphsDB) Insert(spec *proto.TaskGraphSpec) *proto.TaskGraph {
	id := TaskGraphID(uuid.New().String())
	graph := &proto.TaskGraph{
		Id:        string(id),
		Spec:      spec,
		State:     proto.TaskGraphState_TaskGraphWaiting,
		CreatedAt: ptypes.TimestampNow(),
	}
	g.graphs[id] = graph
	g.graphStates[id] = NewMockTasksDB(spec)
	return graph
}
