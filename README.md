# Eventbus
[![Test Coverage](https://codecov.io/gh/nepet/eventbus/master/graph/badge.svg)](https://codecov.io/gh/nepet/eventbus)
#### A simple Implementation of an event bus in golang

This repository is an implementation of the publisher/subscriber pattern in golang. It is meant to be as simple as
possible to cover as many cases as possible.

## Install

```
go get github.com/nepet/eventbus
```


### Event

```
type Event interface {
	IsEvent()
}
```

An event is something that gets published. The interface only requires for a `IsEvent()` method that should do nothing
and just ensures that a struct is meant to be an event. Any struct can implement this method to be used as an event.

### Handler

```
type Handler interface {
	Handle(interface{Event})
}
```

A handler handles an event without a return value. It is up to the developer how the
`Handle(interface{Event})` method is used.

### Mock

The repository also contains a handler mock and a mock for the event bus. This allows you to easily use the bus in unit
tests.

### Examples

The examples provide an exemplary implementation of a handler and an event.
