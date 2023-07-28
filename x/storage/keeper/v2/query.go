package v2

import (
	"context"
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bnb-chain/greenfield/x/storage/types"
	v2 "github.com/bnb-chain/greenfield/x/storage/types/v2"
)

var _ v2.QueryServer = Keeper{}

// Parameters queries the parameters of the module.
func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return k.v1Keeper.Params(ctx, req)
}

// Parameters queries the parameters of the module.
func (k Keeper) QueryParamsByTimestamp(ctx context.Context, req *types.QueryParamsByTimestampRequest) (*types.QueryParamsByTimestampResponse, error) {
	return k.v1Keeper.QueryParamsByTimestamp(ctx, req)
}

// Queries a bucket with specify name.
func (k Keeper) HeadBucket(ctx context.Context, req *types.QueryHeadBucketRequest) (*types.QueryHeadBucketResponse, error) {
	return k.v1Keeper.HeadBucket(ctx, req)
}

// Queries a bucket by id
func (k Keeper) HeadBucketById(ctx context.Context, req *types.QueryHeadBucketByIdRequest) (*types.QueryHeadBucketResponse, error) {
	return k.v1Keeper.HeadBucketById(ctx, req)
}

// Queries a bucket with EIP712 standard metadata info
func (k Keeper) HeadBucketNFT(ctx context.Context, req *types.QueryNFTRequest) (*types.QueryBucketNFTResponse, error) {
	return k.v1Keeper.HeadBucketNFT(ctx, req)
}

// Queries a object with specify name.
func (k Keeper) HeadObject(ctx context.Context, req *types.QueryHeadObjectRequest) (*types.QueryHeadObjectResponse, error) {
	return k.v1Keeper.HeadObject(ctx, req)
}

// Queries an object by id
func (k Keeper) HeadObjectById(ctx context.Context, req *types.QueryHeadObjectByIdRequest) (*types.QueryHeadObjectResponse, error) {
	return k.v1Keeper.HeadObjectById(ctx, req)
}

// Queries a object with EIP712 standard metadata info
func (k Keeper) HeadObjectNFT(ctx context.Context, req *types.QueryNFTRequest) (*types.QueryObjectNFTResponse, error) {
	return k.v1Keeper.HeadObjectNFT(ctx, req)
}

// Queries a list of bucket items.
func (k Keeper) ListBuckets(ctx context.Context, req *types.QueryListBucketsRequest) (*types.QueryListBucketsResponse, error) {
	return k.v1Keeper.ListBuckets(ctx, req)
}

// Queries a list of object items under the bucket.
func (k Keeper) ListObjects(ctx context.Context, req *types.QueryListObjectsRequest) (*types.QueryListObjectsResponse, error) {
	return k.v1Keeper.ListObjects(ctx, req)
}

// Queries a list of object items under the bucket.
func (k Keeper) ListObjectsByBucketId(ctx context.Context, req *types.QueryListObjectsByBucketIdRequest) (*types.QueryListObjectsResponse, error) {
	return k.v1Keeper.ListObjectsByBucketId(ctx, req)
}

// Queries a group with EIP712 standard metadata info
func (k Keeper) HeadGroupNFT(ctx context.Context, req *types.QueryNFTRequest) (*types.QueryGroupNFTResponse, error) {
	return k.v1Keeper.HeadGroupNFT(ctx, req)
}

// Queries a policy which grants permission to account
func (k Keeper) QueryPolicyForAccount(ctx context.Context, req *types.QueryPolicyForAccountRequest) (*types.QueryPolicyForAccountResponse, error) {
	return k.v1Keeper.QueryPolicyForAccount(ctx, req)
}

// Queries a list of VerifyPermission items.
func (k Keeper) VerifyPermission(ctx context.Context, req *types.QueryVerifyPermissionRequest) (*types.QueryVerifyPermissionResponse, error) {
	return k.v1Keeper.VerifyPermission(ctx, req)
}

// Queries a group with specify owner and name .
func (k Keeper) HeadGroup(ctx context.Context, req *types.QueryHeadGroupRequest) (*types.QueryHeadGroupResponse, error) {
	return k.v1Keeper.HeadGroup(ctx, req)
}

// Queries a list of ListGroup items.
func (k Keeper) ListGroup(ctx context.Context, req *types.QueryListGroupRequest) (*types.QueryListGroupResponse, error) {
	return k.v1Keeper.ListGroup(ctx, req)
}

// Queries a list of HeadGroupMember items.
func (k Keeper) HeadGroupMember(goCtx context.Context, req *types.QueryHeadGroupMemberRequest) (*v2.QueryHeadGroupMemberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	member, err := sdk.AccAddressFromHexUnsafe(req.Member)
	if err != nil {
		return nil, err
	}
	owner, err := sdk.AccAddressFromHexUnsafe(req.GroupOwner)
	if err != nil {
		return nil, err
	}
	groupInfo, found := k.v1Keeper.GetGroupInfo(ctx, owner, req.GroupName)
	if !found {
		return nil, types.ErrNoSuchGroup
	}
	groupMember, found := k.v1Keeper.GetPermKeeper().GetGroupMember(ctx, groupInfo.Id, member)
	if !found {
		return nil, types.ErrNoSuchGroupMember
	}
	memberExpirationTime := time.Unix(math.MaxInt64, 0)
	groupMemberExtra, found := k.v1Keeper.GetPermKeeper().GetGroupMemberExtra(ctx, groupInfo.Id, member)
	if found {
		memberExpirationTime = groupMemberExtra.ExpirationTime
	}

	return &v2.QueryHeadGroupMemberResponse{
		Id:             groupMember.Id,
		GroupId:        groupMember.GroupId,
		Member:         groupMember.Member,
		ExpirationTime: memberExpirationTime,
	}, nil
}

// Queries a policy that grants permission to a group
func (k Keeper) QueryPolicyForGroup(ctx context.Context, req *types.QueryPolicyForGroupRequest) (*types.QueryPolicyForGroupResponse, error) {
	return k.v1Keeper.QueryPolicyForGroup(ctx, req)
}

// Queries a policy by policy id
func (k Keeper) QueryPolicyById(ctx context.Context, req *types.QueryPolicyByIdRequest) (*types.QueryPolicyByIdResponse, error) {
	return k.v1Keeper.QueryPolicyById(ctx, req)
}

// Queries lock fee for storing an object
func (k Keeper) QueryLockFee(ctx context.Context, req *types.QueryLockFeeRequest) (*types.QueryLockFeeResponse, error) {
	return k.v1Keeper.QueryLockFee(ctx, req)
}

// Queries a bucket extra info (with gvg bindings and price time) with specify name.
func (k Keeper) HeadBucketExtra(ctx context.Context, req *types.QueryHeadBucketExtraRequest) (*types.QueryHeadBucketExtraResponse, error) {
	return k.v1Keeper.HeadBucketExtra(ctx, req)
}

// Queries whether read and storage prices changed for the bucket.
func (k Keeper) QueryIsPriceChanged(ctx context.Context, req *types.QueryIsPriceChangedRequest) (*types.QueryIsPriceChangedResponse, error) {
	return k.v1Keeper.QueryIsPriceChanged(ctx, req)
}
