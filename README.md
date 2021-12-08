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
- 
- shipment: 
```
DROP TABLE IF EXISTS `shipment`;
CREATE TABLE `shipment` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `UserID` int(6) NOT NULL,
  `Description` varchar(64),
	`Tracking` varchar(16),
	`Comment` varchar(255),
	`date` varchar(16) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```
- user:
```
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `WeChatID` int(6) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```
- Admin: (We don't want to make a complicated authorization, so I didn't encode the password)
```
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `Name` varchar(16) NOT NULL,
	`password` varchar(16) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```


### api
#### frontend
- GET /api/v1/users/me
- GET /api/v1/shippings 

#### admin
- POST /api/admin/v1/users 
- POST /api/admin/v1/shippings 
- POST /api/v1/shippings/actions/resend-notification (for admin view)
- GET /api/admin/v1/shippings 
##### I have changed them to:
#### user
- to be done
#### admin
- POST admin/shipment (to create a new shipment)
- Get admin/shipment (to view all shipments)
- GET admin/shipments/{id} (to view one shipment)
- DELETE admin/shipments/{id} (to delete one shipment)
- PATCH admin/shipments/{id} (to update one shipment)
- to be done

### lib
ant design pro
https://pro.ant.design/

rest api
https://github.com/gorilla/mux

mysql golang driver 
https://github.com/go-sql-driver/mysql
example
https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
