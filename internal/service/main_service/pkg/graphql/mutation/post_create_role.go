package mutation

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/entity"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func PostCreateRoleMutation(
	outputTypes map[string]*graphql.Object,
	db *gorm.DB,
	rolesRepo repository.RolesRepositoryInterface,
	permissionsRepo repository.PermissionsRepositoryInterface,
	rolePermissionsRepo repository.RolePermissionsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "post_create_role",
		Type: outputTypes["role"],
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			input := utils.GetSubMap(params.Source, "role")
			roleName, existRoleName := input["name"].(string)
			if !existRoleName {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"name": utils.ErrorInputFail,
				}, errors.New("name input fail"))
			}

			_, err := rolesRepo.TakeByConditions(false, map[string]interface{}{
				"name": roleName,
			})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if err == nil {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"name": utils.ErrorDuplicate,
				}, errors.New(fmt.Sprintf("roles has name: %s already exist", roleName)))
			}

			permissionsInput := utils.GetSubSliceMap(input, "permissions")
			var inputPermissionIDList []int
			seen := make(map[int]bool)
			for i, permissions := range permissionsInput {
				id := utils.GetSubInteger(permissions, "id")
				if id == nil {
					return nil, wraperror.NewValidationError(map[string]interface{}{
						fmt.Sprintf("role.permissions[%d].id", i): utils.ErrorInputFail,
					}, errors.New(fmt.Sprintf("role.permissions[%d].id is not integer", i)))
				}
				if !seen[*id] {
					seen[*id] = true
					inputPermissionIDList = append(inputPermissionIDList, *id)
				}
			}
			permissionIDList, err := permissionsRepo.PluckIDByConditions(map[string]interface{}{
				"id": inputPermissionIDList,
			})
			if err != nil {
				return nil, err
			}
			msg, ok := utils.CompareIntSlice(inputPermissionIDList, permissionIDList, `roles.permissions[%d].id`)
			if !ok {
				return nil, wraperror.NewValidationError(msg, errors.New("some permission ids not exists"))
			}

			var role entity.Roles
			err = utils.Transaction(db, func(tx *gorm.DB) error {
				role, err = rolesRepo.CreateWithTx(tx, utils.GetOnlyScalar(input))
				if err != nil {
					return err
				}
				rolePermissionsAttr := make([]map[string]interface{}, 0)
				for _, id := range permissionIDList {
					rolePermissionsAttr = append(rolePermissionsAttr, map[string]interface{}{
						"role_id":       role.ID,
						"permission_id": id,
					})
				}
				_, err = rolePermissionsRepo.BatchCreateWithTx(tx, rolePermissionsAttr)
				if err != nil {
					return err
				}

				role, err = rolesRepo.TakePreloadByConditionsWithTx(tx, map[string]interface{}{
					"id": role.ID,
				})
				return err
			})
			return role, err
		},
	}
}
