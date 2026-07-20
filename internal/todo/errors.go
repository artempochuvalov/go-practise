package todo

import "errors"

var EmptyTitleError = errors.New("title cannot be empty")
var TodoNotFound = errors.New("todo not found")
