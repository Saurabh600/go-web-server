package api_controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Saurabh600/go-web-server/internals/config"
	"github.com/Saurabh600/go-web-server/internals/data"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CheckAuth(c *fiber.Ctx) error {
	var bodyData data.EmailAndPassword
	err := c.BodyParser(&bodyData)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte("bad request"))
	}

	dialect := goqu.Dialect("mysql")
	sqlQuery, _, _ := dialect.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"password",
		"password_hash",
		"age",
		"gender",
		"created_at",
		"updated_at",
	).
		From("Users").
		ToSQL()

	db := config.GetDb()
	row := db.QueryRow(sqlQuery)

	var udata data.User
	err = row.Scan(
		&udata.Id,
		&udata.FirstName,
		&udata.LastName,
		&udata.Email,
		&udata.Password,
		&udata.PasswordHash,
		&udata.Age,
		&udata.Gender,
		&udata.CreatedAt,
		&udata.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusBadRequest).
				Send([]byte("user with given email does not exists"))
		}

		return c.Status(http.StatusBadRequest).
			Send([]byte("scan failed!"))

	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(udata.PasswordHash),
		[]byte(bodyData.Password),
	)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(data.JsonResponse{
			Status: false,
			Info:   "invalid credential",
			Data:   []data.User{udata},
		})
	}

	return c.Status(http.StatusOK).JSON(data.JsonResponse{
		Status: true,
		Info:   "valid credential",
		Data:   []data.User{udata},
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	dialect := goqu.Dialect("mysql")
	sqlQuery, _, _ := dialect.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"password",
		"password_hash",
		"age",
		"gender",
		"created_at",
		"updated_at",
	).
		From("Users").
		Order(goqu.I("id").Asc()).
		Limit(100).
		ToSQL()

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(100*time.Millisecond),
	)
	defer cancel()
	db := config.GetDb()
	rows, err := db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"success":   false,
				"message":   "db error",
				"error":     err,
				"sql_query": sqlQuery,
			})
	}

	defer rows.Close()

	var userData []data.User
	for rows.Next() {
		var udata data.User
		err := rows.Scan(
			&udata.Id,
			&udata.FirstName,
			&udata.LastName,
			&udata.Email,
			&udata.Password,
			&udata.PasswordHash,
			&udata.Age,
			&udata.Gender,
			&udata.CreatedAt,
			&udata.UpdatedAt,
		)
		if err != nil {
			return c.Send([]byte(err.Error()))
		}
		userData = append(userData, udata)
	}

	return c.Status(http.StatusOK).JSON(data.JsonResponse{
		Status: true,
		Info:   "user data retrieved successfully",
		Data:   userData,
	})
}

func CreateNewUser(c *fiber.Ctx) error {
	var uPostData data.UserFormData
	err := c.BodyParser(&uPostData)
	if err != nil {
		return c.Send([]byte("invalid body"))
	}

	dialect := goqu.Dialect("mysql")
	hash, _ := bcrypt.GenerateFromPassword([]byte(uPostData.Password), 10)
	insertQuery, _, _ := dialect.Insert("Users").
		Cols("first_name", "last_name", "email", "password", "password_hash", "age", "gender").
		Vals(
			goqu.Vals{
				uPostData.FirstName,
				uPostData.LastName,
				uPostData.Email,
				uPostData.Password,
				string(hash),
				uPostData.Age,
				uPostData.Gender,
			},
		).
		ToSQL()
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(100*time.Millisecond),
	)
	defer cancel()
	db := config.GetDb()
	result, err := db.ExecContext(ctx, insertQuery)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success":   false,
			"message":   "db error",
			"error":     err,
			"sql_query": insertQuery,
		})
	}
	id, _ := result.LastInsertId()
	sqlQuery, _, _ := dialect.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"password",
		"password_hash",
		"age",
		"gender",
		"created_at",
		"updated_at",
	).
		From("Users").
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	ctx, cancel = context.WithDeadline(
		context.Background(),
		time.Now().Add(100*time.Millisecond),
	)
	defer cancel()

	rows, err := db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(fiber.Map{
				"success":   false,
				"message":   "db error",
				"error":     err,
				"sql_query": sqlQuery,
			})
	}

	defer rows.Close()

	var userData []data.User
	for rows.Next() {
		var udata data.User
		err := rows.Scan(
			&udata.Id,
			&udata.FirstName,
			&udata.LastName,
			&udata.Email,
			&udata.Password,
			&udata.PasswordHash,
			&udata.Age,
			&udata.Gender,
			&udata.CreatedAt,
			&udata.UpdatedAt,
		)
		if err != nil {
			return c.Send([]byte(err.Error()))
		}
		userData = append(userData, udata)
	}

	return c.Status(http.StatusOK).JSON(data.JsonResponse{
		Status: true,
		Info:   "new user created successfull",
		Data:   userData,
	})
}
