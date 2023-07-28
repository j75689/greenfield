package v2

import (
	"math/big"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	types "github.com/bnb-chain/greenfield/x/storage/types"
)

func DeserializeCrossChainPackage(rawPack []byte, channelId sdk.ChannelID, packageType sdk.CrossChainPackageType) (interface{}, error) {
	if packageType >= 3 {
		return nil, types.ErrInvalidCrossChainPackage
	}

	pack, err := types.DeserializeRawCrossChainPackage(rawPack)
	if err != nil {
		return nil, err
	}

	operationMap, ok := DeserializeFuncMap[channelId][pack.OperationType]
	if !ok {
		return nil, types.ErrInvalidCrossChainPackage
	}

	return operationMap[packageType](pack.Package)
}

var DeserializeFuncMap = map[sdk.ChannelID]map[uint8][3]types.DeserializeFunc{
	types.BucketChannelId: {
		types.OperationMirrorBucket: {
			types.DeserializeMirrorBucketSynPackage,
			types.DeserializeMirrorBucketAckPackage,
			types.DeserializeMirrorBucketSynPackage,
		},
		types.OperationCreateBucket: {
			types.DeserializeCreateBucketSynPackage,
			types.DeserializeCreateBucketAckPackage,
			types.DeserializeCreateBucketSynPackage,
		},
		types.OperationDeleteBucket: {
			types.DeserializeDeleteBucketSynPackage,
			types.DeserializeDeleteBucketAckPackage,
			types.DeserializeDeleteBucketSynPackage,
		},
	},
	types.ObjectChannelId: {
		types.OperationMirrorObject: {
			types.DeserializeMirrorObjectSynPackage,
			types.DeserializeMirrorObjectAckPackage,
			types.DeserializeMirrorObjectSynPackage,
		},
		types.OperationDeleteObject: {
			types.DeserializeDeleteObjectSynPackage,
			types.DeserializeDeleteObjectAckPackage,
			types.DeserializeDeleteObjectSynPackage,
		},
	},
	types.GroupChannelId: {
		types.OperationMirrorGroup: {
			types.DeserializeMirrorGroupSynPackage,
			types.DeserializeMirrorGroupAckPackage,
			types.DeserializeMirrorGroupSynPackage,
		},
		types.OperationCreateGroup: {
			types.DeserializeCreateGroupSynPackage,
			types.DeserializeCreateGroupAckPackage,
			types.DeserializeCreateGroupSynPackage,
		},
		types.OperationDeleteGroup: {
			types.DeserializeDeleteGroupSynPackage,
			types.DeserializeDeleteGroupAckPackage,
			types.DeserializeDeleteGroupSynPackage,
		},
		types.OperationUpdateGroupMember: {
			DeserializeUpdateGroupMemberSynPackage,
			DeserializeUpdateGroupMemberAckPackage,
			DeserializeUpdateGroupMemberSynPackage,
		},
	},
}

type UpdateGroupMemberSynPackage struct {
	Operator         sdk.AccAddress
	GroupId          *big.Int
	OperationType    uint8
	Members          []sdk.AccAddress
	ExtraData        []byte
	MemberExpiration []uint64
}

type UpdateGroupMemberSynPackageStruct struct {
	Operator         common.Address
	GroupId          *big.Int
	OperationType    uint8
	Members          []common.Address
	ExtraData        []byte
	MemberExpiration []uint64
}

var (
	updateGroupMemberSynPackageType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "Operator", Type: "address"},
		{Name: "GroupId", Type: "uint256"},
		{Name: "OperationType", Type: "uint8"},
		{Name: "Members", Type: "address[]"},
		{Name: "ExtraData", Type: "bytes"},
		{Name: "MemberExpiration", Type: "uint64[]"},
	})

	updateGroupMemberSynPackageArgs = abi.Arguments{
		{Type: updateGroupMemberSynPackageType},
	}
)

func (p UpdateGroupMemberSynPackage) GetMembers() []string {
	members := make([]string, 0, len(p.Members))
	for _, member := range p.Members {
		members = append(members, member.String())
	}
	return members
}

func (p UpdateGroupMemberSynPackage) GetMemberExpiration() []time.Time {
	memberExpiration := make([]time.Time, 0, len(p.MemberExpiration))
	for _, expiration := range p.MemberExpiration {
		memberExpiration = append(memberExpiration, time.Unix(int64(expiration), 0))
	}
	return memberExpiration
}

