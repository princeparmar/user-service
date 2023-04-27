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

// Access defines a struct for access data.
type Access struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// createAccessModel maps Access to Access model.
func createAccessModel(a *Access) *repositories.Access {
	return &repositories.Access{
		ID:   a.ID,
		Name: a.Name,
	}
}

// ParseRequest parses the HTTP request and extracts any relevant data into the Access object.
func (a *Access) ParseRequest(ctx context.IContext, w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		// Unmarshal the request body into the Access object
		err = json.Unmarshal(body, a)
		if err != nil {
			return err
		}
	}

	// Parse ID from the query parameter
	id := r.URL.Query().Get("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid id in query")
	}

	a.ID = i

	return nil
}

// ValidateRequest validates the data in the Access object and returns any errors that occur during validation.
func (a *Access) ValidateRequest(ctx context.IContext) error {
	if a.Name == "" {
		return errors.New("access name is required")
	}
	return nil
}

// CreateAccessExecutor defines an APIExecutor for creating a new access.
type CreateAccessExecutor struct {
	Access
	clienthelper.BaseAPIExecutor
	AccessRepo repositories.AccessRepository
}

// NewCreateAccessExecutor returns a new instance of CreateAccessExecutor.
func NewCreateAccessExecutor(repo repositories.AccessRepository) clienthelper.APIExecutor {
	return &CreateAccessExecutor{
		AccessRepo: repo,
	}
}

// Controller executes the business logic for creating a new access and returns the created access
// and any errors that occur during execution.
func (e *CreateAccessExecutor) Controller(ctx context.IContext) (interface{}, error) {
	access := createAccessModel(&e.Access)
	err := e.AccessRepo.Create(access)
	if err != nil {
		return nil, err
	}

	return access, nil
}

// UpdateAccessExecutor defines an APIExecutor for updating an access by ID.
type UpdateAccessExecutor struct {
	Access
	clienthelper.BaseAPIExecutor
	AccessRepo repositories.AccessRepository
}

// NewUpdateAccessExecutor returns a new instance of UpdateAccessExecutor.
func NewUpdateAccessExecutor(repo repositories.AccessRepository) clienthelper.APIExecutor {
	return &UpdateAccessExecutor{
		AccessRepo: repo,
	}
}

// Controller executes the business logic for updating an access by ID and returns the updated access
// and any errors that occur during execution.
func (e *UpdateAccessExecutor) Controller(ctx context.IContext) (interface{}, error) {
	access := createAccessModel(&e.Access)
	err := e.AccessRepo.Update(access)
	if err != nil {
		return nil, err
	}

	return access, nil
}

// DeleteAccessExecutor defines an APIExecutor for deleting an access mode by ID.
type DeleteAccessExecutor struct {
	Access
	clienthelper.BaseAPIExecutor
	AccessRepo repositories.AccessRepository
}

// NewDeleteAccessExecutor returns a new instance of DeleteAccessExecutor.
func NewDeleteAccessExecutor(repo repositories.AccessRepository) clienthelper.APIExecutor {
	return &DeleteAccessExecutor{
		AccessRepo: repo,
	}
}

// Controller executes the business logic for deleting an access mode by ID and returns any errors that occur during execution.
func (e *DeleteAccessExecutor) Controller(ctx context.IContext) (interface{}, error) {
	id := e.Access.ID
	err := e.AccessRepo.Delete(id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAccessExecutor defines an APIExecutor for getting an access by ID.
type GetAccessExecutor struct {
	Access
	clienthelper.BaseAPIExecutor
	AccessRepo repositories.AccessRepository
}

// NewGetAccessExecutor returns a new instance of GetAccessExecutor.
func NewGetAccessExecutor(repo repositories.AccessRepository) clienthelper.APIExecutor {
	return &GetAccessExecutor{
		AccessRepo: repo,
	}
}

// Controller executes the business logic for getting an access by ID and returns the access
// and any errors that occur during execution.
func (e *GetAccessExecutor) Controller(ctx context.IContext) (interface{}, error) {
	id := e.Access.ID
	access, err := e.AccessRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return access, nil
}

// GetAllAccessesExecutor defines an APIExecutor for getting all accesses.
type GetAllAccessesExecutor struct {
	clienthelper.BaseAPIExecutor
	AccessRepo repositories.AccessRepository
}

// NewGetAllAccessesExecutor returns a new instance of GetAllAccessesExecutor.
func NewGetAllAccessesExecutor(repo repositories.AccessRepository) clienthelper.APIExecutor {
	return &GetAllAccessesExecutor{
		AccessRepo: repo,
	}
}

// Controller executes the business logic for getting all accesses and returns the accesses
// and any errors that occur during execution.
func (e *GetAllAccessesExecutor) Controller(ctx context.IContext) (interface{}, error) {
	return e.AccessRepo.GetAll()
}
