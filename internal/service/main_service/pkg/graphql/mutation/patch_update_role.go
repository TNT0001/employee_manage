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

func PatchUpdateRoleMutation(
	outputTypes map[string]*graphql.Object,
	db *gorm.DB,
	rolesRepo repository.RolesRepositoryInterface,
	permissionsRepo repository.PermissionsRepositoryInterface,
	rolePermissionsRepo repository.RolePermissionsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "patch_update_role",
		Type: outputTypes["role"],
		Args: map[string]*graphql.ArgumentConfig{
			"role_id": {
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			roleID := params.Args["role_id"].(int)
			curRole, err := rolesRepo.TakeByConditions(false, map[string]interface{}{
				"id": roleID,
			})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, wraperror.NewValidationError(map[string]interface{}{
					"role_id": utils.ErrorInputFail,
				}, errors.New(fmt.Sprintf("role with id: %d not exists", roleID)))
			}

			input := utils.GetSubMap(params.Source, "role")
			roleAttrUpdate := make(map[string]interface{})
			roleName, existRoleName := input["name"].(string)
			if existRoleName && roleName != curRole.Name {
				_, err = rolesRepo.TakeByConditions(false, map[string]interface{}{
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
				roleAttrUpdate["name"] = roleName
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
			var permissionIDList []int
			if len(inputPermissionIDList) > 0 {
				permissionIDList, err = permissionsRepo.PluckIDByConditions(map[string]interface{}{
					"id": inputPermissionIDList,
				})
				if err != nil {
					return nil, err
				}
				msg, ok := utils.CompareIntSlice(inputPermissionIDList, permissionIDList, `roles.permissions[%d].id`)
				if !ok {
					return nil, wraperror.NewValidationError(msg, errors.New("some permission ids not exists"))
				}
			}
			var role entity.Roles
			err = utils.Transaction(db, func(tx *gorm.DB) error {
				err := rolePermissionsRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
					"role_id": roleID,
				})
				if err != nil {
					return err
				}

				role, err = rolesRepo.UpdateWithTx(tx, curRole, roleAttrUpdate)
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
				if len(rolePermissionsAttr) > 0 {
					_, err = rolePermissionsRepo.BatchCreateWithTx(tx, rolePermissionsAttr)
					if err != nil {
						return err
					}
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
