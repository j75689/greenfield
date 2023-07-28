package v2

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/greenfield/x/storage/keeper"
	"github.com/bnb-chain/greenfield/x/storage/types"
	v2 "github.com/bnb-chain/greenfield/x/storage/types/v2"
)

type MsgServer struct {
	Keeper
	v1MsgServer keeper.MsgServer
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, v1MsgServer keeper.MsgServer) v2.MsgServer {
	return &MsgServer{Keeper: keeper, v1MsgServer: v1MsgServer}
}

var _ v2.MsgServer = MsgServer{}

// basic operation of bucket
func (server MsgServer) CreateBucket(ctx context.Context, msg *types.MsgCreateBucket) (*types.MsgCreateBucketResponse, error) {
	return server.v1MsgServer.CreateBucket(ctx, msg)
}

func (server MsgServer) DeleteBucket(ctx context.Context, msg *types.MsgDeleteBucket) (*types.MsgDeleteBucketResponse, error) {
	return server.v1MsgServer.DeleteBucket(ctx, msg)
}

func (server MsgServer) UpdateBucketInfo(ctx context.Context, msg *types.MsgUpdateBucketInfo) (*types.MsgUpdateBucketInfoResponse, error) {
	return server.v1MsgServer.UpdateBucketInfo(ctx, msg)
}

func (server MsgServer) MirrorBucket(ctx context.Context, msg *types.MsgMirrorBucket) (*types.MsgMirrorBucketResponse, error) {
	return server.v1MsgServer.MirrorBucket(ctx, msg)
}

func (server MsgServer) DiscontinueBucket(ctx context.Context, msg *types.MsgDiscontinueBucket) (*types.MsgDiscontinueBucketResponse, error) {
	return server.v1MsgServer.DiscontinueBucket(ctx, msg)
}

// basic operation of object
func (server MsgServer) CreateObject(ctx context.Context, msg *types.MsgCreateObject) (*types.MsgCreateObjectResponse, error) {
	return server.v1MsgServer.CreateObject(ctx, msg)
}

func (server MsgServer) SealObject(ctx context.Context, msg *types.MsgSealObject) (*types.MsgSealObjectResponse, error) {
	return server.v1MsgServer.SealObject(ctx, msg)
}

func (server MsgServer) RejectSealObject(ctx context.Context, msg *types.MsgRejectSealObject) (*types.MsgRejectSealObjectResponse, error) {
	return server.v1MsgServer.RejectSealObject(ctx, msg)
}

func (server MsgServer) CopyObject(ctx context.Context, msg *types.MsgCopyObject) (*types.MsgCopyObjectResponse, error) {
	return server.v1MsgServer.CopyObject(ctx, msg)
}

func (server MsgServer) DeleteObject(ctx context.Context, msg *types.MsgDeleteObject) (*types.MsgDeleteObjectResponse, error) {
	return server.v1MsgServer.DeleteObject(ctx, msg)
}

func (server MsgServer) CancelCreateObject(ctx context.Context, msg *types.MsgCancelCreateObject) (*types.MsgCancelCreateObjectResponse, error) {
	return server.v1MsgServer.CancelCreateObject(ctx, msg)
}

func (server MsgServer) MirrorObject(ctx context.Context, msg *types.MsgMirrorObject) (*types.MsgMirrorObjectResponse, error) {
	return server.v1MsgServer.MirrorObject(ctx, msg)
}

func (server MsgServer) DiscontinueObject(ctx context.Context, msg *types.MsgDiscontinueObject) (*types.MsgDiscontinueObjectResponse, error) {
	return server.v1MsgServer.DiscontinueObject(ctx, msg)
}

func (server MsgServer) UpdateObjectInfo(ctx context.Context, msg *types.MsgUpdateObjectInfo) (*types.MsgUpdateObjectInfoResponse, error) {
	return server.v1MsgServer.UpdateObjectInfo(ctx, msg)
}

