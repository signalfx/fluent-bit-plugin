// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: args
func (_m *Logger) Debug(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Debugf provides a mock function with given fields: fmt, args
func (_m *Logger) Debugf(fmt string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, fmt)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: args
func (_m *Logger) Error(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Errorf provides a mock function with given fields: fmt, args
func (_m *Logger) Errorf(fmt string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, fmt)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: args
func (_m *Logger) Info(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Infof provides a mock function with given fields: fmt, args
func (_m *Logger) Infof(fmt string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, fmt)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Panic provides a mock function with given fields: args
func (_m *Logger) Panic(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Panicf provides a mock function with given fields: fmt, args
func (_m *Logger) Panicf(fmt string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, fmt)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: args
func (_m *Logger) Warn(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warnf provides a mock function with given fields: fmt, args
func (_m *Logger) Warnf(fmt string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, fmt)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}