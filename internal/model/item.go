package model

type Item struct {
	Name  string
	Price int
}

var Items = map[string]Item{
	"t-shirt":    {Price: 80},
	"cup":        {Price: 20},
	"book":       {Price: 50},
	"pen":        {Price: 10},
	"powerbank":  {Price: 200},
	"hoody":      {Price: 300},
	"umbrella":   {Price: 200},
	"socks":      {Price: 10},
	"wallet":     {Price: 50},
	"pink-hoody": {Price: 500},
}
