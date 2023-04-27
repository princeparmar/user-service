package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/princeparmar/contact_manager/repositories"
	"github.com/princeparmar/go-helpers/clienthelper"
	"github.com/princeparmar/go-helpers/context"
)

// Role defines a struct for role data.
type Role struct {
	ID   int
	Name string `json:"name"`
}

// createRoleModel maps Role to Role model.
func createRoleModel(r *Role) *repositories.Role {
	return &repositories.Role{
		ID:   r.ID,
		Name: r.Name,
	}
}

// ParseRequest parses the HTTP request and extracts any relevant data into the Role object.
func (r *Role) ParseRequest(ctx context.IContext, w http.ResponseWriter, req *http.Request) error {
	if req.Method == http.MethodPost {
		// Read the request body
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}

		// Unmarshal the request body into the Role object
		err = json.Unmarshal(body, r)
		if err != nil {
			return err
		}
	}

	// Parse ID from the query parameter
	id := req.URL.Query().Get("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id in query")
	}

	r.ID = i

	return nil
}

// ValidateRequest validates the data in the Role object and returns any errors that occur during validation.
func (r *Role) ValidateRequest(ctx context.IContext) error {
	// TODO: add validation if needed
	return nil
}

// CreateRoleExecutor defines an APIExecutor for creating a new role.
type CreateRoleExecutor struct {
	Role
	clienthelper.BaseAPIExecutor
	RoleRepo repositories.RoleRepository
}

// NewCreateRoleExecutor returns a new instance of CreateRoleExecutor.
func NewCreateRoleExecutor(repo repositories.RoleRepository) clienthelper.APIExecutor {
	return &CreateRoleExecutor{
		RoleRepo: repo,
	}
}

// Controller executes the business logic for creating a new role and returns the created role
// and any errors that occur during execution.
func (e *CreateRoleExecutor) Controller(ctx context.IContext) (interface{}, error) {
	role := createRoleModel(&e.Role)
	err := e.RoleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRoleExecutor defines an APIExecutor for deleting a role by ID.
type DeleteRoleExecutor struct {
	Role
	clienthelper.BaseAPIExecutor
	RoleRepo repositories.RoleRepository
}

// NewDeleteRoleExecutor returns a new instance of DeleteRoleExecutor.
func NewDeleteRoleExecutor(repo repositories.RoleRepository) clienthelper.APIExecutor {
	return &DeleteRoleExecutor{
		RoleRepo: repo,
	}
}

// Controller executes the business logic for deleting a role by ID and returns any errors that occur during execution.
func (e *DeleteRoleExecutor) Controller(ctx context.IContext) (interface{}, error) {
	id := e.Role.ID
	err := e.RoleRepo.Delete(id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateRoleExecutor defines an APIExecutor for updating a role by ID.
type UpdateRoleExecutor struct {
	Role
	clienthelper.BaseAPIExecutor
	RoleRepo repositories.RoleRepository
}

// NewUpdateRoleExecutor returns a new instance of UpdateRoleExecutor.
func NewUpdateRoleExecutor(repo repositories.RoleRepository) clienthelper.APIExecutor {
	return &UpdateRoleExecutor{
		RoleRepo: repo,
	}
}

// Controller executes the business logic for updating a role by ID and returns the updated role
// and any errors that occur during execution.
func (e *UpdateRoleExecutor) Controller(ctx context.IContext) (interface{}, error) {
	role := createRoleModel(&e.Role)
	err := e.RoleRepo.Update(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleExecutor defines an APIExecutor for getting a role by ID.
type GetRoleExecutor struct {
	Role
	clienthelper.BaseAPIExecutor
	RoleRepo repositories.RoleRepository
}

// NewGetRoleExecutor returns a new instance of GetRoleExecutor.
func NewGetRoleExecutor(repo repositories.RoleRepository) clienthelper.APIExecutor {
	return &GetRoleExecutor{
		RoleRepo: repo,
	}
}

// Controller executes the business logic for getting a role by ID and returns the role
// and any errors that occur during execution.
func (e *GetRoleExecutor) Controller(ctx context.IContext) (interface{}, error) {
	id := e.Role.ID
	role, err := e.RoleRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetAllRolesExecutor defines an APIExecutor for getting all roles.
type GetAllRolesExecutor struct {
	clienthelper.BaseAPIExecutor
	RoleRepo repositories.RoleRepository
}

// NewGetAllRolesExecutor returns a new instance of GetAllRolesExecutor.
func NewGetAllRolesExecutor(repo repositories.RoleRepository) clienthelper.APIExecutor {
	return &GetAllRolesExecutor{
		RoleRepo: repo,
	}
}

// Controller executes the business logic for getting all roles and returns the roles
// and any errors that occur during execution.
func (e *GetAllRolesExecutor) Controller(ctx context.IContext) (interface{}, error) {
	return e.RoleRepo.GetAll()
}
