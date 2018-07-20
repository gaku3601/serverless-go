package gombda

import "github.com/aws/aws-lambda-go/events"

type Gombda struct {
	request events.APIGatewayProxyRequest
	routes  []*route
}

type route struct {
	method   string
	resource string
	function func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func New(request events.APIGatewayProxyRequest) *Gombda {
	return &Gombda{request: request}
}

func (g *Gombda) GET(resource string, function func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) {
	g.routes = append(g.routes, &route{
		method:   "GET",
		resource: resource,
		function: function,
	})
}

func (g *Gombda) POST(resource string, function func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) {
	g.routes = append(g.routes, &route{
		method:   "POST",
		resource: resource,
		function: function,
	})
}

func (g *Gombda) DELETE(resource string, function func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) {
	g.routes = append(g.routes, &route{
		method:   "DELETE",
		resource: resource,
		function: function,
	})
}

func (g *Gombda) PATCH(resource string, function func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) {
	g.routes = append(g.routes, &route{
		method:   "PATCH",
		resource: resource,
		function: function,
	})
}

func (g *Gombda) OPTIONS(resource string, function func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) {
	g.routes = append(g.routes, &route{
		method:   "OPTIONS",
		resource: resource,
		function: function,
	})
}

func (g *Gombda) Start() (events.APIGatewayProxyResponse, error) {
	for _, v := range g.routes {
		if v.method == g.request.HTTPMethod && v.resource == g.request.Resource {
			return v.function(g.request)
		}
	}
	return events.APIGatewayProxyResponse{}, nil
}
