package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-sto/sto"
	"github.com/ProtoconNet/mitum2/base"
)

type IssueSecurityTokensCommand struct {
	baseCommand
	OperationFlags
	Sender    AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  AddressFlag    `arg:"" name:"contract" help:"contract account address" required:"true"`
	STO       ContractIDFlag `arg:"" name:"sto-id" help:"sto id" required:"true"`
	Receiver  AddressFlag    `arg:"" name:"receiver" help:"token receiver" required:"true"`
	Amount    BigFlag        `arg:"" name:"amount" help:"token amount" required:"true"`
	Partition PartitionFlag  `arg:"" name:"partition" help:"partition" required:"true"`
	Currency  CurrencyIDFlag `arg:"" name:"currency-id" help:"currency id" required:"true"`
	sender    base.Address
	contract  base.Address
	receiver  base.Address
}

func NewIssueSecurityTokensCommand() IssueSecurityTokensCommand {
	cmd := NewbaseCommand()
	return IssueSecurityTokensCommand{
		baseCommand: *cmd,
	}
}

func (cmd *IssueSecurityTokensCommand) Run(pctx context.Context) error {
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *IssueSecurityTokensCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sender, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sender

	contract, err := cmd.Contract.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid contract account format, %q", cmd.Contract.String())
	}
	cmd.contract = contract

	receiver, err := cmd.Receiver.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid receiver format, %q", cmd.Receiver.String())
	}
	cmd.receiver = receiver

	if !cmd.Amount.OverZero() {
		return errors.Wrap(nil, "amount must be over zero")
	}

	return nil
}

func (cmd *IssueSecurityTokensCommand) createOperation() (base.Operation, error) { // nolint:dupl
	var items []sto.IssueSecurityTokensItem

	item := sto.NewIssueSecurityTokensItem(
		cmd.contract,
		cmd.STO.ID,
		cmd.receiver,
		cmd.Amount.Big,
		cmd.Partition.Partition,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := sto.NewIssueSecurityTokensFact([]byte(cmd.Token), cmd.sender, items)

	op, err := sto.NewIssueSecurityTokens(fact)
	if err != nil {
		return nil, errors.Wrap(err, "failed to issue security tokens operation")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to issue security tokens operation")
	}

	return op, nil
}
