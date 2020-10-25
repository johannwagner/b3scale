package cluster

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	// "gitlab.com/infra.run/public/b3scale/pkg/bbb"
	"gitlab.com/infra.run/public/b3scale/pkg/store"
)

// The Controller interfaces with the state of the cluster
// providing methods for retrieving cluster backends and
// frontends.
//
// The controller subscribes to commands.
type Controller struct {
	cmds *store.CommandQueue
	conn *pgxpool.Pool
}

// NewController will initialize the cluster controller
// with a database connection. A BBB client will be created
// which will be used by the backend instances.
func NewController(conn *pgxpool.Pool) *Controller {
	return &Controller{
		cmds: store.NewCommandQueue(conn),
		conn: conn,
	}
}

// Start the controller
func (c *Controller) Start() {
	log.Println("Starting cluster controller")

	// Controller Main Loop
	for {
		// Process commands from queue
		if err := c.cmds.Receive(c.handleCommand); err != nil {
			// Log error and wait a bit
			log.Println(err)
			time.Sleep(1 * time.Second)
		}
	}
}

// Command callback handler: Decode the operation and
// run the command specific handler
func (c *Controller) handleCommand(cmd *store.Command) (interface{}, error) {
	// Invoke command handler
	switch cmd.Action {
	case CmdAddBackend:
		return c.handleAddBackend(cmd)
	case CmdLoadBackendState:
		return c.handleLoadBackendState(cmd)
	}

	return nil, ErrUnknownCommand
}

// Command: AddBackend
// Creates a new backend state and dispatches the initial
// load state.
func (c *Controller) handleAddBackend(cmd *store.Command) (interface{}, error) {
	params, ok := cmd.Params.(map[string]string)
	if !ok {
	}

}

// Command: LoadBackendState
func (c *Controller) handleLoadBackendState(
	cmd *store.Command,
) (interface{}, error) {
	// Get backend from command
	backendID, ok := cmd.Params.(string)
	if !ok {
		return false, fmt.Errorf("invalid backend id: %v", cmd.Params)
	}
	backend, err := c.GetBackend(store.NewQuery().Eq("id", backendID))
	if err != nil {
		return false, err
	}
	if backend == nil {
		return false, fmt.Errorf("backend not found: %s", backendID)
	}
	err = backend.loadBackendState()
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetBackends retrives backends with a store query
func (c *Controller) GetBackends(q *store.Query) ([]*Backend, error) {
	states, err := store.GetBackendStates(c.conn, q)
	if err != nil {
		return nil, err
	}
	// Make cluster backend from each state
	backends := make([]*Backend, 0, len(states))
	for _, s := range states {
		backends = append(backends, NewBackend(s))
	}

	return backends, nil
}

// GetBackend retrievs a single backend by query criteria
func (c *Controller) GetBackend(q *store.Query) (*Backend, error) {
	backends, err := c.GetBackends(q)
	if err != nil {
		return nil, err
	}
	if len(backends) == 0 {
		return nil, nil
	}
	return backends[0], nil
}
