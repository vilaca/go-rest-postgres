package handler

// TODO username must be unique in db
import (
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pauljamescleary/gomin/pkg/common/models"
)

func (h *Handler) CreateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	name_is_valid := len(u.Name) >= 3
	name_is_valid = name_is_valid && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(u.Name)
	name_is_valid = name_is_valid && regexp.MustCompile(`^[a-zA-Z].*$`).MatchString(u.Name)
	if !name_is_valid {
		return echo.NewHTTPError(http.StatusBadRequest, "username must be alphanumeric, the first character must be a letter and at least 3 characters long")
	}

	password_is_valid := len(u.Password) >= 6
	if !password_is_valid {
		return echo.NewHTTPError(http.StatusBadRequest, "password must be at least 6 characters")
	}

	oldUser, err := h.UserRepo.UserExists(u.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if oldUser {
		return echo.NewHTTPError(http.StatusBadRequest, "username already in database")
	}

	u.ID = uuid.New()
	u.Enabled = true
	newUser, err := h.UserRepo.CreateUser(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, newUser)
}

func (h *Handler) CreateSession(c echo.Context) error {
	login := new(models.Login)
	if err := c.Bind(login); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	ok, err := h.UserRepo.Login(login.Name, login.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "username or password don't match")
	}

	//TODO get millisecs
	start_time := 0
	end_time := 0 + 60*60*2
	// TODO check if sessions exist before
	session := new(models.Session)
	session.Id = uuid.New().String()
	session.UserName = login.Name
	session.Started = start_time
	session.Ends = end_time
	h.UserRepo.CreateSession(session)

	return c.JSON(http.StatusCreated, session)
}

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.UserRepo.GetUser(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, id)
	}

	return c.JSON(http.StatusOK, user)
}
