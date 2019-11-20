package model

var (
	userPrefix           = "user_"
	messagePrefix        = "message_"
	groupPrefix          = "group_"
	inGroupPrefix        = "ingroup_"
)

func generateDBid (prefix string, id []byte) []byte {
	return append([]byte(prefix), id...)
}