func (p UpdateGroupMemberSynPackage) ValidateBasic() error {
	if p.OperationType != types.OperationAddGroupMember && p.OperationType != types.OperationDeleteGroupMember {
		return types.ErrInvalidOperationType
	}

	if p.Operator.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if p.GroupId == nil || p.GroupId.Cmp(big.NewInt(0)) < 0 {
		return types.ErrInvalidId
	}

	for _, member := range p.Members {
		if member.Empty() {
			return sdkerrors.ErrInvalidAddress
		}
	}

	if len(p.Members) != len(p.MemberExpiration) {
		return types.ErrInvalidGroupMemberExpiration
	}

	return nil
}

func DeserializeUpdateGroupMemberSynPackage(serializedPackage []byte) (interface{}, error) {
	unpacked, err := updateGroupMemberSynPackageArgs.Unpack(serializedPackage)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidCrossChainPackage, "deserialize delete bucket ack package failed")
	}

	unpackedStruct := abi.ConvertType(unpacked[0], UpdateGroupMemberSynPackageStruct{})
	pkgStruct, ok := unpackedStruct.(UpdateGroupMemberSynPackageStruct)
	if !ok {
		return nil, errors.Wrapf(types.ErrInvalidCrossChainPackage, "reflect delete bucket ack package failed")
	}

	totalMember := len(pkgStruct.Members)
	members := make([]sdk.AccAddress, totalMember)
	for i, member := range pkgStruct.Members {
		members[i] = member.Bytes()
	}
	tp := UpdateGroupMemberSynPackage{
		pkgStruct.Operator.Bytes(),
		pkgStruct.GroupId,
		pkgStruct.OperationType,
		members,
		pkgStruct.ExtraData,
		pkgStruct.MemberExpiration,
	}
	return &tp, nil
}

type UpdateGroupMemberAckPackage struct {
	Status        uint8
	Id            *big.Int
	Operator      sdk.AccAddress
	OperationType uint8
	Members       []sdk.AccAddress
	ExtraData     []byte
}

type UpdateGroupMemberAckPackageStruct struct {
	Status        uint8
	Id            *big.Int
	Operator      common.Address
	OperationType uint8
	Members       []common.Address
	ExtraData     []byte
}

var (
	updateGroupMemberAckPackageType, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "Status", Type: "uint8"},
		{Name: "Id", Type: "uint256"},
		{Name: "Operator", Type: "address"},
		{Name: "OperationType", Type: "uint8"},
		{Name: "Members", Type: "address[]"},
		{Name: "ExtraData", Type: "bytes"},
	})

	updateGroupMemberAckPackageArgs = abi.Arguments{
		{Type: updateGroupMemberAckPackageType},
	}
)

func DeserializeUpdateGroupMemberAckPackage(serializedPackage []byte) (interface{}, error) {
	unpacked, err := updateGroupMemberAckPackageArgs.Unpack(serializedPackage)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidCrossChainPackage, "deserialize update group member ack package failed")
	}

	unpackedStruct := abi.ConvertType(unpacked[0], UpdateGroupMemberAckPackageStruct{})
	pkgStruct, ok := unpackedStruct.(UpdateGroupMemberAckPackageStruct)
	if !ok {
		return nil, errors.Wrapf(types.ErrInvalidCrossChainPackage, "reflect update group member ack package failed")
	}

	totalMember := len(pkgStruct.Members)
	members := make([]sdk.AccAddress, totalMember)
	for i, member := range pkgStruct.Members {
		members[i] = member.Bytes()
	}
	tp := UpdateGroupMemberAckPackage{
		pkgStruct.Status,
		pkgStruct.Id,
		pkgStruct.Operator.Bytes(),
		pkgStruct.OperationType,
		members,
		pkgStruct.ExtraData,
	}
	return &tp, nil
}

func (p UpdateGroupMemberAckPackage) MustSerialize() []byte {
	totalMember := len(p.Members)
	members := make([]common.Address, totalMember)
	for i, member := range p.Members {
		members[i] = common.BytesToAddress(member)
	}

	encodedBytes, err := updateGroupMemberAckPackageArgs.Pack(&UpdateGroupMemberAckPackageStruct{
		p.Status,
		types.SafeBigInt(p.Id),
		common.BytesToAddress(p.Operator),
		p.OperationType,
		members,
		p.ExtraData,
	})
	if err != nil {
		panic("encode delete group ack package error")
	}
	return encodedBytes
}
