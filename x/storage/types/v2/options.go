package v2

import (
	types "github.com/bnb-chain/greenfield/x/storage/types"
)

type UpdateGroupMemberOptions struct {
	SourceType               types.SourceType
	MembersToAdd             []string
	MembersToDelete          []string
	MemberWithExpirationTime []*MsgGroupMember
}
