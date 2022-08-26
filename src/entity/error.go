package entity

import "errors"

//ErrNotFound not found
var ErrNotFound = errors.New("ERR: Not found")

//ErrEmptyDB Empty DB or List
var ErrEmptyDB = errors.New("ERR: Empty DB")

//ErrInvalidEntity invalid entity
var ErrInvalidEntity = errors.New("ERR: Invalid entity")

//ErrCannotBeCreated cannot be Created
var ErrCannotBeCreated = errors.New("ERR: Cannot Be Created")

//ErrCannotBeDeleted cannot be deleted
var ErrCannotBeDeleted = errors.New("ERR: Cannot Be Deleted")

//ErrCannotBeUpdated cannot be Updated
var ErrCannotBeUpdated = errors.New("ERR: Cannot Be Updated")

//ErrNotExist not found
var ErrNotExist = errors.New("ERR: Not Exist")

//ErrIsEmpty not found
var ErrIsEmpty = errors.New("ERR: is Empty")

//ErrCannotConvertJSON Cannot Convert to Json Format
var ErrCannotConvertJSON = errors.New("JSON: Cannot Convert to Json Format")

//ErrIsEmpty not found
var ErrInvalidID = errors.New("ID: Invalid")
