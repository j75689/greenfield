package v2

import (
	types "github.com/bnb-chain/greenfield/x/storage/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&types.MsgCreateBucket{}, "storage/CreateBucket", nil)
	cdc.RegisterConcrete(&types.MsgDeleteBucket{}, "storage/DeleteBucket", nil)
	cdc.RegisterConcrete(&types.MsgCreateObject{}, "storage/CreateObject", nil)
	cdc.RegisterConcrete(&types.MsgSealObject{}, "storage/SealObject", nil)
	cdc.RegisterConcrete(&types.MsgRejectSealObject{}, "storage/RejectSealObject", nil)
	cdc.RegisterConcrete(&types.MsgDeleteObject{}, "storage/DeleteObject", nil)
	cdc.RegisterConcrete(&MsgCreateGroup{}, "storage/v2/CreateGroup", nil)
	cdc.RegisterConcrete(&types.MsgDeleteGroup{}, "storage/DeleteGroup", nil)
	cdc.RegisterConcrete(&MsgUpdateGroupMember{}, "storage/v2/UpdateGroupMember", nil)
	cdc.RegisterConcrete(&types.MsgUpdateGroupExtra{}, "storage/UpdateGroupExtra", nil)
	cdc.RegisterConcrete(&types.MsgLeaveGroup{}, "storage/LeaveGroup", nil)
	cdc.RegisterConcrete(&types.MsgCopyObject{}, "storage/CopyObject", nil)
	cdc.RegisterConcrete(&types.MsgUpdateBucketInfo{}, "storage/UpdateBucketInfo", nil)
	cdc.RegisterConcrete(&types.MsgCancelCreateObject{}, "storage/CancelCreateObject", nil)
	cdc.RegisterConcrete(&types.MsgDeletePolicy{}, "storage/DeletePolicy", nil)
	cdc.RegisterConcrete(&types.MsgMigrateBucket{}, "storage/MigrateBucket", nil)
	cdc.RegisterConcrete(&types.MsgCompleteMigrateBucket{}, "storage/CompleteMigrateBucket", nil)
	cdc.RegisterConcrete(&types.MsgCancelMigrateBucket{}, "storage/CancelMigrateBucket", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgCreateBucket{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgDeleteBucket{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgUpdateBucketInfo{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgMirrorBucket{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgDiscontinueBucket{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgCreateObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgSealObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgRejectSealObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgCopyObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgDeleteObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgCancelCreateObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgMirrorObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgDiscontinueObject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgUpdateObjectInfo{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateGroup{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgDeleteGroup{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateGroupMember{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgUpdateGroupExtra{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgLeaveGroup{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgMirrorGroup{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgPutPolicy{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgDeletePolicy{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgUpdateParams{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgMigrateBucket{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgCompleteMigrateBucket{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&types.MsgCancelMigrateBucket{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
