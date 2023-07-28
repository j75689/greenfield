package v2

import (
	"math"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	permtypes "github.com/bnb-chain/greenfield/x/permission/types"
	"github.com/bnb-chain/greenfield/x/storage/keeper"
	"github.com/bnb-chain/greenfield/x/storage/types"
	v2 "github.com/bnb-chain/greenfield/x/storage/types/v2"
)

type Keeper struct {
	v1Keeper keeper.Keeper
}

func NewKeeper(v1Keeper keeper.Keeper) Keeper {
	return Keeper{v1Keeper: v1Keeper}
}

type CreateGroupOptions struct {
	Members                  []string
	SourceType               types.SourceType
	Extra                    string
	MemberWithExpirationTime []*v2.MsgGroupMember
}

// CreateGroup creates a group with the given name and members.
func (k Keeper) CreateGroup(
	ctx sdk.Context, owner sdk.AccAddress,
	groupName string, opts CreateGroupOptions) (sdkmath.Uint, error) {
	store := ctx.KVStore(k.v1Keeper.GetStoreKey())

	groupInfo := types.GroupInfo{
		Owner:      owner.String(),
		SourceType: opts.SourceType,
		Id:         k.v1Keeper.GenNextGroupId(ctx),
		GroupName:  groupName,
		Extra:      opts.Extra,
	}

	// Can not create a group with the same name.
	groupKey := types.GetGroupKey(owner, groupName)
	if store.Has(groupKey) {
		return sdkmath.ZeroUint(), types.ErrGroupAlreadyExists
	}

	gbz := k.v1Keeper.GetCdc().MustMarshal(&groupInfo)
	store.Set(groupKey, k.v1Keeper.GetGroupSeq().EncodeSequence(groupInfo.Id))
	store.Set(types.GetGroupByIDKey(groupInfo.Id), gbz)

	groupMemberInfoEvents := make([]*v2.MsgGroupMember, 0, len(opts.Members)+len(opts.MemberWithExpirationTime))

	// need to limit the size of Msg.Members to avoid taking too long to execute the msg
	for _, member := range opts.Members {
		memberAddress, err := sdk.AccAddressFromHexUnsafe(member)
		if err != nil {
			return sdkmath.ZeroUint(), err
		}
		_, err = k.v1Keeper.GetPermKeeper().AddGroupMember(ctx, groupInfo.Id, memberAddress)
		if err != nil {
			return sdkmath.Uint{}, err
		}

		groupMemberInfoEvents = append(groupMemberInfoEvents, &v2.MsgGroupMember{
			Member:         member,
			ExpirationTime: time.Unix(int64(math.MaxInt64), 0), // if the member has no expiration time, set it to MaxInt64
		})

	}

	for _, member := range opts.MemberWithExpirationTime {
		memberAddress, err := sdk.AccAddressFromHexUnsafe(member.Member)
		if err != nil {
			return sdkmath.ZeroUint(), err
		}

		err = k.v1Keeper.GetPermKeeper().AddGroupMemberWithExpiration(ctx, groupInfo.Id, memberAddress, member.ExpirationTime.UTC())
		if err != nil {
			return sdkmath.Uint{}, err
		}

		groupMemberInfoEvents = append(groupMemberInfoEvents, member)
	}

	if err := ctx.EventManager().EmitTypedEvents(&v2.EventCreateGroup{
		Owner:         groupInfo.Owner,
		GroupName:     groupInfo.GroupName,
		GroupId:       groupInfo.Id,
		SourceType:    groupInfo.SourceType,
		Members:       opts.Members,
		Extra:         opts.Extra,
		MembersDetail: groupMemberInfoEvents,
	}); err != nil {
		return sdkmath.ZeroUint(), err
	}
	return groupInfo.Id, nil
}

// UpdateGroupMember updates the members of a group.
func (k Keeper) UpdateGroupMember(ctx sdk.Context, operator sdk.AccAddress, groupInfo *types.GroupInfo, opts v2.UpdateGroupMemberOptions) error {
	if groupInfo.SourceType != opts.SourceType {
		return types.ErrSourceTypeMismatch
	}

	// check permission
	effect := k.v1Keeper.VerifyGroupPermission(ctx, groupInfo, operator, permtypes.ACTION_UPDATE_GROUP_MEMBER)
	if effect != permtypes.EFFECT_ALLOW {
		return types.ErrAccessDenied.Wrapf(
			"The operator(%s) has no UpdateGroupMember permission of the group(%s), operator(%s)",
			operator.String(), groupInfo.GroupName, groupInfo.Owner)
	}

	groupMemberInfoEvents := make([]*v2.MsgGroupMember, 0, len(opts.MembersToAdd)+len(opts.MemberWithExpirationTime))
	for _, member := range opts.MembersToAdd {
		memberAcc, err := sdk.AccAddressFromHexUnsafe(member)
		if err != nil {
			return err
		}
		_, err = k.v1Keeper.GetPermKeeper().AddGroupMember(ctx, groupInfo.Id, memberAcc)
		if err != nil {
			return err
		}

		groupMemberInfoEvents = append(groupMemberInfoEvents, &v2.MsgGroupMember{
			Member:         member,
			ExpirationTime: time.Unix(int64(math.MaxInt64), 0), // if the member has no expiration time, set it to MaxInt64
		})
	}

	for _, member := range opts.MembersToDelete {
		memberAcc, err := sdk.AccAddressFromHexUnsafe(member)
		if err != nil {
			return err
		}
		err = k.v1Keeper.GetPermKeeper().RemoveGroupMember(ctx, groupInfo.Id, memberAcc)
		if err != nil {
			return err
		}
	}

	for _, member := range opts.MemberWithExpirationTime {
		memberAcc, err := sdk.AccAddressFromHexUnsafe(member.Member)
		if err != nil {
			return err
		}

		err = k.v1Keeper.GetPermKeeper().AddGroupMemberWithExpiration(ctx, groupInfo.Id, memberAcc, member.ExpirationTime.UTC())
		if err != nil {
			return err
		}

		groupMemberInfoEvents = append(groupMemberInfoEvents, member)
	}

	if err := ctx.EventManager().EmitTypedEvents(&v2.EventUpdateGroupMember{
		Operator:           operator.String(),
		Owner:              groupInfo.Owner,
		GroupName:          groupInfo.GroupName,
		GroupId:            groupInfo.Id,
		MembersToAdd:       opts.MembersToAdd,
		MembersToDelete:    opts.MembersToDelete,
		AddedMembersDetail: groupMemberInfoEvents,
	}); err != nil {
		return err
	}
	return nil
}

func (k Keeper) LeaveGroup(
	ctx sdk.Context, member sdk.AccAddress, owner sdk.AccAddress,
	groupName string, opts types.LeaveGroupOptions) error {

	groupInfo, err := k.v1Keeper.LeaveGroup(ctx, member, owner, groupName, opts)
	if err != nil {
		return err
	}

	// remove group member extra
	return k.v1Keeper.GetPermKeeper().RemoveGroupMemberExtra(ctx, groupInfo.Id, member)
}
