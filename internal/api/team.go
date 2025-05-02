package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/open-move/intercord/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
	baseURL     string
}

func NewTeamHandler(teamService *services.TeamService, baseURL string) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		baseURL:     baseURL,
	}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var input services.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input: " + err.Error()})
		return
	}

	userID := c.GetInt64("userID")
	team, err := h.teamService.Create(c.Request.Context(), input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetTeams(c *gin.Context) {
	userID := c.GetInt64("userID")
	teams, err := h.teamService.GetTeamsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid team ID"})
		return
	}

	team, err := h.teamService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Team not found"})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) InviteToTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid team ID"})
		return
	}

	var input services.InviteToTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input: " + err.Error()})
		return
	}

	input.TeamID = id
	userID := c.GetInt64("userID")
	err = h.teamService.InviteToTeam(c.Request.Context(), input, userID, h.baseURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Invitation sent"})
}

func (h *TeamHandler) JoinTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid team ID"})
		return
	}

	userID := c.GetInt64("userID")
	err = h.teamService.JoinTeam(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Joined team successfully"})
}

func (h *TeamHandler) LeaveTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid team ID"})
		return
	}

	userID := c.GetInt64("userID")
	err = h.teamService.LeaveTeam(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Left team successfully"})
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid team ID"})
		return
	}

	userID := c.GetInt64("userID")
	err = h.teamService.DeleteTeam(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Team deleted successfully"})
}
