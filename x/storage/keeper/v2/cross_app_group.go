package v2

import (
	"encoding/hex"

	"cosmossdk.io/math"
	"github.com/bnb-chain/greenfield/x/storage/keeper"
	"github.com/bnb-chain/greenfield/x/storage/types"
	v2 "github.com/bnb-chain/greenfield/x/storage/types/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.CrossChainApplication = &GroupApp{}

type GroupApp struct {
	*keeper.GroupApp
	v2Keeper *Keeper
}

func NewGroupApp(v1GroupApp *keeper.GroupApp, storageKeeper *Keeper) *GroupApp {
	return &GroupApp{
		GroupApp: v1GroupApp,
		v2Keeper: storageKeeper,
	}
}

func (app *GroupApp) ExecuteSynPackage(ctx sdk.Context, appCtx *sdk.CrossChainAppContext, payload []byte) sdk.ExecuteResult {
	pack, err := v2.DeserializeCrossChainPackage(payload, types.GroupChannelId, sdk.SynCrossChainPackageType)
	if err != nil {
		app.GetStorageKeeper().Logger(ctx).Error("deserialize group cross chain package error", "payload", hex.EncodeToString(payload), "error", err.Error())
		return sdk.ExecuteResult{
			Payload: types.UpdateGroupMemberAckPackage{
				Status: types.StatusFail,
			}.MustSerialize(),
			Err: types.ErrInvalidCrossChainPackage,
		}
	}

	var operationType uint8
	var result sdk.ExecuteResult
	switch p := pack.(type) {
	case *types.MirrorGroupSynPackage:
		operationType = types.OperationMirrorGroup
		result = app.HandleMirrorGroupSynPackage(ctx, appCtx, p)
	case *types.CreateGroupSynPackage:
		operationType = types.OperationCreateGroup
		result = app.HandleCreateGroupSynPackage(ctx, appCtx, p)
	case *types.DeleteGroupSynPackage:
		operationType = types.OperationDeleteGroup
		result = app.HandleDeleteGroupSynPackage(ctx, appCtx, p)
	case *v2.UpdateGroupMemberSynPackage:
		operationType = types.OperationUpdateGroupMember
		result = app.HandleUpdateGroupMemberSynPackage(ctx, appCtx, p)
	default:
		return sdk.ExecuteResult{
			Err: types.ErrInvalidCrossChainPackage,
		}
	}

	if len(result.Payload) != 0 {
		wrapPayload := types.CrossChainPackage{
			OperationType: operationType,
			Package:       result.Payload,
		}
		result.Payload = wrapPayload.MustSerialize()
	}

	return result
}

func (app *GroupApp) HandleUpdateGroupPackageOperation(ctx sdk.Context, pkg *v2.UpdateGroupMemberSynPackage) (v2.UpdateGroupMemberOptions, error) {
	options := v2.UpdateGroupMemberOptions{
		SourceType: types.SOURCE_TYPE_BSC_CROSS_CHAIN,
	}
	switch pkg.OperationType {
	case types.OperationAddGroupMember:
		members := pkg.GetMembers()
		memberExpiration := pkg.GetMemberExpiration()
		memberWithExpirationTime := make([]*v2.MsgGroupMember, 0, len(members))
		for i := range members {
			memberWithExpirationTime = append(memberWithExpirationTime, &v2.MsgGroupMember{
				Member:         members[i],
				ExpirationTime: memberExpiration[i],
			})
		}

		options.MemberWithExpirationTime = memberWithExpirationTime
	case types.OperationDeleteGroupMember:
		options.MembersToDelete = pkg.GetMembers()
	}

	return options, nil
}

func (app *GroupApp) HandleUpdateGroupMemberSynPackage(ctx sdk.Context, header *sdk.CrossChainAppContext, updateGroupPackage *v2.UpdateGroupMemberSynPackage) sdk.ExecuteResult {
	err := updateGroupPackage.ValidateBasic()
	if err != nil {
		return sdk.ExecuteResult{
			Payload: types.UpdateGroupMemberAckPackage{
				Status:    types.StatusFail,
				Operator:  updateGroupPackage.Operator,
				ExtraData: updateGroupPackage.ExtraData,
			}.MustSerialize(),
			Err: err,
		}
	}

	groupInfo, found := app.GetStorageKeeper().GetGroupInfoById(ctx, math.NewUintFromBigInt(updateGroupPackage.GroupId))
	if !found {
		return sdk.ExecuteResult{
			Payload: types.UpdateGroupMemberAckPackage{
				Status:    types.StatusFail,
				Operator:  updateGroupPackage.Operator,
				ExtraData: updateGroupPackage.ExtraData,
			}.MustSerialize(),
			Err: types.ErrNoSuchGroup,
		}
	}

	options, err := app.HandleUpdateGroupPackageOperation(ctx, updateGroupPackage)
	if err != nil {
		return sdk.ExecuteResult{
			Payload: types.UpdateGroupMemberAckPackage{
				Status:    types.StatusFail,
				Operator:  updateGroupPackage.Operator,
				ExtraData: updateGroupPackage.ExtraData,
			}.MustSerialize(),
			Err: err,
		}
	}

	err = app.v2Keeper.UpdateGroupMember(
		ctx,
		updateGroupPackage.Operator,
		groupInfo,
		options,
	)
	if err != nil {
		return sdk.ExecuteResult{
			Payload: types.UpdateGroupMemberAckPackage{
				Status:    types.StatusFail,
				Operator:  updateGroupPackage.Operator,
				ExtraData: updateGroupPackage.ExtraData,
			}.MustSerialize(),
			Err: err,
		}
	}

	return sdk.ExecuteResult{
		Payload: types.UpdateGroupMemberAckPackage{
			Status:        types.StatusSuccess,
			Id:            groupInfo.Id.BigInt(),
			Operator:      updateGroupPackage.Operator,
			OperationType: updateGroupPackage.OperationType,
			Members:       updateGroupPackage.Members,
			ExtraData:     updateGroupPackage.ExtraData,
		}.MustSerialize(),
	}
}
