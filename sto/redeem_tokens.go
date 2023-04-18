package sto

import (
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	RedeemTokensFactHint = hint.MustNewHint("mitum-sto-redeem-tokens-operation-fact-v0.0.1")
	RedeemTokensHint     = hint.MustNewHint("mitum-sto-redeem-tokens-operation-v0.0.1")
)

type RedeemTokensItem interface {
	hint.Hinter
	util.IsValider
	Bytes() []byte
	Amount() currency.Amount
	Account() base.Address
	Partition() Partition
	Addresses() []base.Address
}

type RedeemTokensFact struct {
	base.BaseFact
	sender base.Address
	items  []RedeemTokensItem
}

func NewRedeemTokensFact(token []byte, receiver base.Address, items []RedeemTokensItem) RedeemTokensFact {
	bf := base.NewBaseFact(RedeemTokensFactHint, token)
	fact := RedeemTokensFact{
		BaseFact: bf,
		sender:   receiver,
		items:    items,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact RedeemTokensFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact RedeemTokensFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact RedeemTokensFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact RedeemTokensFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return util.ErrInvalid.Errorf("empty items")
	}

	if err := util.CheckIsValiders(nil, false, fact.sender); err != nil {
		return err
	}

	for i := range fact.items {
		if err := util.CheckIsValiders(nil, false, fact.items[i]); err != nil {
			return err
		}
	}

	return nil
}

func (fact RedeemTokensFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact RedeemTokensFact) Sender() base.Address {
	return fact.sender
}

func (fact RedeemTokensFact) Items() []RedeemTokensItem {
	return fact.items
}

func (fact RedeemTokensFact) Addresses() ([]base.Address, error) {
	as := []base.Address{}

	adrMap := make(map[string]struct{})
	for i := range fact.items {
		for j := range fact.items[i].Addresses() {
			if _, found := adrMap[fact.items[i].Addresses()[j].String()]; !found {
				adrMap[fact.items[i].Addresses()[j].String()] = struct{}{}
				as = append(as, fact.items[i].Addresses()[j])
			}
		}
	}
	as = append(as, fact.sender)

	return as, nil
}

type RedeemTokens struct {
	currency.BaseOperation
}

func NewRedeemTokens(fact RedeemTokensFact) (RedeemTokens, error) {
	return RedeemTokens{BaseOperation: currency.NewBaseOperation(IssueSecurityTokensHint, fact)}, nil
}

func (op *RedeemTokens) HashSign(priv base.Privatekey, networkID base.NetworkID) error {
	err := op.Sign(priv, networkID)
	if err != nil {
		return err
	}

	return nil
}
