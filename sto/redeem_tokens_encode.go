package sto

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *RedeemTokensFact) unpack(enc encoder.Encoder, ow string, bit []byte) error {
	e := util.StringErrorFunc("failed to unmarshal RedeemTokensFact")

	switch a, err := base.DecodeAddress(ow, enc); {
	case err != nil:
		return e(err, "")
	default:
		fact.sender = a
	}

	hit, err := enc.DecodeSlice(bit)
	if err != nil {
		return e(err, "")
	}

	items := make([]RedeemTokensItem, len(hit))
	for i := range hit {
		j, ok := hit[i].(RedeemTokensItem)
		if !ok {
			return e(util.ErrWrongType.Errorf("expected RedeemTokensItem, not %T", hit[i]), "")
		}

		items[i] = j
	}
	fact.items = items

	return nil
}