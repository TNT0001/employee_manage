package query

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/users/pkg/repository"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/output"
)

func GetUsersQuery(
	outputTypes map[string]*graphql.Object,
	userRepository repository.UsersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "get_users",
		Type: outputTypes["user_list"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			users, err := userRepository.GetUsers(nil)
			if err != nil {
				return nil, err
			}
			total, err := userRepository.CountUsers(nil)
			if err != nil {
				return nil, err
			}

			return output.UsersList{
				Total: total,
				List:  users,
			}, nil
		},
	}
}
