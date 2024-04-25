package chain

import (
	"encoding/base64"
	"errors"
	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type DistriEvent interface {
	*MachineEvent | *OrderEvent | *TaskEvent | *RewardEvent

	UnmarshalWithDecoder(decoder *bin.Decoder) error
}

func decodeDistriEvent[T DistriEvent](data string) (T, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(bytes) < 8 {
		return nil, errors.New("data too short")
	}
	var event T
	if err := event.UnmarshalWithDecoder(bin.NewBorshDecoder(bytes[8:])); err != nil {
		return nil, err
	}
	return event, nil
}

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
