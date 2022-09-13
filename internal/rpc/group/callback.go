package group

import (
	pbGroup "micro_servers/pkg/proto/group"
)

func callbackBeforeCreateGroup(req *pbGroup.CreateGroupReq) (bool, error) {
	return true, nil
}

func callbackAfterCreateGroup(req *pbGroup.CreateGroupReq) {

}
