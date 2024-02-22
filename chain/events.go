package chain

import (
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type MachineEvent struct {
	Owner solana.PublicKey
	Uuid  [16]uint8
}

func (obj *MachineEvent) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Deserialize `Owner`:
	err = decoder.Decode(&obj.Owner)
	if err != nil {
		return err
	}
	// Deserialize `Uuid`:
	err = decoder.Decode(&obj.Uuid)
	if err != nil {
		return err
	}
	return nil
}

type OrderEvent struct {
	OrderId   [16]uint8
	Buyer     solana.PublicKey
	Seller    solana.PublicKey
	MachineId [16]uint8
}

func (obj *OrderEvent) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Deserialize `OrderId`:
	err = decoder.Decode(&obj.OrderId)
	if err != nil {
		return err
	}
	// Deserialize `Buyer`:
	err = decoder.Decode(&obj.Buyer)
	if err != nil {
		return err
	}
	// Deserialize `Seller`:
	err = decoder.Decode(&obj.Seller)
	if err != nil {
		return err
	}
	// Deserialize `MachineId`:
	err = decoder.Decode(&obj.MachineId)
	if err != nil {
		return err
	}
	return nil
}

type TaskEvent struct {
	Uuid      [16]uint8
	Period    uint32
	Owner     solana.PublicKey
	MachineId [16]uint8
}

func (obj *TaskEvent) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Deserialize `Uuid`:
	err = decoder.Decode(&obj.Uuid)
	if err != nil {
		return err
	}
	// Deserialize `Period`:
	err = decoder.Decode(&obj.Period)
	if err != nil {
		return err
	}
	// Deserialize `Owner`:
	err = decoder.Decode(&obj.Owner)
	if err != nil {
		return err
	}
	// Deserialize `MachineId`:
	err = decoder.Decode(&obj.MachineId)
	if err != nil {
		return err
	}
	return nil
}

type RewardEvent struct {
	Period    uint32
	Owner     solana.PublicKey
	MachineId [16]uint8
}

func (obj *RewardEvent) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// Deserialize `Period`:
	err = decoder.Decode(&obj.Period)
	if err != nil {
		return err
	}
	// Deserialize `Owner`:
	err = decoder.Decode(&obj.Owner)
	if err != nil {
		return err
	}
	// Deserialize `MachineId`:
	err = decoder.Decode(&obj.MachineId)
	if err != nil {
		return err
	}
	return nil
}
