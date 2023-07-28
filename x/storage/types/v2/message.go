package v2

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	gnfderrors "github.com/bnb-chain/greenfield/types/errors"
	"github.com/bnb-chain/greenfield/types/s3util"
	types "github.com/bnb-chain/greenfield/x/storage/types"
)

const (
	// For group
	TypeMsgCreateGroup       = "create_group"
	TypeMsgUpdateGroupMember = "update_group_member"
)

var (
	// For group
	_ sdk.Msg = &MsgCreateGroup{}
	_ sdk.Msg = &MsgUpdateGroupMember{}
)

func NewMsgCreateGroup(creator sdk.AccAddress, groupName string, membersAcc []sdk.AccAddress, extra string, membersWithExpiration []*MsgGroupMember) *MsgCreateGroup {
	var members []string
	for _, member := range membersAcc {
		members = append(members, member.String())
	}
	return &MsgCreateGroup{
		Creator:               creator.String(),
		GroupName:             groupName,
		Members:               members,
		Extra:                 extra,
		MembersWithExpiration: membersWithExpiration,
	}
}

// Route implements the sdk.Msg interface.
func (msg *MsgCreateGroup) Route() string {
	return types.RouterKey
}

// Type implements the sdk.Msg interface.
func (msg *MsgCreateGroup) Type() string {
	return TypeMsgCreateGroup
}

// GetSigners implements the sdk.Msg interface.
func (msg *MsgCreateGroup) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromHexUnsafe(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (msg *MsgCreateGroup) GetSignBytes() []byte {
	bz := types.ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg *MsgCreateGroup) ValidateBasic() error {
	_, err := sdk.AccAddressFromHexUnsafe(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	err = s3util.CheckValidGroupName(msg.GroupName)
	if err != nil {
		return gnfderrors.ErrInvalidGroupName.Wrapf("invalid groupName (%s)", err)
	}
	if len(msg.Members)+len(msg.MembersWithExpiration) > types.MaxGroupMemberLimitOnce {
		return gnfderrors.ErrInvalidParameter.Wrapf("Once update group member limit exceeded")
	}
	if len(msg.Extra) > types.MaxGroupExtraInfoLimit {
		return errors.Wrapf(gnfderrors.ErrInvalidParameter, "extra is too long with length %d, limit to %d", len(msg.Extra), types.MaxGroupExtraInfoLimit)
	}

	return nil
}

func NewMsgUpdateGroupMember(
	operator sdk.AccAddress, groupOwner sdk.AccAddress, groupName string, membersToAdd []sdk.AccAddress,
	membersToDelete []sdk.AccAddress, membersWithExpiration []*MsgGroupMember) *MsgUpdateGroupMember {
	var membersAddrToAdd, membersAddrToDelete []string
	for _, member := range membersToAdd {
		membersAddrToAdd = append(membersAddrToAdd, member.String())
	}
	for _, member := range membersToDelete {
		membersAddrToDelete = append(membersAddrToDelete, member.String())
	}
	return &MsgUpdateGroupMember{
		Operator:              operator.String(),
		GroupOwner:            groupOwner.String(),
		GroupName:             groupName,
		MembersToAdd:          membersAddrToAdd,
		MembersToDelete:       membersAddrToDelete,
		MembersWithExpiration: membersWithExpiration,
	}
}

// Route implements the sdk.Msg interface.
func (msg *MsgUpdateGroupMember) Route() string {
	return types.RouterKey
}

// Type implements the sdk.Msg interface.
func (msg *MsgUpdateGroupMember) Type() string {
	return TypeMsgUpdateGroupMember
}

// GetSigners implements the sdk.Msg interface.
func (msg *MsgUpdateGroupMember) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromHexUnsafe(msg.Operator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the message bytes to sign over.
func (msg *MsgUpdateGroupMember) GetSignBytes() []byte {
	bz := types.ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg *MsgUpdateGroupMember) ValidateBasic() error {
	_, err := sdk.AccAddressFromHexUnsafe(msg.Operator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid operator address (%s)", err)
	}

	_, err = sdk.AccAddressFromHexUnsafe(msg.GroupOwner)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid group owner address (%s)", err)
	}

	err = s3util.CheckValidGroupName(msg.GroupName)
	if err != nil {
		return err
	}

	if len(msg.MembersToAdd)+len(msg.MembersToDelete)+len(msg.MembersWithExpiration) > types.MaxGroupMemberLimitOnce {
		return gnfderrors.ErrInvalidParameter.Wrapf("Once update group member limit exceeded")
	}
	for _, member := range msg.MembersToAdd {
		_, err = sdk.AccAddressFromHexUnsafe(member)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid member address (%s)", err)
		}
	}
	for _, member := range msg.MembersToDelete {
		_, err = sdk.AccAddressFromHexUnsafe(member)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid member address (%s)", err)
		}
	}
	for _, member := range msg.MembersWithExpiration {
		_, err = sdk.AccAddressFromHexUnsafe(member.Member)
		if err != nil {
			return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid member address (%s)", err)
		}
	}
	return nil
}
