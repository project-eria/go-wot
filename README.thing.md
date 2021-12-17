# Thing (producer)
## Create a thing object
```
	mything, err := thing.New(
		"dev:ops:my-actuator-1234",
		"actuator1",
		"Actuator1 Example",
		"An actuator 1st example",
		[]string{},
	)
```

## Add a security scheme
```
	noSecurityScheme := securityScheme.NewNoSecurity()
	mything.AddSecurity("no_sec", noSecurityScheme)
```

## Add a boolean property
```
	booleanData := dataSchema.NewBoolean(false)
	property := interaction.NewProperty(
		"on",
		"On/Off",
		"Whether the device is turned on",
		booleanData,
	)
	mything.AddProperty(&property)
```

## Add an action
```
	action := interaction.NewAction(
		"a",
		"No Input, No Output",
		"",
		nil,
		nil,
	)
	mything.AddAction(action)
```

## Launch the service
```
	myProducer := producer.New("127.0.0.1", 8888, false)
	exposedThing := myProducer.Produce(mything)
```

## Set property handlers
```
	exposedThing.SetPropertyReadHandler("on", propertyRead)
	exposedThing.SetPropertyWriteHandler("on", propertyWrite)
	
	...

	var (
		_on bool
	)
	
	_on = false

	func propertyRead() (interface{}, error) {
		return _on, nil
	}

	func propertyWrite(value interface{}) error {
		_on = value.(bool)
		return nil
	}
```

## Set action handler
```
	exposedThing.SetActionHandler("a", handlerA)

	...

	func handlerA(interface{}) (interface{}, error) {
		println("a action")
		return nil, nil
	}
```

## Expose the services
```
	myProducer.Expose()
```