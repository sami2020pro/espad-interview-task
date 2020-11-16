# aspad interview task
<b><var>aspad interview task</var></b>

<div>
  <img
      src="/data/task-completed.png"
      alt="aspad interview task | task completed"
      style="max-width:100%;"
  />
</div>

# Preview and Codes
<div>
  <img
      src="/data/tree-the-task-codes.png"
      alt="aspad interview task | tree the task codes"
      style="max-width:100%;"
  />
</div>

# What is this project
This project is the task of <strong>Spad Company</strong>

# TODOs
We need create and write:

        1) Generate a unque alias for the provided address

        2) The service should redirect users to original URLs when they access a short link 
        
        3) The short link hace a lifetime
        
        4) The service should have a visits count 
        
        5) Use camelCase
        
        6) Clean code

# Does this project have any dependencies
***Yes***

# What technologies did we use
We used these technologies

  1. **Golang** 
  
  2. **Golang libraries**
  
  3. **Redis**

# How to create a unique identifier for URL
<div>
We paid attention to these
  
  1. **Length**
  
  2. **Uniqueness**
  
  3. **Search**
</div>

# What settings did we use for this project
We put all the <var>Redis</var> settings in one **Json** file

And you can find all the settings in the ``configuration.json`` file
  
And You can change your settings in the ``configuration.json`` file

# How we worked with storage and Redis
We wrote a <var>Golang</var> code to do our job with **Redis**

And you can find this code in the ``storage.go`` file in the **storage directory**

# How we handled the whole project
We divided the tasks and wrote a *special code* for each task and wanted each section to do its job **efficiently**.

Base:
> base62.go

Config: 
> configuration.go

Handler:
> handler.go

Storage:
> storage.go

Main:
> main.go

# Why didn't we use Postman
We also had the ability to use <var>Postman</var>, but we didn't want it to get too crowded, so we used **Redis** for fun.

# How to use this project
Run service

```golang
go run main.go
```

Create a short link 

```shell
curl -L -X POST 'localhost:8080/encode' -H 'Content-Type: application/json' --data-raw '{
    "url": "WRITE-YOUR-URL-IN-HERE",
    "expires": "WRITE-YOUR-EXPIRE-TIME-IN-HERE"
}'
```

we will get a **response** like:
```json
{
   "success": true,
   "shortUrl": "http://localhost:8080/THE-SHORT-URL"
}
```

Any short url can be included in the **THE-SHORT-URL** 

Get detailed **information** for the short link

```shell
curl -L -X GET 'http://localhost:8080/THE-SHORT-URL'
```

Open the browser and follow the **your short link**
```html
http://localhost:8080/THE-SHORT-URL
```

# Examples
