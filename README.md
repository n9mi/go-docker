# go-docker
Simple CRUD with role-based JWT authorization

## **Packages used**
- github.com/joho/godotenv/cmd/godotenv@latest
- github.com/labstack/echo/v4
- github.com/go-playground/validator/v10
- gorm.io/driver/postgres
- gorm.io/gorm 
- github.com/casbin/casbin/v2

## **Database structure**

## **Role**
- Users can have more than one role

|Role|List the blogs|View full version of the blog|Create a blog|Update a blog|Delete a blog|
|--------|--------|--------|--------|--------|--------|
|Admin|[x]|[x]|[x]|[x]|[x]|
|Creator|[x]|[x]|[x]|[x]|[ ]|
|Subscriber|[x]|[x]|[ ]|[ ]|[ ]|
|Guest|[x]|[ ]|[ ]|[ ]|[ ]|

- New user can't register as an admin
- Admin can register existing account as another admin [TODO]
- Creator only can perform update and delete on their own blog [TODO]


## **Structure**
Based on repository pattern, this project use:
- Repository layer: For accessing db in the behalf of project to store/update/delete data
- Service layer: Contains set of logic/action needed to process data/orchestrate those data
- Models layer: Contains set of entity/actual data attribute
- Controller layer: Acts to mapping users input/request and presented it back to user as relevant responses

## **API Endpoints**
Swagger docs 
```
    127.0.0.1:8080
```