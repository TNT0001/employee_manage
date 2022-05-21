package container

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/mutation"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/output"
	"tungnt/emmployee_manage/internal/service/main_service/pkg/graphql/query"
	"tungnt/emmployee_manage/pkg/share/keycloak"
)

// NewGraphQLSchema func
func NewGraphQLSchema(
	repositories *RepositoryContainer,
	services *ServicesContainer,
	db *gorm.DB,
	logger *logrus.Logger,
	client *keycloak.Client,
) (graphql.Schema, error) {
	outputTypes := make(map[string]*graphql.Object)
	for _, graphqlObject := range []*graphql.Object{
		output.NewUserType(
			outputTypes,
			repositories.OrganizationContainer.TeamsRepository,
			repositories.UsersContainer.PersonalInfosRepository,
			repositories.JobsContainer.JobsRepository,
			repositories.AuthorizationsContainer.UserRolesRepository,
		),
		output.NewUserListType(outputTypes),
		output.NewUserTokenType(outputTypes),
		output.NewTeamListType(outputTypes),
		output.NewTeamType(outputTypes),
		output.NewPermissionType(outputTypes),
		output.NewPermissionListType(outputTypes),
		output.NewRoleListType(outputTypes),
		output.NewRoleType(outputTypes),
		output.NewjobType(outputTypes),
		output.NewPersonalInfoType(),
		output.NewUserRoleType(outputTypes, repositories.AuthorizationsContainer.RolesRepository),
	} {
		outputTypes[graphqlObject.Name()] = graphqlObject
	}

	voidType := output.NewVoidType()

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"get_users": query.GetUsersQuery(
					outputTypes,
					repositories.UsersContainer.UsersRepository,
				),
				"get_teams": query.GetTeamsQuery(
					outputTypes,
					repositories.OrganizationContainer.TeamsRepository,
				),
				"get_permissions": query.GetPermissionsQuery(
					outputTypes,
					repositories.AuthorizationsContainer.PermissionsRepository,
				),
				"get_roles": query.GetRolesQuery(
					outputTypes,
					repositories.AuthorizationsContainer.RolesRepository,
				),
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"post_admin_user_login": mutation.PostAdminUserLoginMutation(
					outputTypes,
					client,
				),
				"post_user_login": mutation.PostUserLoginMutation(
					outputTypes,
					client,
				),
				"post_user_login_otp": mutation.PostUserLoginOTPMutation(
					outputTypes,
					client,
				),
				"post_create_team": mutation.PostCreateTeamMutation(
					outputTypes,
					client,
					repositories.OrganizationContainer.TeamsRepository,
				),
				"post_create_permission": mutation.PostCreatePermissionsMutation(
					outputTypes,
					services.AuthorizationsContainer.PermissionsService,
					repositories.AuthorizationsContainer.PermissionsRepository,
				),
				"post_create_role": mutation.PostCreateRoleMutation(
					outputTypes,
					db,
					repositories.AuthorizationsContainer.RolesRepository,
					repositories.AuthorizationsContainer.PermissionsRepository,
					repositories.AuthorizationsContainer.RolePermissionsRepository,
				),
				"delete_role": mutation.DeleteRoleMutation(
					voidType,
					db,
					repositories.AuthorizationsContainer.RolesRepository,
					repositories.AuthorizationsContainer.RolePermissionsRepository,
					repositories.AuthorizationsContainer.UserRolesRepository,
				),
				"patch_update_role": mutation.PatchUpdateRoleMutation(
					outputTypes,
					db,
					repositories.AuthorizationsContainer.RolesRepository,
					repositories.AuthorizationsContainer.PermissionsRepository,
					repositories.AuthorizationsContainer.RolePermissionsRepository,
				),
				"post_create_user": mutation.PostCreateUserMutation(
					outputTypes,
					db,
					client,
					repositories.OrganizationContainer.TeamsRepository,
					services.UserContainer.UsersService,
					repositories.UsersContainer.UsersRepository,
					repositories.AuthorizationsContainer.UserRolesRepository,
					repositories.JobsContainer.JobsRepository,
					repositories.UsersContainer.PersonalInfosRepository,
				),
			},
		}),
	})
}
