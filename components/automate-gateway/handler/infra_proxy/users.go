package infra_proxy

import (
	"context"

	gwreq "github.com/chef/automate/api/external/infra_proxy/request"

	gwres "github.com/chef/automate/api/external/infra_proxy/response"
	infra_req "github.com/chef/automate/api/interservice/infra_proxy/request"
	infra_res "github.com/chef/automate/api/interservice/infra_proxy/response"
)

//GetUsersList: fetches a list of an existing users in organization
func (c *InfraProxyServer) GetOrgUsersList(ctx context.Context, r *gwreq.OrgUsers) (*gwres.OrgUsers, error) {
	req := &infra_req.OrgUsers{
		OrgId:    r.OrgId,
		ServerId: r.ServerId,
	}
	res, err := c.client.GetOrgUsersList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &gwres.OrgUsers{
		Users: fromUpstreamOrgUsers(res.Users),
	}, nil
}

func fromUpstreamOrgUsers(users []*infra_res.UsersListItem) []*gwres.UsersListItem {
	us := make([]*gwres.UsersListItem, len(users))

	for i, user := range users {
		us[i] = &gwres.UsersListItem{
			Username: user.GetUsername(),
		}
	}

	return us
}
