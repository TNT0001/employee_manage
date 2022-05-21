package query

import (
	"github.com/graphql-go/graphql"
	"tungnt/emmployee_manage/internal/modules/organizations/pkg/repository"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/output"
)

func GetTeamsQuery(
	outputTypes map[string]*graphql.Object,
	teamsRepository repository.TeamsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "get_teams",
		Type: outputTypes["team_list"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			teams, err := teamsRepository.GetByConditions(nil)
			if err != nil {
				return nil, err
			}
			total, err := teamsRepository.CountByConditions(nil)
			if err != nil {
				return nil, err
			}

			return output.TeamList{
				Total: total,
				List:  teams,
			}, nil
		},
	}
}
