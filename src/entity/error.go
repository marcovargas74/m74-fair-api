package entity

import "errors"

//ErrNotFound not found
var ErrNotFound = errors.New("Not found")

//ErrInvalidEntity invalid entity
var ErrInvalidEntity = errors.New("Invalid entity")

//ErrCannotBeDeleted cannot be deleted
var ErrCannotBeDeleted = errors.New("Cannot Be Deleted")

//ErrNotExist not found
var ErrNotExist = errors.New("Not Exist")

//ErrIsEmpty not found
var ErrIsEmpty = errors.New("is Empty")
