package sto

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *STODesign) unpack(enc encoder.Encoder, ht hint.Hint, sto string, gra uint64, bpo []byte) error {
	e := util.StringErrorFunc("failed to decode bson of STODesign")

	de.BaseHinter = hint.NewBaseHinter(ht)
	de.stoID = extensioncurrency.ContractID(sto)
	de.granularity = gra

	if hinter, err := enc.Decode(bpo); err != nil {
		return e(err, "")
	} else if po, ok := hinter.(STOPolicy); !ok {
		return e(util.ErrWrongType.Errorf("expected STOPolicy, not %T", hinter), "")
	} else {
		de.policy = po
	}

	return nil
}
