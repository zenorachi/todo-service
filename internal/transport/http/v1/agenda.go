package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/todo-service/internal/entity"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-Jan-02"

func (h *Handler) initAgendaRoutes(api *gin.RouterGroup) {
	agenda := api.Group("/agenda", h.userIdentity)
	{
		agenda.POST("/create", h.createTask)
		agenda.GET("/:task_id", h.getTaskByID)
		agenda.PUT("/set_status", h.setTaskStatus)
		agenda.DELETE("/delete_by_id", h.deleteTaskByID)
		agenda.DELETE("/delete_all", h.deleteUserTasks)
		agenda.GET("/get_all", h.getUserTasks)
		agenda.GET("/get_by_date", h.getTasksByDataAndStatus)
	}
}

/* -- CREATE TASK -- */

type createTaskInput struct {
	Title       string `json:"title"    binding:"required,min=2,max=64"`
	Description string `json:"description"`
	Date        string `json:"date" binding:"required,min=6,max=64"`
	Status      string `json:"status"`
}

type createTaskResponse struct {
	ID int `json:"id"`
}

// @Summary Create task
// @Description create task
// @Tags agenda
// @Accept json
// @Produce json
// @Param input body createTaskInput true "input"
// @Success 201 {object} createTaskResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/create [post]
func (h *Handler) createTask(c *gin.Context) {
	var input createTaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	date, err := time.Parse(dateFormat, input.Date)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidData.Error())
		return
	}

	id, err := h.services.Agenda.CreateTask(c, entity.Task{
		UserID:      c.GetInt(userCtx),
		Title:       input.Title,
		Description: input.Description,
		Date:        date,
		Status:      input.Status,
	})

	if err != nil {
		if errors.Is(err, entity.ErrTaskAlreadyExist) {
			newErrorResponse(c, http.StatusConflict, err.Error())
		} else if errors.Is(err, entity.ErrInvalidStatus) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusCreated, createTaskResponse{ID: id})
}

/* -- GET TASK BY ID-- */

type getTaskByIDResponse struct {
	Task entity.Task `json:"task"`
}

// @Summary Get Task By ID
// @Security Bearer
// @Description getting task by id
// @Tags agenda
// @Produce json
// @Param task_id path int true "Task ID"
// @Success 200 {object} getTaskByIDResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/:task_id [get]
func (h *Handler) getTaskByID(c *gin.Context) {
	paramId := strings.Trim(c.Param("task_id"), "/")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid parameter (id)")
		return
	}

	task, err := h.services.Agenda.GetTaskByID(c, id)
	if err != nil {
		if errors.Is(err, entity.ErrTaskDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusOK, getTaskByIDResponse{Task: task})
}

/* --- SET TASK STATUS --- */

type setTaskStatusInput struct {
	ID     int    `json:"task_id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// @Summary Update Task Status
// @Security Bearer
// @Description updating task status by id
// @Tags agenda
// @Accept json
// @Param input body setTaskStatusInput true "input"
// @Success 204 "No Content"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/set_status [put]
func (h *Handler) setTaskStatus(c *gin.Context) {
	var input setTaskStatusInput
	if c.BindJSON(&input) != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	err := h.services.Agenda.SetTaskStatus(c, input.ID, input.Status)
	if err != nil {
		if errors.Is(err, entity.ErrTaskDoesNotExist) || errors.Is(err, entity.ErrInvalidStatus) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusNoContent, nil)
}

/* --- DELETE TASK BY ID --- */

type deleteTaskByIdInput struct {
	ID int `json:"task_id" binding:"required"`
}

// @Summary Delete Task By ID
// @Security Bearer
// @Description deleting task by id
// @Tags agenda
// @Accept json
// @Param input body deleteTaskByIdInput true "input"
// @Success 204 "No Content"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/delete_by_id [delete]
func (h *Handler) deleteTaskByID(c *gin.Context) {
	var input deleteTaskByIdInput
	if c.BindJSON(&input) != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	err := h.services.Agenda.DeleteTaskByID(c, input.ID)
	if err != nil {
		if errors.Is(err, entity.ErrTaskDoesNotExist) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	newResponse(c, http.StatusNoContent, nil)
}

/* --- DELETE ALL USER TASKS --- */

// @Summary Delete All User Tasks
// @Security Bearer
// @Description deleting all user tasks by user id
// @Tags agenda
// @Success 204 "No Content"
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/delete_all [delete]
func (h *Handler) deleteUserTasks(c *gin.Context) {
	err := h.services.Agenda.DeleteUserTasks(c, c.GetInt(userCtx))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(c, http.StatusNoContent, nil)
}

/* --- GET ALL USER TASKS --- */

type getAllUserTasksResponse struct {
	Tasks []entity.Task `json:"tasks"`
}

// @Summary Get All User Tasks
// @Security Bearer
// @Description getting all user tasks by user id
// @Tags agenda
// @Produce json
// @Success 200 {object} getAllUserTasksResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/get_all [get]
func (h *Handler) getUserTasks(c *gin.Context) {
	tasks, err := h.services.Agenda.GetUserTasks(c, c.GetInt(userCtx))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(c, http.StatusOK, getAllUserTasksResponse{Tasks: tasks})
}

/* --- GET ALL USER TASKS BY DATE --- */

type getAllUserTasksByDataInput struct {
	Date   string `json:"date"`
	Limit  int    `json:"limit"  binding:"required,min=1"`
	Offset int    `json:"offset" binding:"required,min=1"`
}

// @Summary Get All User Tasks
// @Security Bearer
// @Description getting all user tasks by user id and date or status (also support pagination)
// @Tags agenda
// @Accept json
// @Produce json
// @Param input body getAllUserTasksByDataInput true "input"
// @Success 200 {object} getAllUserTasksResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/agenda/get_by_date [get]
func (h *Handler) getTasksByDataAndStatus(c *gin.Context) {
	var input getAllUserTasksByDataInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidInput.Error())
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidPaginationSizes.Error())
		return
	}

	var (
		date   = time.Time{}
		status = c.DefaultQuery("status", "not done")
	)

	if len(input.Date) != 0 {
		date, err = time.Parse(dateFormat, input.Date)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, entity.ErrInvalidData.Error())
			return
		}
	}

	tasks, err := h.services.Agenda.GetByDateAndStatus(c, c.GetInt(userCtx), status, date, input.Limit, (page-1)*input.Offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(c, http.StatusOK, getAllUserTasksResponse{Tasks: tasks})
}
