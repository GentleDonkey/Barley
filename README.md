# notifications
Project “Barley” because its euphony is “大卖”. 


## project description
- provide them a way to view their purchase history quickly, instead of checking it by looking back our chat history one by one. 
- allows customers to track the package easily.

## features
admin
- admin creates user account
- admin create shipping information (general product description)
- admin can trigger shipping notification by sms

user
- login
- check shipping history

## plans
### frontend user
- user authorization (http://example.com/?code=jwt)
- list shipping history page filter, sort 
- list shipping history page - action button - redirect to tracking website

### frontend admin
- admin login 
- admin creates user account
- list shipping history page filter, sort 
- list shipping history page - action button

### database
- design schema

### api
#### frontend
- GET /api/v1/users/me
- GET /api/v1/shippings 

#### admin
- POST /api/admin/v1/users 
- POST /api/admin/v1/shippings 
- POST /api/v1/shippings/actions/resend-notification (for admin view)
- GET /api/admin/v1/shippings

### lib
ant design pro
https://pro.ant.design/

rest api
https://github.com/gorilla/mux

mysql golang driver 
https://github.com/go-sql-driver/mysql
example
https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
