package cluster

import (
	"encoding/json"
	"strconv"
	"time"
)

// Node represents a farm with endpoint of One. Each node has an Address
// (in the form <scheme>://<host>:<port>/RPC2) and map with arbritary
// metadata.
type Node struct {
	Address        string
	Region         string `json:"_id"`
	Healing        HealingData
	Metadata       map[string]string
	Clusters       map[string]map[string]string
	CreationStatus string
}

type HealingData struct {
	LockedUntil time.Time
	IsFailure   bool
}

type NodeList []Node

const (
	NodeStatusWaiting             = "waiting"
	NodeStatusReady               = "ready"
	NodeStatusRetry               = "ready for retry"
	NodeStatusTemporarilyDisabled = "temporarily disabled"
	NodeStatusHealing             = "healing"

	NodeCreationStatusCreated  = "created"
	NodeCreationStatusError    = "error"
	NodeCreationStatusPending  = "pending"
	NodeCreationStatusDisabled = "disabled"
)

func (a NodeList) Len() int           { return len(a) }
func (a NodeList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NodeList) Less(i, j int) bool { return a[i].Address < a[j].Address }

func (n Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Region":   n.Region,
		"Address":  n.Address,
		"Metadata": n.Metadata,
		"Status":   n.Status(),
	})
}

func (n *Node) HasSuccess() bool {
	_, hasSuccess := n.Metadata["LastSuccess"]
	return hasSuccess
}

func (n *Node) Status() string {
	if n.CreationStatus != "" && n.CreationStatus != NodeCreationStatusCreated {
		return n.CreationStatus
	}
	if n.isHealing() {
		return NodeStatusHealing
	}
	if n.Metadata == nil {
		return NodeStatusWaiting
	}
	if n.isEnabled() {
		_, hasFailures := n.Metadata["Failures"]
		if hasFailures {
			return NodeStatusRetry
		}
		if !n.HasSuccess() {
			return NodeStatusWaiting
		}
		return NodeStatusReady
	}
	return NodeStatusTemporarilyDisabled
}

func (n *Node) isEnabled() bool {
	if n.CreationStatus != "" && n.CreationStatus != NodeCreationStatusCreated {
		return false
	}
	if n.isHealing() {
		return false
	}
	if n.Metadata == nil {
		return true
	}
	disabledStr, _ := n.Metadata["DisabledUntil"]
	t, _ := time.Parse(time.RFC3339, disabledStr)
	return time.Now().After(t)
}

func (n *Node) isHealing() bool {
	return (!n.Healing.LockedUntil.IsZero()) && n.Healing.IsFailure
}

func (nodes NodeList) filterDisabled() NodeList {
	filtered := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		if node.isEnabled() {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

func (n *Node) updateError(lastErr error, incrementFailures bool) {
	if n.Metadata == nil {
		n.Metadata = make(map[string]string)
	}
	if incrementFailures {
		n.Metadata["Failures"] = strconv.Itoa(n.FailureCount() + 1)
	}
	n.Metadata["LastError"] = lastErr.Error()
}

func (n *Node) updateDisabled(disabledUntil time.Time) {
	if n.Metadata == nil {
		n.Metadata = make(map[string]string)
	}
	n.Metadata["DisabledUntil"] = disabledUntil.Format(time.RFC3339)
}

func (n *Node) updateSuccess() {
	n.ResetFailures()
	n.Metadata["LastSuccess"] = time.Now().Format(time.RFC3339)
}

func (n *Node) FailureCount() int {
	if n.Metadata == nil {
		return 0
	}
	metaFail, _ := n.Metadata["Failures"]
	failures, _ := strconv.Atoi(metaFail)
	return failures
}

func (n *Node) ResetFailures() {
	if n.Metadata == nil {
		n.Metadata = make(map[string]string)
	}
	delete(n.Metadata, "Failures")
	delete(n.Metadata, "DisabledUntil")
	delete(n.Metadata, "LastError")
}

func (n *Node) CleanMetadata() map[string]string {
	paramsCopy := make(map[string]string)
	for k, v := range n.Metadata {
		paramsCopy[k] = v
	}
	delete(paramsCopy, "Failures")
	delete(paramsCopy, "DisabledUntil")
	delete(paramsCopy, "LastError")
	delete(paramsCopy, "LastSuccess")
	return paramsCopy
}
