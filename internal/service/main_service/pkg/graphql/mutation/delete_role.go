package mutation

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"tungnt/emmployee_manage/internal/modules/authorizations/pkg/repository"
	"tungnt/emmployee_manage/pkg/share/utils"
	"tungnt/emmployee_manage/pkg/share/wraperror"
)

func DeleteRoleMutation(
	voidType *graphql.Scalar,
	db *gorm.DB,
	rolesRepo repository.RolesRepositoryInterface,
	rolePermissionsRepo repository.RolePermissionsRepositoryInterface,
	userRolesRepo repository.UserRolesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Name: "delete_role",
		Type: voidType,
		Args: map[string]*graphql.ArgumentConfig{
			"role_id": {
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			roleID := params.Args["role_id"].(int)
			_, err := rolesRepo.TakeByConditions(false, map[string]interface{}{
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

			err = utils.Transaction(db, func(tx *gorm.DB) error {
				err := userRolesRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
					"role_id": roleID,
				})
				if err != nil {
					return err
				}

				err = rolePermissionsRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
					"role_id": roleID,
				})
				if err != nil {
					return err
				}

				return rolesRepo.DeleteByConditionsWithTx(tx, map[string]interface{}{
					"id": roleID,
				})
			})

			return nil, err
		},
	}
}
