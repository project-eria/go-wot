# Client (consumer)

## Connect client
```
	myConsumer := consumer.New()
	td, _ := myConsumer.ConsumeURL("http://127.0.0.1:8888/")
```

## Read thing description
```
	fmt.Println(td.GetThingDescription().Title)
```

## Read thing property
```
    value, _ := td.ReadProperty("on")
```