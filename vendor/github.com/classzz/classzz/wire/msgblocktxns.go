package wire

import (
	"errors"
	"github.com/bourbaki-czz/classzz/chaincfg/chainhash"
	"io"
)

// MsgBlockTxns implements the Message interface and represents a Bitcoin blocktxn
// message.  It is sent in response to the getblocktxn message.
type MsgBlockTxns struct {
	BlockHash chainhash.Hash
	Txs       []*MsgTx
}

// CzzDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgBlockTxns) CzzDecode(r io.Reader, pver uint32, enc MessageEncoding) error {

	if err := readElement(r, &msg.BlockHash); err != nil {
		return err
	}

	txCount, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	for i := uint64(0); i < txCount; i++ {
		tx := MsgTx{}
		err = tx.CzzDecode(r, pver, enc)
		if err != nil {
			return err
		}
		msg.Txs = append(msg.Txs, &tx)
	}
	return nil
}

// CzzEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgBlockTxns) CzzEncode(w io.Writer, pver uint32, enc MessageEncoding) error {

	if err := writeElement(w, &msg.BlockHash); err != nil {
		return err
	}

	if err := WriteVarInt(w, pver, uint64(len(msg.Txs))); err != nil {
		return err
	}

	for _, tx := range msg.Txs {
		if err := tx.CzzEncode(w, pver, enc); err != nil {
			return err
		}
	}
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgBlockTxns) Command() string {
	return CmdBlockTxns
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgBlockTxns) MaxPayloadLength(pver uint32) uint32 {
	// In practice this will always be less than the payload but the number
	// of txs in a block can vary so we really don't know the real max.
	return maxMessagePayload()
}

// AbsoluteIndexes takes in the requested differential indexes from a MsgGetBlockTxns
// message and returns a map of the absolution position of a tx in the block to the tx.
func (msg *MsgBlockTxns) AbsoluteIndexes(requestedDiffIndexes []uint32) (map[uint32]*MsgTx, error) {
	if len(requestedDiffIndexes) != len(msg.Txs) {
		return nil, errors.New("requestedDiffIndexes length does not match length of txs in blocktxn message")
	}
	m := make(map[uint32]*MsgTx)
	lastIndex := uint32(0)
	for i, tx := range msg.Txs {
		index := requestedDiffIndexes[i]
		m[index+lastIndex] = tx
		lastIndex += index + 1
	}
	return m, nil
}

// NewMsgBlockTxns returns a new bitcoin blocktxn message that conforms to the
// Message interface using the passed parameters and defaults for the remaining
// fields.
func NewMsgBlockTxns(blockHash chainhash.Hash, txs []*MsgTx) *MsgBlockTxns {
	return &MsgBlockTxns{BlockHash: blockHash, Txs: txs}
}
