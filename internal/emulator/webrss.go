package emulator

import (
	"encoding/hex"
	"fmt"
	glb "myruscoint/internal/globals"
	"myruscoint/views"
	"slices"
	"strconv"
	"time"
)

func (wb *EmulatorWeb) RssLogSend(i int, msg string) {
	m := renderLogRow(wb.ctx, i, time.Now().Format(glb.LOG_DATE_FORMAT)+msg)
	e := NewRssEvent().WithEvent([]byte(glb.RSS_LOG_EVENT)).WithData(m)
	wb.rssChan <- *e
}

func (wb *EmulatorWeb) RssLogInfoSend(msg string, a ...any) {
	wb.RssLogSend(glb.LOG_LVL_INFO, fmt.Sprintf(msg, a...))
}
func (wb *EmulatorWeb) RssLogErrorSend(msg string, a ...any) {
	wb.RssLogSend(glb.LOG_LVL_ERROR, fmt.Sprintf(msg, a...))
}
func (wb *EmulatorWeb) RssLogOKSend(msg string, a ...any) {
	wb.RssLogSend(glb.LOG_LVL_OK, fmt.Sprintf(msg, a...))
}
func (wb *EmulatorWeb) RssLogEvilSend(msg string, a ...any) {
	wb.RssLogSend(glb.LOG_LVL_EVIL, fmt.Sprintf(msg, a...))
}

func (wb *EmulatorWeb) rssSendMinerStatusUpdate(name string, isMain bool) {
	msg := renderViewToBytes(wb.ctx, views.NodeMode(isMain))
	e := NewRssEvent().
		WithEvent([]byte(name + glb.RSS_EVENT_MINER_SET)).
		WithData(msg)
	wb.rssChan <- *e
}

func (wb *EmulatorWeb) RssSendMainMinerUpdates(name string) {
	wb.rssSendMinerStatusUpdate(name, true)
	idleNodes := wb.RcMngr.NodeNames()
	idleNodes = slices.DeleteFunc(idleNodes, func(n string) bool { return n == name })
	for _, n := range idleNodes {
		wb.rssSendMinerStatusUpdate(n, false)
	}
}

func (wb *EmulatorWeb) RssSendMinerSelect() {
	if n := wb.RcMngr.mainNode; n != nil {
		name := n.Name
		wb.RssLogOKSend("Miner set to " + name)
		wb.RssSendMainMinerUpdates(name)
	}
}

func (wb *EmulatorWeb) RssSendNodeWalletUpdate(id string) {
	node, ok := wb.RcMngr.Nodes[id]
	if !ok {
		return
	}
	e := NewRssEvent().
		WithEvent([]byte(node.Name + glb.RSS_EVENT_WALLET_COINS)).
		WithData([]byte(strconv.Itoa(node.Wallet.Balance())))
	wb.rssChan <- *e
}

func (wb *EmulatorWeb) RssSendNodeLastBlock(id string) {
	n, ok := wb.RcMngr.Nodes[id]
	if !ok {
		return
	}
	b := n.GetLastBlock()
	if b == nil {
		return
	}
	bsm := views.BlockInfoSmallItem{
		Height:   strconv.Itoa(b.Header.Height),
		Coinbase: strconv.Itoa(b.Body.Coinbase),
		Nonce:    strconv.Itoa(b.Header.Nonce),
		Hash:     hex.EncodeToString(b.Header.Hash),
		Root:     hex.EncodeToString(b.Header.Root),
	}
	msg := renderViewToBytes(wb.ctx, views.BlockInfoSmall(bsm))
	e := NewRssEvent().
		WithEvent([]byte(n.Name + glb.RSS_EVENT_LASTBLOCK)).
		WithData([]byte(msg))
	wb.rssChan <- *e
}

func (wb *EmulatorWeb) RssSendNodeCoinbase(id string) {
	n, ok := wb.RcMngr.Nodes[id]
	if !ok {
		return
	}
	msg := NewRssEvent().
		WithEvent([]byte(n.Name + glb.RSS_EVENT_NODE_COINBASE)).
		WithData([]byte(fmt.Sprintf("<span>%d</span>", n.CoinbaseUtxoAmount())))
	wb.rssChan <- *msg
}

// Send all RSS messages related to node with given id
func (wb *EmulatorWeb) RssNodeAllUpdates(id string) {
	wb.RssSendNodeCoinbase(id)
	wb.RssSendNodeLastBlock(id)
	wb.RssSendNodeWalletUpdate(id)
}

// Send RSS messages about all nodes
func (wb *EmulatorWeb) RssAllNodesUpdates() {
	for id := range wb.RcMngr.Nodes {
		wb.RssSendNodeCoinbase(id)
		wb.RssSendNodeLastBlock(id)
		wb.RssSendNodeWalletUpdate(id)
	}
}

func (wb *EmulatorWeb) RssTick() {
	msg := NewRssEvent().
		WithEvent([]byte(glb.RSS_EVENT_TICK)).
		WithData([]byte(strconv.Itoa(wb.RcMngr.Tick)))
	wb.rssChan <- *msg
}
