package entity

import "errors"

//ErrNotFound not found
var ErrNotFound = errors.New("ERR: Not found")

//ErrInvalidEntity invalid entity
var ErrInvalidEntity = errors.New("ERR: Invalid entity")

//ErrCannotBeDeleted cannot be deleted
var ErrCannotBeDeleted = errors.New("ERR: Cannot Be Deleted")

//ErrCannotBeUpdated cannot be deleted
var ErrCannotBeUpdated = errors.New("ERR: Cannot Be Updated")

//ErrNotExist not found
var ErrNotExist = errors.New("ERR: Not Exist")

//ErrIsEmpty not found
var ErrIsEmpty = errors.New("ERR: is Empty")

//ErrCannotConvertJSON Cannot Convert to Json Format
var ErrCannotConvertJSON = errors.New("JSON: Cannot Convert to Json Format")

//ErrIsEmpty not found
var ErrInvalidID = errors.New("ID: Invalid")
