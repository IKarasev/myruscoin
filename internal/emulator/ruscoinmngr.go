package emulator

import (
	"fmt"
	"math/rand"
	"myruscoint/internal/ruscoin"
)

type RuscoinMngr struct {
	Nodes     map[string]*ruscoin.Node
	Wallets   map[string]*ruscoin.Wallet
	mainNode  *ruscoin.Node
	EvilBlock *ruscoin.Block
	Tick      int
}

func NewRuscoinMngr() *RuscoinMngr {
	return &RuscoinMngr{
		Nodes:    make(map[string]*ruscoin.Node),
		Wallets:  make(map[string]*ruscoin.Wallet),
		mainNode: nil,
	}
}

func DefaultRuscoinMngr() *RuscoinMngr {
	rm := &RuscoinMngr{
		Nodes:    make(map[string]*ruscoin.Node),
		Wallets:  make(map[string]*ruscoin.Wallet),
		mainNode: nil,
	}
	for _, name := range []string{"Node1", "Node2", "Node3"} {
		rm.NewNode(name)
	}
	rm.NewWallet("User1")
	return rm
}

func (rm *RuscoinMngr) SelectMainNode() *ruscoin.Node {
	if rm.mainNode != nil {
		rm.mainNode.BlockCandidate = nil
	}
	l := len(rm.Nodes)
	ids := make([]string, 0, l)
	for k := range rm.Nodes {
		ids = append(ids, k)
	}
	rm.mainNode = rm.Nodes[ids[rand.Intn(l)]]
	rm.mainNode.NewBlockCandidate()
	return rm.mainNode
}

func (rm *RuscoinMngr) GetSetMainNode() *ruscoin.Node {
	if rm.mainNode != nil {
		return rm.mainNode
	}
	return rm.SelectMainNode()
}

func (rm *RuscoinMngr) MainNode() *ruscoin.Node {
	return rm.mainNode
}

func (rm *RuscoinMngr) NodeNames() []string {
	names := make([]string, len(rm.Nodes))
	i := 0
	for _, n := range rm.Nodes {
		names[i] = n.Name
		i++
	}
	return names
}

func (rm *RuscoinMngr) NewWallet(name string) (*ruscoin.Wallet, error) {
	w, err := ruscoin.NewWallet(name)
	if err != nil {
		return nil, err
	}
	rm.Wallets[w.Addr] = w
	return w, nil
}

func (rm *RuscoinMngr) AddWallet(w *ruscoin.Wallet) {
	rm.Wallets[w.Addr] = w
}

func (rm *RuscoinMngr) NewNode(name string) (*ruscoin.Node, error) {
	n, err := ruscoin.NewNode(name)
	if err != nil {
		return nil, err
	}
	rm.Nodes[n.Id] = n
	rm.AddWallet(n.Wallet)
	return n, err
}

func (rm *RuscoinMngr) Mine() (*ruscoin.Block, error) {
	b, err := rm.GetSetMainNode().Mine()
	if err != nil {
		return nil, err
	}
	rm.UpdateWalletsUtxo(b)
	return b, nil
}

func (rm *RuscoinMngr) EvryNode(f func(n *ruscoin.Node) error) error {
	for _, n := range rm.Nodes {
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (rm *RuscoinMngr) EvryNonMainNode(f func(n *ruscoin.Node) error) error {
	if rm.mainNode == nil {
		return fmt.Errorf("Emulator Server: main node not set")
	}
	for k, n := range rm.Nodes {
		if k == rm.mainNode.Id {
			continue
		}
		if err := f(n); err != nil {
			return err
		}
	}
	return nil
}

func (rm *RuscoinMngr) GetNodeBlock(nodeId string, height int) (*ruscoin.Block, error) {
	n, ok := rm.Nodes[nodeId]
	if !ok {
		return nil, fmt.Errorf("RuscoinMngr: Node [%s] not found", nodeId)
	}
	if height == -1 {
		if n.BlockCandidate == nil {
			return nil, fmt.Errorf("RuscoinMngr: Node [%s]: no block candidate", nodeId)
		}
		return n.BlockCandidate, nil
	} else if height < 0 || height >= len(n.BlockChain) {
		return nil, fmt.Errorf("RuscoinMngr: Node [%s]: invalid block height", nodeId)
	}
	return n.BlockChain[height], nil
}

func (rm *RuscoinMngr) UpdateWalletsUtxo(b *ruscoin.Block) {
	for _, t := range b.Body.Transactions {
		for id, u := range t.InputUtxo {
			if w, ok := rm.Wallets[u.Addr]; ok {
				w.RemoveUtxo(id)
			}
		}
		for id, u := range t.OutputUtxo {
			if w, ok := rm.Wallets[u.Addr]; ok {
				w.AddUtxo(id, u.Addr, u.Amount)
			}
		}
	}
}

func (rm *RuscoinMngr) consensusCheck(b *ruscoin.Block) (bool, []error) {
	nl := len(rm.Nodes)
	errs := make([]error, 0, nl)
	for _, n := range rm.Nodes {
		if err := n.VerifyBlock(b); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > nl/2 {
		return false, errs
	}
	return true, errs
}
