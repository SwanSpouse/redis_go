package database

type TSet interface {
	SCard()
	SIsMember()
	SMembers()
	SRandMember()
	SPop()
	SRem()
}
