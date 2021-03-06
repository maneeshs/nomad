package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/nomad/helper/raftutil"
	"github.com/posener/complete"
)

type OperatorRaftInfoCommand struct {
	Meta
}

func (c *OperatorRaftInfoCommand) Help() string {
	helpText := `
Usage: nomad operator raft _info <path to nomad data dir>

  Displays info about the raft logs in the data directory.
  
  This is a low-level debugging tool and not subject to Nomad's usual backward
  compatibility guarantees.
`
	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftInfoCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *OperatorRaftInfoCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OperatorRaftInfoCommand) Synopsis() string {
	return "Display info of the raft log"
}

func (c *OperatorRaftInfoCommand) Name() string { return "operator raft _info" }

func (c *OperatorRaftInfoCommand) Run(args []string) int {
	if len(args) != 1 {
		c.Ui.Error("This command takes one argument: <path>")
		c.Ui.Error(commandErrorText(c))

		return 1
	}

	p := args[0]
	if fi, err := os.Stat(p); err == nil && fi.IsDir() {
		p = filepath.Join(args[0], "server", "raft", "raft.db")
	}

	store, firstIdx, lastIdx, err := raftutil.RaftStateInfo(p)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to open raft logs: %v", err))
		return 1
	}
	defer store.Close()

	c.Ui.Output(fmt.Sprintf("path:        %v", p))
	c.Ui.Output(fmt.Sprintf("length:      %v", lastIdx-firstIdx+1))
	c.Ui.Output(fmt.Sprintf("first index: %v", firstIdx))
	c.Ui.Output(fmt.Sprintf("last index:  %v", lastIdx))

	return 0
}
