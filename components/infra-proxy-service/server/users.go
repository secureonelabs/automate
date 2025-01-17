package server

import (
	"context"

	"github.com/chef/automate/api/interservice/infra_proxy/request"
	"github.com/chef/automate/api/interservice/infra_proxy/response"
	"github.com/go-chef/chef"
)

// GetUsersList Get a list of all users in an organization
func (s *Server) GetOrgUsersList(ctx context.Context, req *request.OrgUsers) (*response.OrgUsers, error) {
	c, err := s.createClient(ctx, req.OrgId, req.ServerId)
	if err != nil {
		return nil, err
	}
	usersList, err := c.client.Associations.List()
	if err != nil {
		return nil, ParseAPIError(err)
	}

	return &response.OrgUsers{
		Users: fromAPIToListOrgUsers(usersList),
	}, nil

}

func fromAPIToListOrgUsers(list []chef.OrgUserListEntry) []*response.UsersListItem {
	users := make([]*response.UsersListItem, 0)
	for _, user := range list {
		item := &response.UsersListItem{
			Username: user.User.Username,
		}
		users = append(users, item)
	}
	return users
}
