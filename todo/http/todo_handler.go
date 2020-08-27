package http

import (
	"context"
	"log"
	"net/http"

	"github.com/kemalnw/todo-with-go/todo/model"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ResponseSucces represent the reseponse error struct
type ResponseSucces struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TodoHandler represent the httphandler for todo
type TodoHandler struct {
	Db *mongo.Database
}

// NewTodoHandler will initialize the todo endpoint
func NewTodoHandler(e *echo.Echo, Db *mongo.Database) {
	handler := &TodoHandler{Db}

	e.GET("/todos", handler.fetchTodo)
	e.POST("/todos", handler.store)
	e.GET("/todos/:id", handler.getByID)
	e.DELETE("/todos/:id", handler.delete)
}

func (t *TodoHandler) fetchTodo(c echo.Context) error {
	var todos []*model.Todo

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	cursor, err := t.Db.Collection("todo").Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	for cursor.Next(context.Background()) {
		var todo model.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, &todo)
	}

	return c.JSON(http.StatusOK, ResponseSucces{Message: "successful retrieve list of todo", Data: todos})
}

func (t *TodoHandler) store(c echo.Context) error {
	var todo model.Todo
	err := c.Bind(&todo)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	todo.AddTimeStamps()

	r, err := t.Db.Collection("todo").InsertOne(ctx, todo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	todo.ID = r.InsertedID.(primitive.ObjectID)

	return c.JSON(http.StatusCreated, ResponseSucces{Message: "todo was created", Data: todo})
}

func (t *TodoHandler) getByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	todo := model.Todo{}
	err = t.Db.Collection("todo").FindOne(ctx, bson.M{"_id": id}).Decode(&todo)
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseSucces{Message: "successful retrieve todo", Data: todo})
}

func (t *TodoHandler) delete(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if r := t.Db.Collection("todo").FindOneAndDelete(ctx, bson.M{"_id": id}); r.Err() != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseSucces{Message: "successful remove todo", Data: nil})
}
