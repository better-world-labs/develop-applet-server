package entity

import (
	"encoding/json"
	"errors"
)

const (
	Fen  = 1
	Jiao = Fen * 10
	Yuan = Jiao * 10
)

type CNY struct {
	fen int64
}

func FromYuan(rmb float32) CNY {
	return CNY{
		int64(rmb * Yuan),
	}
}

func FromJiao(rmb float32) CNY {
	return CNY{
		int64(rmb * Jiao),
	}
}

func FromFen(rmb int64) CNY {
	return CNY{
		rmb,
	}
}

func (r CNY) Fen() int64 {
	return r.fen
}

func (r CNY) Jiao() float32 {
	return float32(r.fen) / Jiao
}

func (r CNY) Yuan() float32 {
	return float32(r.fen) / Yuan
}

func (r CNY) Add(another CNY) CNY {
	return FromFen(r.fen + another.fen)
}

func (r CNY) Sub(another CNY) (CNY, error) {
	sub := r.fen - another.fen
	if sub < 0 {
		return CNY{}, errors.New("CNY must > 0")
	}

	return FromFen(sub), nil
}

func (r CNY) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Yuan())
}

func (r *CNY) UnmarshalJSON(b []byte) error {
	var rmb float32
	err := json.Unmarshal(b, &rmb)
	if err != nil {
		return err
	}

	*r = FromYuan(rmb)
	return err
}

func (r *CNY) FromDB(b []byte) error {
	var rmb int64
	err := json.Unmarshal(b, &rmb)
	if err != nil {
		return err
	}

	*r = FromFen(rmb)
	return err
}

func (r CNY) ToDB() ([]byte, error) {
	return json.Marshal(r.Fen())
}
