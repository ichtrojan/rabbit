# Rabbit

![rabbit-hero](https://github.com/ichtrojan/rabbit/assets/5338836/1dcbd594-5ce2-4a1a-a848-5dc6b9e2e07d)

## Introduction

If you ever want a way to trigger your Laravel jobs from Golang to leverage the features that [Laravel Horizon](https://laravel.com/docs/10.x/horizon) provides. Rabbit is the package for you.

## Installation

Run the following command to install `rabbit` on your Go project:

```bash
go get github.com/ichtrojan/rabbit
```

## Config

```go
queue := rabbit.Config{
    AppName: "horizon",
    Conn:    client,
    Job:     "App\\Jobs\\SendEmail",
    Queue:   "default",
    Delay:   10,
}
```

| Attribute	   | Description 	                                             |
|--------------|-----------------------------------------------------------|
| 	  `AppName` | 	Your `APP_NAME` as set in your laravel's app `.env` file |  
| 	   `Conn`   | 	    Redis connection                                     |
| 	   `Job`    | 	    The laravel job namespace                            | 
| 	   `Queue`  | 	    The queue name, the default queue name is `default`  |
| 	   `Delay`  | 	    Delay in seconds                                     |

## Parameters

The parameters should mirror the construct on your laravel job. For example, if you have a job with a construct like this:

```php
<?php
...
class ExampleJob {
    public $id;
    private $address;
    protected $message;

    public function __construct($id, $address, $message)
    {
        $this->id = $id;
        $this->address = $address;
        $this->message = $message;
    }
    ...
}
```

Your horizon parameter definition should look like this:

```go
params := []rabbit.Param{
    {Type: "public", Name: "id", Value: "0000-000000-00000-00000"},
    {Type: "private", Name: "address", Value: "1 Apple Park Way, Cupertino, California, USA"},
    {Type: "protected", Name: "message", Value: "stay foolish, stay hungry"},
}
```

> **NOTE**
> * Supported parameter types are `public`, `private` and `protected`.
> * Objects/Models cannot be passed as parameters and modifying your Job on laravel may be required.
> * Declaring an invalid parameter type would default to `public`.

## Usage

Ensure you have [Laravel Horizon](https://laravel.com/docs/10.x/horizon) installed and running on your laravel application, you can do that by running:

```bash
php artisan horizon
```

Alternatively, this would also work if you use the default queue worker that comes prepackaged with laravel.

```bash
php artisan queue:work redis
```

Here's a sample code with everything defined:

```go
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ichtrojan/rabbit"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", "127.0.0.1", "6379"),
		Password: "",
		DB:       0,
	})

	queue := rabbit.Config{
		AppName: "rabbit",
		Conn:    client,
		Job:     "App\\Jobs\\SendEmail",
		Queue:   "default",
		Delay:   10,
	}

	params := []rabbit.Param{
		{Type: "public", Name: "id", Value: "0000-000000-00000-00000"},
		{Type: "private", Name: "address", Value: "1 Apple Park Way, Cupertino, California, USA"},
		{Type: "protected", Name: "message", Value: "stay foolish, stay hungry"},
	}

	if err := queue.Dispatch(params...); err != nil {
		fmt.Println(err)
	}
}

```

## Contributor(s)

- Ibukun Ajimoti - [GitHub](https://github.com/ajimoti) [Twitter](https://x.com/ajimotea)
