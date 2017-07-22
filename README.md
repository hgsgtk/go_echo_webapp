#  Basic
- [echo getting started](https://github.com/labstack/echo)
```
$ go run server.go

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.2.1
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1323
```

```
curl -F "name=Joe Smith" -F "email=joe@labstack.com" http://localhost:1323/save
name: Joe Smith, email: joe@labstack.com
```
#　template rendering
- [template rendering](https://echo.labstack.com/guide/templates)
