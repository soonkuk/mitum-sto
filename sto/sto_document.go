package sto

import (
	"net/url"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type URI string

func (uri URI) Bytes() []byte {
	return []byte(uri)
}

func (uri URI) String() string {
	return string(uri)
}

func (uri URI) IsValid([]byte) error {
	if _, err := url.Parse(string(uri)); err != nil {
		return err
	}
	return nil
}

var (
	DocumentHint = hint.MustNewHint("mitum-sto-document-v0.0.1")
)

type Document struct {
	hint.BaseHinter
	stoID extensioncurrency.ContractID
	title string
	hash  string
	uri   URI
}

func NewDocument(stoID extensioncurrency.ContractID, title, hash string, uri URI) Document {
	return Document{
		BaseHinter: hint.NewBaseHinter(STODesignHint),
		stoID:      stoID,
		title:      title,
		hash:       hash,
		uri:        uri,
	}
}

func (s Document) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		s.BaseHinter,
		s.stoID,
		s.uri,
	); err != nil {
		return util.ErrInvalid.Errorf("invalid STODesign: %w", err)
	}

	return nil
}

func (s Document) Bytes() []byte {
	return util.ConcatBytesSlice(
		s.stoID.Bytes(),
		[]byte(s.title),
		[]byte(s.hash),
		s.uri.Bytes(),
	)
}

func (s Document) STO() extensioncurrency.ContractID {
	return s.stoID
}
