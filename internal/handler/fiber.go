package handler

import (
	"fmt"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db"
	"github.com/gofiber/fiber/v2"

	"github.com/google/uuid"
)

var dbConn *db.DB

func init() {
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}
	dbConn = conn
}

func buildRespone(article interface{}, nextPageKey string) fiber.Map {
	if article == nil {
		return fiber.Map{
			"nextPageKey": nextPageKey,
		}
	}
	return fiber.Map{
		"article":     article,
		"nextPageKey": nextPageKey,
	}
}

func GetHead() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		key := c.Params("key")
		u, err := uuid.Parse(key)
		if err != nil {
			err = fmt.Errorf("invalid key: %s", key)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		pageNode, err := dbConn.GetPageListBegin(u)
		if err != nil {
			err := fmt.Errorf("failed to get page")
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		response := buildRespone(nil, pageNode.Key.String())
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func GetPage() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		key := c.Params("key")
		u, err := uuid.Parse(key)
		if err != nil {
			err = fmt.Errorf("invalid key: %s", key)
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		pageNode, err := dbConn.GetPageNodeByKey(u)
		if err != nil {
			err := fmt.Errorf("failed to get page")
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		article, err := dbConn.GetPageByID(pageNode.PageID)
		if err != nil {
			err := fmt.Errorf("failed to get page")
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		response := buildRespone(article, pageNode.NextKey.String())
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