// basic operation of group
func (server MsgServer) CreateGroup(goCtx context.Context, msg *v2.MsgCreateGroupV2) (*types.MsgCreateGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ownerAcc := sdk.MustAccAddressFromHex(msg.Creator)

	id, err := server.Keeper.CreateGroup(ctx, ownerAcc, msg.GroupName, CreateGroupOptions{Members: msg.Members, Extra: msg.Extra, MemberWithExpirationTime: msg.MembersWithExpiration})
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateGroupResponse{
		GroupId: id,
	}, nil
}

func (server MsgServer) DeleteGroup(ctx context.Context, msg *types.MsgDeleteGroup) (*types.MsgDeleteGroupResponse, error) {
	return server.v1MsgServer.DeleteGroup(ctx, msg)
}

func (server MsgServer) UpdateGroupMember(goCtx context.Context, msg *v2.MsgUpdateGroupMemberV2) (*types.MsgUpdateGroupMemberResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operator := sdk.MustAccAddressFromHex(msg.Operator)

	groupOwner := sdk.MustAccAddressFromHex(msg.GroupOwner)

	groupInfo, found := server.v1MsgServer.GetGroupInfo(ctx, groupOwner, msg.GroupName)
	if !found {
		return nil, types.ErrNoSuchGroup
	}
	err := server.Keeper.UpdateGroupMember(ctx, operator, groupInfo, v2.UpdateGroupMemberOptions{
		SourceType:               types.SOURCE_TYPE_ORIGIN,
		MembersToAdd:             msg.MembersToAdd,
		MembersToDelete:          msg.MembersToDelete,
		MemberWithExpirationTime: msg.MembersWithExpiration,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateGroupMemberResponse{}, nil
}

func (server MsgServer) UpdateGroupExtra(ctx context.Context, msg *types.MsgUpdateGroupExtra) (*types.MsgUpdateGroupExtraResponse, error) {
	return server.v1MsgServer.UpdateGroupExtra(ctx, msg)
}

func (server MsgServer) LeaveGroup(goCtx context.Context, msg *types.MsgLeaveGroup) (*types.MsgLeaveGroupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	memberAcc := sdk.MustAccAddressFromHex(msg.Member)
	ownerAcc := sdk.MustAccAddressFromHex(msg.GroupOwner)
	err := server.Keeper.LeaveGroup(ctx, memberAcc, ownerAcc, msg.GroupName, types.LeaveGroupOptions{SourceType: types.SOURCE_TYPE_ORIGIN})
	if err != nil {
		return nil, err
	}

	return &types.MsgLeaveGroupResponse{}, nil
}

func (server MsgServer) MirrorGroup(ctx context.Context, msg *types.MsgMirrorGroup) (*types.MsgMirrorGroupResponse, error) {
	return server.v1MsgServer.MirrorGroup(ctx, msg)
}

// basic operation of policy
func (server MsgServer) PutPolicy(ctx context.Context, msg *types.MsgPutPolicy) (*types.MsgPutPolicyResponse, error) {
	return server.v1MsgServer.PutPolicy(ctx, msg)
}

func (server MsgServer) DeletePolicy(ctx context.Context, msg *types.MsgDeletePolicy) (*types.MsgDeletePolicyResponse, error) {
	return server.v1MsgServer.DeletePolicy(ctx, msg)
}

// basic operation of params
func (server MsgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return server.v1MsgServer.UpdateParams(ctx, msg)
}

// basic operation of migrate
func (server MsgServer) MigrateBucket(ctx context.Context, msg *types.MsgMigrateBucket) (*types.MsgMigrateBucketResponse, error) {
	return server.v1MsgServer.MigrateBucket(ctx, msg)
}

func (server MsgServer) CompleteMigrateBucket(ctx context.Context, msg *types.MsgCompleteMigrateBucket) (*types.MsgCompleteMigrateBucketResponse, error) {
	return server.v1MsgServer.CompleteMigrateBucket(ctx, msg)
}

func (server MsgServer) CancelMigrateBucket(ctx context.Context, msg *types.MsgCancelMigrateBucket) (*types.MsgCancelMigrateBucketResponse, error) {
	return server.v1MsgServer.CancelMigrateBucket(ctx, msg)
}
