package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"main/forms"
	"main/middleware"
	"main/models"
	"main/repository"
	"main/utils"
	"net/http"
)

func renderResponse[Model any, Form any](model *Model, includeStatus bool) (interface{}, error) {
	var form Form
	err := copier.Copy(&form, model)
	if err != nil {
		return nil, err
	}

	if includeStatus {
		err = copier.Copy(&form, forms.Success)
		if err != nil {
			return nil, err
		}
	}
	return form, nil
}

func renderResponses[Model any, ItemForm any, ItemsForm forms.ListResponse](items []*Model) (interface{}, error) {
	renderedItems := make([]*ItemForm, 0, len(items))
	for _, item := range items {
		renderedItem, err := renderResponse[Model, ItemForm](item, true)
		if err != nil {
			return nil, err
		}
		result := renderedItem.(ItemForm)
		renderedItems = append(renderedItems, &result)
	}

	var form ItemsForm
	utils.SetField(&form, form.ListField(), renderedItems)
	err := copier.Copy(&form, forms.Success)
	if err != nil {
		return nil, err
	}
	return form, nil
}

func setUser[Model any](model *Model, user *models.User) (ok bool) {
	var other interface{} = model
	settable, ok := other.(models.UserSettable)
	if !ok {
		return ok
	}
	settable.SetUser(user)
	return true
}

func handleCreateUserModel[Model models.UserGettable, Form any, ResponseForm any](
	parseForm func(*gin.Context) (*Form, error),
	createModel func(*Model) (*Model, error),
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		form, err := parseForm(ctx)
		if err != nil {
			return
		}

		user, err := middleware.GetUserFromRequest(ctx)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}

		var model Model
		err = copier.Copy(&model, form)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}

		ok := setUser[Model](&model, user)
		if !ok {
			ctx.AbortWithStatus(500)
			return
		}

		result, err := createModel(&model)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Internal error",
				Error:  "Could not create",
			})
			return
		}

		response, err := renderResponse[Model, ResponseForm](result, true)
		ctx.JSON(http.StatusCreated, response)
	}
}

func handlePatchModel[Model models.UserGettable, Form any, ResponseForm any](
	parseForm func(*gin.Context) (*Form, error),
	patchModel func(*Model) (*Model, error),
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		form, err := parseForm(ctx)
		if err != nil {
			return
		}

		itemID, exists := middleware.GetFromContext[*models.ModelID](ctx, middleware.ModelIdKey)
		if !exists {
			ctx.AbortWithStatus(500)
			return
		}

		model, err := repository.GetModelByID[Model](*itemID)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}

		err = copier.Copy(&model, form)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}

		item, err := patchModel(model)
		if err != nil {
			ctx.AbortWithStatus(500)
		}

		response, err := renderResponse[Model, ResponseForm](item, true)
		ctx.JSON(http.StatusOK, response)
	}
}

func handleGetUserModel[Model any, ResponseForm any](
	getUserModel func(id models.ModelID) (*Model, error),
) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		modelID, exists := middleware.GetFromContext[*models.ModelID](ctx, middleware.ModelIdKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, forms.ErrorResponse{
				Status: "Bad Request",
				Error:  "Model not found",
			})
			return
		}

		model, err := getUserModel(*modelID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, forms.ErrorResponse{
				Status: "Server error",
				Error:  err.Error(),
			})
			return
		}

		response, err := renderResponse[Model, ResponseForm](model, true)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}

func handleModelFromContext[Model any, ResponseForm any](ctx *gin.Context) {
	model, exists := middleware.GetFromContext[*Model](ctx, middleware.ModelKey)
	if !exists {
		return
	}

	response, err := renderResponse[Model, ResponseForm](model, true)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func handleGetUserModels[Model any, ItemResponse any, ItemsResponse forms.ListResponse](
	getUserModels func(id models.ModelID) ([]*Model, error),
) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user, err := middleware.GetUserFromRequest(ctx)
		if err != nil {
			return
		}

		items, err := getUserModels(user.ID)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}

		response, err := renderResponses[Model, ItemResponse, ItemsResponse](items)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}

func handleDeleteModel[Model any, ItemResponse any](
	deleteModel func(id models.ModelID) (*Model, error),
) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		itemID, exists := middleware.GetFromContext[*models.ModelID](ctx, middleware.ModelIdKey)
		if !exists {
			ctx.AbortWithStatus(500)
			return
		}

		item, err := deleteModel(*itemID)
		if err != nil {
			ctx.AbortWithStatusJSON(500, forms.ErrorResponse{
				Status: "Server error",
				Error:  err.Error(),
			})
			return
		}

		response, err := renderResponse[Model, ItemResponse](item, true)
		ctx.JSON(http.StatusOK, response)
	}
}

func handleGetAllModels[Model any, ItemResponse any, ItemsResponse forms.ListResponse](
	getAllModels func() ([]*Model, error),
) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		items, err := getAllModels()
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}

		response, err := renderResponses[Model, ItemResponse, ItemsResponse](items)
		if err != nil {
			ctx.AbortWithStatus(500)
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}
