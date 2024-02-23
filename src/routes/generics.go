package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"main/forms"
	"main/models"
	"main/repository"
	"main/utils"
	"net/http"
)

func modelPermissionRequired[Model models.UserProtected](ctx *gin.Context) {
	user, err := utils.GetUserFromRequest(ctx)
	if err != nil {
		ctx.Abort()
		return
	}

	modelID, err := utils.ShouldParseID(ctx)
	if err != nil {
		ctx.Abort()
		return
	}

	model, err := repository.GetModelByID[Model](modelID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, forms.ErrorResponse{
			Status: "Not found",
			Error:  "No item found with this ID",
		})
		return
	}

	if (*model).GetUser().ID != user.ID {
		ctx.JSON(http.StatusForbidden, forms.ErrorResponse{
			Status: "Forbidden",
			Error:  "Operations on this model are forbidden",
		})
		ctx.Abort()
		return
	}

	ctx.Set("item", model)
	ctx.Set("itemID", modelID)

	ctx.Next()
}

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

func handleCreateModel[Model any, Form any, ResponseForm any](ctx *gin.Context) {
	form, err := utils.ShouldGetForm[Form](ctx)
	if err != nil {
		return
	}

	user, err := utils.GetUserFromRequest(ctx)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	var model Model
	err = copier.Copy(&model, form)
	if err != nil {
		ctx.AbortWithStatus(500)
	}
	utils.SetField(&model, "UserID", user.ID)

	result, err := repository.CreateModel[Model](&model)
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

func handlePatchModel[Model models.UserProtected, Form any, ResponseForm any](ctx *gin.Context) {
	form, err := utils.ShouldGetForm[Form](ctx)
	if err != nil {
		return
	}

	itemID, exists := utils.GetFromContext[uint](ctx, "itemID")
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

	item, err := repository.PatchModel[Model](model)
	if err != nil {
		ctx.AbortWithStatus(500)
	}

	response, err := renderResponse[Model, ResponseForm](item, true)
	ctx.JSON(http.StatusOK, response)
}

func handleGetUserModel[Model any, ResponseForm any](ctx *gin.Context) {
	item, exists := utils.GetFromContext[*Model](ctx, "item")
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	response, err := renderResponse[Model, ResponseForm](*item, true)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func handleGetUserModels[Model any, ItemResponse any, ItemsResponse forms.ListResponse](ctx *gin.Context) {
	user, err := utils.GetUserFromRequest(ctx)
	if err != nil {
		return
	}

	items, err := repository.GetUserModels[Model](user.ID)
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

func handleDeleteModel[Model any, ItemResponse any](ctx *gin.Context) {
	itemID, exists := utils.GetFromContext[uint](ctx, "itemID")
	if !exists {
		ctx.AbortWithStatus(500)
		return
	}

	item, err := repository.DeleteModel[Model](*itemID)
	if err != nil {
		ctx.AbortWithStatusJSON(500, forms.ErrorResponse{
			Status: "Error",
			Error:  "Could not delete",
		})
		return
	}

	response, err := renderResponse[Model, ItemResponse](item, true)
	ctx.JSON(http.StatusOK, response)
}

func handleGetAllModels[Model any, ItemResponse any, ItemsResponse forms.ListResponse](ctx *gin.Context) {
	items, err := repository.GetAllModels[Model]()
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
